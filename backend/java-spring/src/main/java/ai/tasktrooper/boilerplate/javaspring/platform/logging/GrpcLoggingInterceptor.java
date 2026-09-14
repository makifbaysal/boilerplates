package ai.tasktrooper.boilerplate.javaspring.platform.logging;

import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;

import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.ForwardingServerCall.SimpleForwardingServerCall;
import io.grpc.Status;

import net.devh.boot.grpc.server.interceptor.GrpcGlobalServerInterceptor;

/**
 * gRPC equivalent of RequestLoggingFilter — one structured log line per
 * call, MDC request_id for the duration of the call. The go-fiber
 * sibling is logger.UnaryServerInterceptor.
 */
@GrpcGlobalServerInterceptor
public class GrpcLoggingInterceptor implements ServerInterceptor {

    private static final Logger log = LoggerFactory.getLogger("grpc_request");

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        String requestId = UUID.randomUUID().toString();
        MDC.put("request_id", requestId);
        long start = System.currentTimeMillis();

        ServerCall<ReqT, RespT> wrapped = new SimpleForwardingServerCall<>(call) {
            @Override
            public void close(Status status, Metadata trailers) {
                log.info("method={} status={} duration_ms={}",
                        call.getMethodDescriptor().getFullMethodName(), status.getCode(),
                        System.currentTimeMillis() - start);
                MDC.remove("request_id");
                super.close(status, trailers);
            }
        };
        return next.startCall(wrapped, headers);
    }
}
