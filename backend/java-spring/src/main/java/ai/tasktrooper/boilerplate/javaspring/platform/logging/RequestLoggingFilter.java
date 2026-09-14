package ai.tasktrooper.boilerplate.javaspring.platform.logging;

import java.io.IOException;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.stereotype.Component;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpFilter;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

/**
 * REST equivalent of backend/go-fiber's FiberMiddleware: one structured
 * log line per request, with a request_id in MDC for the duration of
 * the request so anything logged downstream picks it up automatically
 * (SLF4J/Logback's MDC is the request-scoped-logger mechanism here,
 * where go-fiber attaches a logger to context.Context instead).
 */
@Component
public class RequestLoggingFilter extends HttpFilter {

    private static final Logger log = LoggerFactory.getLogger("http_request");
    private static final String REQUEST_ID_HEADER = "X-Request-ID";

    @Override
    protected void doFilter(HttpServletRequest req, HttpServletResponse res, FilterChain chain)
            throws IOException, ServletException {
        String requestId = req.getHeader(REQUEST_ID_HEADER);
        if (requestId == null || requestId.isBlank()) {
            requestId = UUID.randomUUID().toString();
        }
        res.setHeader(REQUEST_ID_HEADER, requestId);
        MDC.put("request_id", requestId);

        long start = System.currentTimeMillis();
        try {
            chain.doFilter(req, res);
        } finally {
            log.info("method={} path={} status={} duration_ms={}",
                    req.getMethod(), req.getRequestURI(), res.getStatus(), System.currentTimeMillis() - start);
            MDC.remove("request_id");
        }
    }
}
