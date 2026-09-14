import { NextResponse, type NextRequest } from "next/server";
import { RateLimiter } from "@/lib/resilience/rate-limit";

// Rate limit every API route by client IP. Defaults are deliberately generous
// (a normal user never hits them) and configurable per deployment; the point is
// that a fresh project is not born with an unmetered public endpoint.
const limiter = new RateLimiter(
  Number(process.env.RATE_LIMIT_MAX ?? 120),
  Number(process.env.RATE_LIMIT_WINDOW_MS ?? 60_000),
);

export function middleware(request: NextRequest) {
  const ip =
    request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() ??
    request.headers.get("x-real-ip") ??
    "unknown";

  const result = limiter.check(ip);
  if (!result.allowed) {
    return NextResponse.json(
      { error: "rate limit exceeded" },
      {
        status: 429,
        headers: { "Retry-After": String(Math.ceil(result.resetMs / 1000)) },
      },
    );
  }

  const response = NextResponse.next();
  response.headers.set("X-RateLimit-Remaining", String(result.remaining));
  return response;
}

export const config = {
  // Only API routes: page/static requests are served from the edge cache and
  // do not need (or want) a per-IP counter.
  matcher: ["/api/:path*"],
};
