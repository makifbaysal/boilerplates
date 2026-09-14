package ai.tasktrooper.boilerplate.javaspring.platform.resilience;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.Duration;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

/**
 * Per-IP inbound rate limit, on by default. An API without one is a single
 * misbehaving client away from an outage, and that client is usually our own
 * retry loop.
 *
 * <p>Tune with app.rate-limit.max / app.rate-limit.window-seconds, disable with
 * app.rate-limit.enabled=false (local load tests).
 */
@Component
public class RateLimitFilter extends OncePerRequestFilter {

    private final boolean enabled;
    private final RateLimiter limiter;

    public RateLimitFilter(
            @Value("${app.rate-limit.enabled:true}") boolean enabled,
            @Value("${app.rate-limit.max:120}") int max,
            @Value("${app.rate-limit.window-seconds:60}") int windowSeconds) {
        this.enabled = enabled;
        this.limiter = new RateLimiter(max, Duration.ofSeconds(windowSeconds));
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        // Platform health probes must never be throttled: a 429 there reads as
        // an outage and gets the instance restarted.
        String path = request.getRequestURI();
        return !enabled || path.startsWith("/actuator") || path.equals("/healthz");
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain chain)
            throws ServletException, IOException {
        RateLimiter.Decision decision = limiter.check(clientKey(request));
        response.setHeader("X-RateLimit-Remaining", String.valueOf(decision.remaining()));
        if (!decision.allowed()) {
            response.setStatus(HttpStatus.TOO_MANY_REQUESTS.value());
            response.setHeader("Retry-After", String.valueOf(decision.retryAfterSeconds()));
            response.setContentType("application/json");
            response.getWriter().write("{\"error\":\"rate limit exceeded\"}");
            return;
        }
        chain.doFilter(request, response);
    }

    private String clientKey(HttpServletRequest request) {
        String forwarded = request.getHeader("X-Forwarded-For");
        if (forwarded != null && !forwarded.isBlank()) {
            return forwarded.split(",")[0].trim();
        }
        return request.getRemoteAddr() != null ? request.getRemoteAddr() : "unknown";
    }
}
