# flutter boilerplate

Flutter + Riverpod Task CRUD app calling the same REST shape every
`backend/*` boilerplate exposes. Freezed/json_serializable models
(code-gen — see `.ai/common-tasks.md`), atoms/molecules/organisms
layering under `lib/widgets/` matching `frontend/ts-react-tailwind-next`'s
convention.

```
lib/
  models/task.dart          freezed model (+ .freezed.dart/.g.dart, generated)
  api/task_api_client.dart   http-based client — the only I/O boundary
  providers/task_providers.dart   Riverpod AsyncNotifier
  widgets/{atoms,molecules,organisms}/
  screens/task_screen.dart
```

## Commands

```
flutter pub get
dart run build_runner build   # after editing models/task.dart
flutter analyze
flutter test
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

Set the API base URL via `apiBaseUrlProvider` override (see
`main.dart`/tests) — point it at go-fiber/java-spring/quarkus.
