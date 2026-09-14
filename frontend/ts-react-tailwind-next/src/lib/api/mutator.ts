import { CircuitBreaker } from "@/lib/resilience/circuit-breaker";

// The fetch implementation orval's generated hooks call through —
// swap the base URL resolution or add auth headers here, once, rather
// than in every generated hook.
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// One breaker for the whole API: when it is down, every hook fails fast
// instead of queueing doomed requests behind a 30s timeout.
const breaker = new CircuitBreaker({
  failureThreshold: Number(process.env.NEXT_PUBLIC_BREAKER_FAILURES ?? 5),
  openMs: Number(process.env.NEXT_PUBLIC_BREAKER_OPEN_MS ?? 30_000),
});

export const customFetch = async <T>(url: string, options: RequestInit = {}): Promise<T> => {
  // Only transport errors and 5xx count against the breaker: a 404 or a 422 is
  // this request being wrong, not the API being down, and must not trip it.
  const response = await breaker.run(async () => {
    const res = await fetch(`${API_BASE_URL}${url}`, {
      ...options,
      headers: { "Content-Type": "application/json", ...options.headers },
    });
    if (res.status >= 500) {
      throw new ApiError(res.status, await res.text());
    }
    return res;
  });

  if (!response.ok) {
    throw new ApiError(response.status, await response.text());
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
};

export default customFetch;
