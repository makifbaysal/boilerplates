package ai.tasktrooper.boilerplate.quarkus.platform.logging;

import java.util.UUID;

import jakarta.inject.Singleton;

import org.jboss.logging.Logger;
import org.jboss.logging.MDC;

import io.grpc.ForwardingServerCall.SimpleForwardingServerCall;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.Status;
import io.quarkus.grpc.GlobalInterceptor;

/** gRPC equivalent of RequestLoggingFilter — one structured log line
 * per call, MDC request_id for the duration of the call. */
@Singleton
@GlobalInterceptor
public class GrpcLoggingInterceptor implements ServerInterceptor {

    private static final Logger LOG = Logger.getLogger("grpc_request");

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        String requestId = UUID.randomUUID().toString();
        MDC.put("request_id", requestId);
        long start = System.currentTimeMillis();

        ServerCall<ReqT, RespT> wrapped = new SimpleForwardingServerCall<>(call) {
            @Override
            public void close(Status status, Metadata trailers) {
                LOG.infof("method=%s status=%s duration_ms=%d",
                        call.getMethodDescriptor().getFullMethodName(), status.getCode(),
                        System.currentTimeMillis() - start);
                MDC.remove("request_id");
                super.close(status, trailers);
            }
        };
        return next.startCall(wrapped, headers);
    }
}
