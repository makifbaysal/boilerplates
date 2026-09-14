import Foundation
import Testing
@testable import TaskTrooperIOS

/// Mocks only the true I/O boundary (URLProtocol intercepting
/// URLSession) and exercises the real URLSessionTaskAPIClient on top —
/// same principle as MockWebServer in mobile/android-native and
/// http.Client mocking in mobile/flutter.
final class MockURLProtocol: URLProtocol {
    nonisolated(unsafe) static var handler: ((URLRequest) throws -> (HTTPURLResponse, Data))?

    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }

    override func startLoading() {
        guard let handler = Self.handler else {
            fatalError("MockURLProtocol.handler not set")
        }
        do {
            let (response, data) = try handler(request)
            client?.urlProtocol(self, didReceive: response, cacheStoragePolicy: .notAllowed)
            client?.urlProtocol(self, didLoad: data)
            client?.urlProtocolDidFinishLoading(self)
        } catch {
            client?.urlProtocol(self, didFailWithError: error)
        }
    }

    override func stopLoading() {}
}

@Suite(.serialized)
struct TaskAPIClientTests {
    private func makeClient() -> URLSessionTaskAPIClient {
        let config = URLSessionConfiguration.ephemeral
        config.protocolClasses = [MockURLProtocol.self]
        let session = URLSession(configuration: config)
        return URLSessionTaskAPIClient(baseURL: URL(string: "https://example.test")!, session: session)
    }

    @Test
    func listTasksDecodesSnakeCaseFields() async throws {
        let json = """
        [{"id":"1","title":"Write tests","done":false,
          "created_at":"2026-01-01T10:00:00Z","updated_at":"2026-01-01T10:00:00Z"}]
        """.data(using: .utf8)!

        MockURLProtocol.handler = { request in
            #expect(request.httpMethod == "GET")
            #expect(request.url?.path == "/tasks")
            let response = HTTPURLResponse(
                url: request.url!, statusCode: 200,
                httpVersion: nil, headerFields: nil)!
            return (response, json)
        }

        let tasks = try await makeClient().listTasks()
        #expect(tasks.count == 1)
        #expect(tasks[0].id == "1")
        #expect(tasks[0].title == "Write tests")
        #expect(tasks[0].done == false)
    }

    @Test
    func createTaskSendsTitleAndReturnsCreatedTask() async throws {
        let json = """
        {"id":"2","title":"New task","done":false,
         "created_at":"2026-01-01T10:00:00Z","updated_at":"2026-01-01T10:00:00Z"}
        """.data(using: .utf8)!

        MockURLProtocol.handler = { request in
            #expect(request.httpMethod == "POST")
            let body = try #require(request.httpBodyOrStream())
            let decoded = try JSONDecoder().decode(CreateTaskRequest.self, from: body)
            #expect(decoded.title == "New task")
            let response = HTTPURLResponse(
                url: request.url!, statusCode: 201,
                httpVersion: nil, headerFields: nil)!
            return (response, json)
        }

        let task = try await makeClient().createTask(title: "New task")
        #expect(task.id == "2")
    }

    @Test
    func httpErrorStatusThrows() async throws {
        MockURLProtocol.handler = { request in
            let response = HTTPURLResponse(
                url: request.url!, statusCode: 404,
                httpVersion: nil, headerFields: nil)!
            return (response, Data())
        }

        await #expect(throws: TaskAPIError.http(status: 404)) {
            _ = try await makeClient().listTasks()
        }
    }

    @Test
    func deleteTaskSendsCorrectPathAndMethod() async throws {
        MockURLProtocol.handler = { request in
            #expect(request.httpMethod == "DELETE")
            #expect(request.url?.path == "/tasks/42")
            let response = HTTPURLResponse(
                url: request.url!, statusCode: 204,
                httpVersion: nil, headerFields: nil)!
            return (response, Data())
        }

        try await makeClient().deleteTask(id: "42")
    }
}

private extension URLRequest {
    /// URLProtocol sees httpBody as nil for streamed bodies; URLSession
    /// swaps it into the stream form before dispatching. Read it back
    /// out for assertions.
    func httpBodyOrStream() -> Data? {
        if let body = httpBody { return body }
        guard let stream = httpBodyStream else { return nil }
        stream.open()
        defer { stream.close() }
        var data = Data()
        let bufferSize = 4096
        var buffer = [UInt8](repeating: 0, count: bufferSize)
        while stream.hasBytesAvailable {
            let read = stream.read(&buffer, maxLength: bufferSize)
            if read > 0 { data.append(buffer, count: read) }
        }
        return data
    }
}
