import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/task_api_client.dart';
import '../models/task.dart';

/// `flutter run --dart-define=API_BASE_URL=http://...` overrides this at
/// compile time; tests override the provider directly instead (see
/// widget_test.dart) since --dart-define isn't available there.
const _apiBaseUrlDefine = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');

final apiBaseUrlProvider = Provider<String>((ref) => _apiBaseUrlDefine);

final taskApiClientProvider = Provider<TaskApiClient>(
  (ref) => TaskApiClient(baseUrl: ref.watch(apiBaseUrlProvider)),
);

final taskListProvider =
    AsyncNotifierProvider<TaskListNotifier, List<Task>>(TaskListNotifier.new);

class TaskListNotifier extends AsyncNotifier<List<Task>> {
  @override
  Future<List<Task>> build() {
    return ref.read(taskApiClientProvider).listTasks();
  }

  Future<void> create(String title) async {
    await ref.read(taskApiClientProvider).createTask(title);
    ref.invalidateSelf();
    await future;
  }

  Future<void> setDone(String id, bool done) async {
    await ref.read(taskApiClientProvider).setTaskDone(id, done);
    ref.invalidateSelf();
    await future;
  }

  Future<void> delete(String id) async {
    await ref.read(taskApiClientProvider).deleteTask(id);
    ref.invalidateSelf();
    await future;
  }
}
