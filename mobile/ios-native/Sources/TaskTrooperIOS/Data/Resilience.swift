import Foundation

/// Thrown instead of calling the API while the breaker is open.
struct CircuitOpenError: Error, Equatable {
    let retryAfter: TimeInterval
}

/// Thrown when the client-side throttle rejects a call.
struct RateLimitedError: Error, Equatable {}

/// Client-side circuit breaker + token-bucket throttle in one actor, so the
/// state is safe under Swift concurrency without locks.
///
/// A phone that keeps retrying a dead API drains the battery and delays the
/// moment the UI can honestly say "offline", so failing fast is a UX decision
/// here, not only an infrastructure one.
///
/// The guard deliberately does NOT run the request itself: handing a closure
/// across the actor boundary would drag Sendable requirements into every call
/// site. Callers ask for permission, do their own I/O, and report the outcome.
actor NetworkGuard {
    enum State: Equatable { case closed, open, halfOpen }

    private let failureThreshold: Int
    private let openDuration: TimeInterval
    private let perSecond: Double
    private let burst: Double
    private let now: @Sendable () -> Date

    private var state: State = .closed
    private var failures = 0
    private var openedAt: Date?
    private var tokens: Double
    private var lastRefill: Date?

    init(
        failureThreshold: Int = 4,
        openDuration: TimeInterval = 20,
        perSecond: Double = 8,
        burst: Double = 16,
        now: @escaping @Sendable () -> Date = Date.init
    ) {
        self.failureThreshold = failureThreshold
        self.openDuration = openDuration
        self.perSecond = perSecond
        self.burst = burst
        self.now = now
        self.tokens = burst
    }

    var currentState: State {
        refresh()
        return state
    }

    /// Throws when the circuit is open or the throttle is empty; returns
    /// normally when the caller may perform the request.
    func acquire() throws {
        refresh()
        if state == .open, let openedAt {
            throw CircuitOpenError(retryAfter: openDuration - now().timeIntervalSince(openedAt))
        }
        guard consumeToken() else { throw RateLimitedError() }
    }

    func recordSuccess() {
        state = .closed
        failures = 0
    }

    /// Only call this for a genuine dependency failure (transport error or
    /// 5xx). A 404 or a decoding bug says nothing about the backend's health.
    func recordFailure() {
        if state == .halfOpen {
            state = .open
            openedAt = now()
            return
        }
        failures += 1
        if failures >= failureThreshold {
            state = .open
            openedAt = now()
        }
    }

    private func consumeToken() -> Bool {
        guard perSecond > 0 else { return true }
        let current = now()
        let last = lastRefill ?? current
        tokens = min(burst, tokens + current.timeIntervalSince(last) * perSecond)
        lastRefill = current
        guard tokens >= 1 else { return false }
        tokens -= 1
        return true
    }

    private func refresh() {
        if state == .open, let openedAt, now().timeIntervalSince(openedAt) >= openDuration {
            state = .halfOpen
        }
    }
}
