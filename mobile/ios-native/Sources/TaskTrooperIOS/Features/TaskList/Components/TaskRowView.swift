import SwiftUI

/// Pure view: props in, callbacks out — mirrors mobile/android-native's
/// TaskItem composable and mobile/flutter's TaskRow widget.
struct TaskRowView: View {
    let task: TaskItem
    let onToggle: () -> Void

    var body: some View {
        Button(action: onToggle) {
            HStack {
                Image(systemName: task.done ? "checkmark.circle.fill" : "circle")
                    .foregroundStyle(task.done ? .green : .secondary)
                Text(task.title)
                    .strikethrough(task.done)
                    .foregroundStyle(task.done ? .secondary : .primary)
                Spacer()
            }
        }
        .buttonStyle(.plain)
    }
}

#Preview {
    TaskRowView(
        task: TaskItem(id: "1", title: "Write SwiftUI boilerplate", done: false,
                       createdAt: .now, updatedAt: .now),
        onToggle: {}
    )
}
