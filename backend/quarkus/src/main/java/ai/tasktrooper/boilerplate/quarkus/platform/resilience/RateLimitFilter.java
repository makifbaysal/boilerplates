package ai.tasktrooper.boilerplate.quarkus.platform.resilience;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.container.ContainerRequestFilter;
import jakarta.ws.rs.container.ContainerResponseContext;
import jakarta.ws.rs.container.ContainerResponseFilter;
import jakarta.ws.rs.core.Context;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.Provider;
import java.time.Duration;
import org.eclipse.microprofile.config.inject.ConfigProperty;

/**
 * Per-IP inbound rate limit, on by default. An API without one is a single
 * misbehaving client away from an outage.
 *
 * <p>Tune with app.rate-limit.max / app.rate-limit.window-seconds, disable with
 * app.rate-limit.enabled=false.
 */
@Provider
@ApplicationScoped
public class RateLimitFilter implements ContainerRequestFilter, ContainerResponseFilter {

    private static final String REMAINING_PROPERTY = "app.rate-limit.remaining";

    private final boolean enabled;
    private final RateLimiter limiter;

    public RateLimitFilter(
            @ConfigProperty(name = "app.rate-limit.enabled", defaultValue = "true") boolean enabled,
            @ConfigProperty(name = "app.rate-limit.max", defaultValue = "120") int max,
            @ConfigProperty(name = "app.rate-limit.window-seconds", defaultValue = "60") int windowSeconds) {
        this.enabled = enabled;
        this.limiter = new RateLimiter(max, Duration.ofSeconds(windowSeconds));
    }

    @Override
    public void filter(ContainerRequestContext requestContext) {
        String path = requestContext.getUriInfo().getPath();
        // Health probes are the platform asking whether we are alive; a 429
        // there reads as an outage.
        if (!enabled || path.startsWith("q/health") || path.equals("healthz")) {
            return;
        }
        RateLimiter.Decision decision = limiter.check(clientKey(requestContext));
        requestContext.setProperty(REMAINING_PROPERTY, decision.remaining());
        if (!decision.allowed()) {
            requestContext.abortWith(Response.status(Response.Status.TOO_MANY_REQUESTS)
                    .header("Retry-After", decision.retryAfterSeconds())
                    .type(MediaType.APPLICATION_JSON)
                    .entity("{\"error\":\"rate limit exceeded\"}")
                    .build());
        }
    }

    @Override
    public void filter(ContainerRequestContext requestContext, ContainerResponseContext responseContext) {
        Object remaining = requestContext.getProperty(REMAINING_PROPERTY);
        if (remaining != null) {
            responseContext.getHeaders().putSingle("X-RateLimit-Remaining", String.valueOf(remaining));
        }
    }

    private String clientKey(ContainerRequestContext ctx) {
        String forwarded = ctx.getHeaderString("X-Forwarded-For");
        if (forwarded != null && !forwarded.isBlank()) {
            return forwarded.split(",")[0].trim();
        }
        String realIp = ctx.getHeaderString("X-Real-IP");
        return realIp != null && !realIp.isBlank() ? realIp : "unknown";
    }
}
