package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor is the gRPC equivalent of FiberMiddleware: one
// structured line per call, request-scoped logger attached to ctx.
func UnaryServerInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		reqLog := log.With("request_id", uuid.NewString())
		ctx = WithContext(ctx, reqLog)

		resp, err := handler(ctx, req)

		attrs := []any{
			"method", info.FullMethod,
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if err != nil {
			attrs = append(attrs, "error", err.Error())
			reqLog.Error("grpc_request", attrs...)
		} else {
			reqLog.Info("grpc_request", attrs...)
		}
		return resp, err
	}
}
