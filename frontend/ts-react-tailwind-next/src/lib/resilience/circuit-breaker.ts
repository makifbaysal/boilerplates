// A dependency-free circuit breaker for outbound calls. A frontend that keeps
// hammering a dead API burns the user's battery, fills the console with noise
// and delays the moment the UI can say "this is down" — so failing fast is a
// UX decision as much as an infrastructure one.

export type BreakerState = "closed" | "open" | "half-open";

export class CircuitOpenError extends Error {
  constructor(readonly retryAfterMs: number) {
    super("circuit breaker is open");
    this.name = "CircuitOpenError";
  }
}

export interface BreakerOptions {
  /** Consecutive failures that open the circuit. */
  failureThreshold?: number;
  /** How long the circuit stays open before a probe is allowed. */
  openMs?: number;
  /** Probe successes needed to close it again. */
  successThreshold?: number;
  /** Injectable clock for tests. */
  now?: () => number;
}

export class CircuitBreaker {
  private state: BreakerState = "closed";
  private failures = 0;
  private successes = 0;
  private openedAt = 0;
  private readonly failureThreshold: number;
  private readonly openMs: number;
  private readonly successThreshold: number;
  private readonly now: () => number;

  constructor(options: BreakerOptions = {}) {
    this.failureThreshold = options.failureThreshold ?? 5;
    this.openMs = options.openMs ?? 30_000;
    this.successThreshold = options.successThreshold ?? 1;
    this.now = options.now ?? Date.now;
  }

  currentState(): BreakerState {
    this.refresh();
    return this.state;
  }

  async run<T>(fn: () => Promise<T>): Promise<T> {
    this.refresh();
    if (this.state === "open") {
      throw new CircuitOpenError(this.openMs - (this.now() - this.openedAt));
    }
    try {
      const result = await fn();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }

  private refresh(): void {
    if (this.state === "open" && this.now() - this.openedAt >= this.openMs) {
      this.state = "half-open";
      this.successes = 0;
    }
  }

  private onSuccess(): void {
    if (this.state === "half-open") {
      this.successes += 1;
      if (this.successes >= this.successThreshold) {
        this.state = "closed";
        this.failures = 0;
      }
      return;
    }
    this.failures = 0;
  }

  private onFailure(): void {
    if (this.state === "half-open") {
      this.state = "open";
      this.openedAt = this.now();
      return;
    }
    this.failures += 1;
    if (this.failures >= this.failureThreshold) {
      this.state = "open";
      this.openedAt = this.now();
    }
  }
}
