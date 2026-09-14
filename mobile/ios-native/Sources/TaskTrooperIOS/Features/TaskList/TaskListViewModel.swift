import Foundation
import Observation

@MainActor
@Observable
final class TaskListViewModel {
    private(set) var tasks: [TaskItem] = []
    private(set) var isLoading = false
    var errorMessage: String?

    private let apiClient: TaskAPIClientProtocol

    init(apiClient: TaskAPIClientProtocol) {
        self.apiClient = apiClient
    }

    func load() async {
        isLoading = true
        defer { isLoading = false }
        do {
            tasks = try await apiClient.listTasks()
            errorMessage = nil
        } catch {
            errorMessage = "Failed to load tasks."
        }
    }

    func addTask(title: String) async {
        let trimmed = title.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else { return }
        do {
            let created = try await apiClient.createTask(title: trimmed)
            tasks.append(created)
            errorMessage = nil
        } catch {
            errorMessage = "Failed to create task."
        }
    }

    func toggleDone(_ task: TaskItem) async {
        do {
            let updated = try await apiClient.setTaskDone(id: task.id, done: !task.done)
            if let index = tasks.firstIndex(where: { $0.id == updated.id }) {
                tasks[index] = updated
            }
            errorMessage = nil
        } catch {
            errorMessage = "Failed to update task."
        }
    }

    func delete(_ task: TaskItem) async {
        do {
            try await apiClient.deleteTask(id: task.id)
            tasks.removeAll { $0.id == task.id }
            errorMessage = nil
        } catch {
            errorMessage = "Failed to delete task."
        }
    }
}
