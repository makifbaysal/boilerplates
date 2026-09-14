import SwiftUI

struct TaskListView: View {
    @State private var viewModel: TaskListViewModel
    @State private var newTitle = ""

    init(apiClient: TaskAPIClientProtocol) {
        _viewModel = State(initialValue: TaskListViewModel(apiClient: apiClient))
    }

    var body: some View {
        NavigationStack {
            VStack(spacing: 0) {
                if let message = viewModel.errorMessage {
                    Text(message)
                        .font(.footnote)
                        .foregroundStyle(.red)
                        .padding(.horizontal)
                }

                List {
                    ForEach(viewModel.tasks) { task in
                        TaskRowView(task: task) {
                            Task { await viewModel.toggleDone(task) }
                        }
                    }
                    .onDelete { offsets in
                        for index in offsets {
                            let task = viewModel.tasks[index]
                            Task { await viewModel.delete(task) }
                        }
                    }
                }
                .listStyle(.plain)
                .overlay {
                    if viewModel.isLoading && viewModel.tasks.isEmpty {
                        ProgressView()
                    } else if viewModel.tasks.isEmpty {
                        ContentUnavailableView(
                            "No tasks yet", systemImage: "checklist",
                            description: Text("Add one below."))
                    }
                }

                HStack {
                    TextField("New task", text: $newTitle)
                        .textFieldStyle(.roundedBorder)
                    Button("Add") {
                        let title = newTitle
                        newTitle = ""
                        Task { await viewModel.addTask(title: title) }
                    }
                    .disabled(newTitle.trimmingCharacters(in: .whitespaces).isEmpty)
                }
                .padding()
            }
            .navigationTitle("Tasks")
            .task { await viewModel.load() }
            .refreshable { await viewModel.load() }
        }
    }
}

#Preview {
    TaskListView(apiClient: URLSessionTaskAPIClient(baseURL: URL(string: "http://localhost:8080")!))
}
