# Architecture

SwiftUI + Swift 6 (strict concurrency), `@Observable` view models
(Observation framework, not Combine/`ObservableObject`), URLSession for
the only I/O boundary. The Xcode project is **generated, not
committed** — see Project generation below.

```
Sources/TaskTrooperIOS/
  App/TaskTrooperIOSApp.swift    @main entry point
  App/Config.swift                 default API base URL
  Data/TaskItem.swift               Codable model + request DTOs
  Data/TaskAPIClient.swift          port (protocol) + URLSession adapter
  Features/TaskList/TaskListViewModel.swift   @Observable, @MainActor
  Features/TaskList/TaskListView.swift          screen
  Features/TaskList/Components/TaskRowView.swift  pure view, props/callbacks
Tests/TaskTrooperIOSTests/
  TaskAPIClientTests.swift        Swift Testing + URLProtocol mock
```

## Project generation

`project.yml` (XcodeGen spec) is checked in; `TaskTrooperIOS.xcodeproj`
is not — it's generated and gitignored, same reasoning every other
XcodeGen-based iOS codebase uses: a generated `.xcodeproj` never causes
merge conflicts. Run `xcodegen generate` after `git clone` and again
**any time you add, remove, or rename a source file** — the generated
project's file list is a snapshot of disk state at generate-time; a
stale one fails the build with "Build input file cannot be found"
pointing at the old filename. CI regenerates it on every run for the
same reason.

## Real bugs found building this

- **A domain type named `Task` shadows Swift Concurrency's own `Task`.**
  Every other boilerplate in this repo names the domain entity `Task`;
  in Swift that collides with `_Concurrency.Task`, and every
  `Task { await ... }` call site in the same module silently resolves
  to the *domain* type's initializer instead (a confusing "trailing
  closure passed to parameter of type `Decoder`" error, not an obvious
  "redeclaration" one). Renamed to `TaskItem` — the one Swift-specific
  deviation from the shared cross-boilerplate name, and worth knowing
  before you rename it back.
- **`ISO8601DateFormatter` isn't `Sendable`**, so a `static let` array of
  them fails Swift 6's default strict-concurrency check ("static
  property ... may have shared mutable state"). These two formatter
  instances are configured once and only ever read via `.date(from:)`
  afterward, so `nonisolated(unsafe)` is the correct escape hatch here
  — see `TaskAPIClient.swift`.
- **`CreateTaskRequest`/`SetDoneRequest` need `Codable`, not just
  `Encodable`**, even though the app only ever encodes them — tests
  that assert on the exact JSON body sent (rather than string-matching
  raw bytes) need to decode it back, which requires `Decodable` too.

## Config

`AppConfig.defaultAPIBaseURL` (`App/Config.swift`) defaults to
`http://localhost:8080`. Unlike Android's emulator, the iOS Simulator
shares the host machine's network namespace, so plain `localhost`
reaches a `backend/*` boilerplate running on your dev machine directly
— no NAT alias needed. Point it at a device or deployed backend by
editing this constant (a real project would likely wire this to
per-scheme build settings; kept as a plain constant here to stay
focused on the CRUD example).

## Testing

`Tests/TaskTrooperIOSTests` uses **Swift Testing** (`import Testing`,
`@Test`, `#expect`), not XCTest — Xcode 26's idiomatic default. Tests
mock only the true I/O boundary: a `URLProtocol` subclass intercepts
requests at the `URLSession` level, and the real
`URLSessionTaskAPIClient` runs on top of it unmodified — same principle
as MockWebServer in `mobile/android-native` and `http.Client` mocking in
`mobile/flutter`. No UI tests are included (they'd need
`XCUIApplication` and add real runtime cost in CI); add them under a
`TaskTrooperIOSUITests` target via `project.yml` if your project needs
them.

Running the test target inside the CI/local simulator also launches the
full app (SwiftUI Testing bundles run in-process in the host app), so
you'll see one harmless connection-refused log line from
`TaskListView`'s `.task { await viewModel.load() }` trying to reach
`localhost:8080` with nothing listening — the view model catches it and
sets `errorMessage`, it doesn't fail the test run.
