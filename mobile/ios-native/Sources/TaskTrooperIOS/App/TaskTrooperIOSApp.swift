import SwiftUI

@main
struct TaskTrooperIOSApp: App {
    private let apiClient = URLSessionTaskAPIClient(baseURL: AppConfig.defaultAPIBaseURL)

    var body: some Scene {
        WindowGroup {
            TaskListView(apiClient: apiClient)
        }
    }
}
