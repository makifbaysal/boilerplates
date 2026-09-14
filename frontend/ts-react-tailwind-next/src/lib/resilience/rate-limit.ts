// In-memory fixed-window rate limiter for Next.js route handlers and
// middleware. It is per-instance on purpose: it stops a single client from
// hammering one server without adding Redis to a starter project. Swap the
// store for Redis (or the platform's own limiter) before running more than one
// instance behind a load balancer.

export interface RateLimitResult {
  allowed: boolean;
  remaining: number;
  resetMs: number;
}

interface Bucket {
  count: number;
  resetAt: number;
}

export class RateLimiter {
  private readonly buckets = new Map<string, Bucket>();

  constructor(
    private readonly max = 60,
    private readonly windowMs = 60_000,
    private readonly now: () => number = Date.now,
  ) {}

  check(key: string): RateLimitResult {
    const now = this.now();
    const bucket = this.buckets.get(key);

    if (!bucket || bucket.resetAt <= now) {
      this.buckets.set(key, { count: 1, resetAt: now + this.windowMs });
      this.sweep(now);
      return { allowed: true, remaining: this.max - 1, resetMs: this.windowMs };
    }

    bucket.count += 1;
    const remaining = Math.max(0, this.max - bucket.count);
    return { allowed: bucket.count <= this.max, remaining, resetMs: bucket.resetAt - now };
  }

  /** Drops expired buckets so a long-running process does not leak memory. */
  private sweep(now: number): void {
    if (this.buckets.size < 1000) return;
    for (const [key, bucket] of this.buckets) {
      if (bucket.resetAt <= now) this.buckets.delete(key);
    }
  }
}
