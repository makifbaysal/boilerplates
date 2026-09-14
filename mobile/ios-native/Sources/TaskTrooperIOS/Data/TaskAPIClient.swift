import Foundation

enum TaskAPIError: Error, Equatable {
    case invalidResponse
    case http(status: Int)
    case decoding
}

/// Port: the only I/O boundary this app talks to. Tests mock at the
/// URLProtocol level (see TaskAPIClientTests) and exercise the real
/// URLSessionTaskAPIClient on top, never a hand-rolled fake of this
/// protocol — same "mock only the true I/O boundary" rule as every
/// other boilerplate in this repo.
protocol TaskAPIClientProtocol: Sendable {
    func listTasks() async throws -> [TaskItem]
    func createTask(title: String) async throws -> TaskItem
    func setTaskDone(id: String, done: Bool) async throws -> TaskItem
    func deleteTask(id: String) async throws
}

struct URLSessionTaskAPIClient: TaskAPIClientProtocol {
    let baseURL: URL
    let session: URLSession
    private let guardRail: NetworkGuard

    init(baseURL: URL, session: URLSession = .shared, guardRail: NetworkGuard = NetworkGuard()) {
        self.baseURL = baseURL
        self.session = session
        self.guardRail = guardRail
    }

    /// Every request goes through the throttle + breaker. A 4xx is this
    /// request being wrong and must not open the circuit; transport errors and
    /// 5xx are the backend being down and must.
    private func guarded(_ request: URLRequest) async throws -> (Data, URLResponse) {
        try await guardRail.acquire()
        let data: Data
        let response: URLResponse
        do {
            (data, response) = try await session.data(for: request)
        } catch {
            await guardRail.recordFailure()
            throw error
        }
        if let http = response as? HTTPURLResponse, http.statusCode >= 500 {
            await guardRail.recordFailure()
        } else {
            await guardRail.recordSuccess()
        }
        return (data, response)
    }

    private static let decoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .custom { decoderContainer in
            let container = try decoderContainer.singleValueContainer()
            let raw = try container.decode(String.self)
            for formatter in dateFormatters {
                if let date = formatter.date(from: raw) {
                    return date
                }
            }
            throw DecodingError.dataCorruptedError(
                in: container, debugDescription: "Unrecognized date: \(raw)")
        }
        return decoder
    }()

    // ISO8601DateFormatter isn't Sendable, but these two instances are
    // configured once and only ever read from (.date(from:)) — safe to
    // opt out of Swift 6's static-property concurrency check here.
    nonisolated(unsafe) private static let dateFormatters: [ISO8601DateFormatter] = {
        let withFractional = ISO8601DateFormatter()
        withFractional.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        let plain = ISO8601DateFormatter()
        plain.formatOptions = [.withInternetDateTime]
        return [withFractional, plain]
    }()

    private static let encoder = JSONEncoder()

    func listTasks() async throws -> [TaskItem] {
        try await send(request: makeRequest(path: "/tasks", method: "GET"))
    }

    func createTask(title: String) async throws -> TaskItem {
        var request = makeRequest(path: "/tasks", method: "POST")
        request.httpBody = try Self.encoder.encode(CreateTaskRequest(title: title))
        return try await send(request: request)
    }

    func setTaskDone(id: String, done: Bool) async throws -> TaskItem {
        var request = makeRequest(path: "/tasks/\(id)", method: "PATCH")
        request.httpBody = try Self.encoder.encode(SetDoneRequest(done: done))
        return try await send(request: request)
    }

    func deleteTask(id: String) async throws {
        let request = makeRequest(path: "/tasks/\(id)", method: "DELETE")
        let (_, response) = try await guarded(request)
        try Self.validate(response)
    }

    private func makeRequest(path: String, method: String) -> URLRequest {
        var request = URLRequest(url: baseURL.appendingPathComponent(path))
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        return request
    }

    private func send<T: Decodable>(request: URLRequest) async throws -> T {
        let (data, response) = try await guarded(request)
        try Self.validate(response)
        do {
            return try Self.decoder.decode(T.self, from: data)
        } catch {
            throw TaskAPIError.decoding
        }
    }

    private static func validate(_ response: URLResponse) throws {
        guard let http = response as? HTTPURLResponse else {
            throw TaskAPIError.invalidResponse
        }
        guard (200..<300).contains(http.statusCode) else {
            throw TaskAPIError.http(status: http.statusCode)
        }
    }
}
