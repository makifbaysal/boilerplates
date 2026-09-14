import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/task_providers.dart';
import '../molecules/task_tile.dart';

class TaskListView extends ConsumerWidget {
  const TaskListView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tasksAsync = ref.watch(taskListProvider);

    return tasksAsync.when(
      data: (tasks) {
        if (tasks.isEmpty) {
          return const Center(child: Text('No tasks yet.'));
        }
        return ListView.builder(
          itemCount: tasks.length,
          itemBuilder: (context, i) {
            final task = tasks[i];
            return TaskTile(
              task: task,
              onToggleDone: (done) => ref.read(taskListProvider.notifier).setDone(task.id, done),
              onDelete: () => ref.read(taskListProvider.notifier).delete(task.id),
            );
          },
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (err, _) => Center(child: Text('Could not load tasks: $err')),
    );
  }
}
