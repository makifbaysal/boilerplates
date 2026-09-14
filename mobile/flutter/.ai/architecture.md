# Architecture

Layering: `api/` (I/O boundary — only this file talks to the network) →
`providers/` (Riverpod `AsyncNotifier`, business/state logic) →
`widgets/{atoms,molecules,organisms}/` → `screens/`. Same shape as
`frontend/ts-react-tailwind-next`: atoms/molecules are pure
props-in/callbacks-out, organisms are the only widgets that read
providers.

## Code generation

`models/task.dart` is `@freezed` + `@JsonSerializable` — run
`dart run build_runner build` after editing it (regenerates
`task.freezed.dart`/`task.g.dart`, both committed). Keep
`freezed_annotation`/`json_annotation` in `dependencies` (not
`dev_dependencies`) — the generated code imports them at runtime, only
`build_runner`/`freezed`/`json_serializable` themselves are
dev-only/build-time.

## Testing

Mock only the true I/O boundary (`http.Client`, via
`package:http/testing.dart`'s `MockClient`) — same principle as every
other boilerplate in this repo. `test/api/task_api_client_test.dart`
tests the client directly; `test/widget_test.dart` overrides
`taskApiClientProvider` with a client backed by the same `MockClient`
and renders the real `TaskScreen`/providers/widgets end to end.

## Add a new endpoint

Add the method to `TaskApiClient`, a corresponding method to
`TaskListNotifier` (call the client, then `ref.invalidateSelf()` +
`await future`), then whichever widget needs to call it. Add a test
for the client method and, if UI changed, a widget test.
