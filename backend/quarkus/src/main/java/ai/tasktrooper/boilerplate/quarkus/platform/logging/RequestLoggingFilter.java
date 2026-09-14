package ai.tasktrooper.boilerplate.quarkus.platform.logging;

import java.util.UUID;

import jakarta.annotation.Priority;
import jakarta.ws.rs.Priorities;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.container.ContainerRequestFilter;
import jakarta.ws.rs.container.ContainerResponseContext;
import jakarta.ws.rs.container.ContainerResponseFilter;
import jakarta.ws.rs.ext.Provider;

import org.jboss.logging.Logger;
import org.jboss.logging.MDC;

/**
 * REST equivalent of backend/go-fiber's FiberMiddleware / java-spring's
 * RequestLoggingFilter: one structured log line per request, with a
 * request_id in MDC for the duration of the request.
 */
@Provider
@Priority(Priorities.HEADER_DECORATOR)
public class RequestLoggingFilter implements ContainerRequestFilter, ContainerResponseFilter {

    private static final Logger LOG = Logger.getLogger("http_request");
    private static final String REQUEST_ID_HEADER = "X-Request-ID";
    private static final String START_TIME_PROPERTY = "requestStartTime";

    @Override
    public void filter(ContainerRequestContext requestContext) {
        String requestId = requestContext.getHeaderString(REQUEST_ID_HEADER);
        if (requestId == null || requestId.isBlank()) {
            requestId = UUID.randomUUID().toString();
        }
        MDC.put("request_id", requestId);
        requestContext.setProperty("requestId", requestId);
        requestContext.setProperty(START_TIME_PROPERTY, System.currentTimeMillis());
    }

    @Override
    public void filter(ContainerRequestContext requestContext, ContainerResponseContext responseContext) {
        String requestId = (String) requestContext.getProperty("requestId");
        responseContext.getHeaders().add(REQUEST_ID_HEADER, requestId);

        Long start = (Long) requestContext.getProperty(START_TIME_PROPERTY);
        long durationMs = start == null ? -1 : System.currentTimeMillis() - start;

        LOG.infof("method=%s path=%s status=%d duration_ms=%d",
                requestContext.getMethod(), requestContext.getUriInfo().getPath(),
                responseContext.getStatus(), durationMs);
        MDC.remove("request_id");
    }
}
