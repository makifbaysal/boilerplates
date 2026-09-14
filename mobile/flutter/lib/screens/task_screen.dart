import 'package:flutter/material.dart';
import '../widgets/organisms/create_task_form.dart';
import '../widgets/organisms/task_list_view.dart';

class TaskScreen extends StatelessWidget {
  const TaskScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Tasks')),
      body: const Column(
        children: [
          CreateTaskForm(),
          Expanded(child: TaskListView()),
        ],
      ),
    );
  }
}
