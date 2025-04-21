package interceptors

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AccessInterceptor struct {
	logger *zap.SugaredLogger
}

func NewAccessInterceptor(logger *zap.SugaredLogger) *AccessInterceptor {
	return &AccessInterceptor{logger: logger}
}

func (i *AccessInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctxLogger := i.logger.With(zap.String("method", info.FullMethod))
		newCtx := logger.LoggerToContext(ctx, ctxLogger)

		startTime := time.Now()

		ctxLogger.Infow("gRPC request received",
			"method", info.FullMethod,
			"request", req,
		)

		resp, err := handler(newCtx, req)

		duration := time.Since(startTime)

		if err != nil {
			st, _ := status.FromError(err)
			ctxLogger.Errorw("gRPC request failed",
				"method", info.FullMethod,
				"code", st.Code(),
				"error", err,
				"duration", duration,
			)
		} else {
			ctxLogger.Infow("gRPC request completed",
				"method", info.FullMethod,
				"duration", duration,
			)
		}

		return resp, err
	}
}
