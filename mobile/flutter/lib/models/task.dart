import 'package:freezed_annotation/freezed_annotation.dart';

part 'task.freezed.dart';
part 'task.g.dart';

/// Matches the Task shape every backend boilerplate in this repo
/// exposes (go-fiber/java-spring/quarkus REST /tasks).
@freezed
abstract class Task with _$Task {
  const factory Task({
    required String id,
    required String title,
    required bool done,
    @JsonKey(name: 'created_at') required DateTime createdAt,
    @JsonKey(name: 'updated_at') required DateTime updatedAt,
  }) = _Task;

  factory Task.fromJson(Map<String, dynamic> json) => _$TaskFromJson(json);
}
