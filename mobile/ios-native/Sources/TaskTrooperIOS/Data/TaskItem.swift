import Foundation

// Named `TaskItem`, not `Task` — Swift's own concurrency primitive is
// `Task`, and a same-module type named `Task` shadows it, silently
// breaking every `Task { await ... }` call site (they'd try to bind to
// this struct's memberwise/Decodable initializer instead). This is the
// one Swift-specific deviation from the shared `Task` domain name used
// by every other boilerplate in this repo.
struct TaskItem: Identifiable, Codable, Equatable, Sendable {
    let id: String
    var title: String
    var done: Bool
    let createdAt: Date
    let updatedAt: Date

    enum CodingKeys: String, CodingKey {
        case id, title, done
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}

struct CreateTaskRequest: Codable, Sendable {
    let title: String
}

struct SetDoneRequest: Codable, Sendable {
    let done: Bool
}
