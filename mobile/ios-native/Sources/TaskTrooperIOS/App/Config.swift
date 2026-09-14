import Foundation

enum AppConfig {
    /// iOS Simulator shares the host machine's network namespace, so
    /// plain `localhost` reaches a backend/* boilerplate running on
    /// your dev machine directly — unlike Android's emulator, no NAT
    /// alias (10.0.2.2) is needed. See .ai/architecture.md#config for
    /// how to point this at a device or a deployed backend instead.
    static let defaultAPIBaseURL = URL(string: "http://localhost:8080")!
}
