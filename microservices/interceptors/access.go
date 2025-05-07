package interceptors

import (
	"context"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	metrics "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AccessInterceptor struct {
	logger  *zap.SugaredLogger
	metrics *metrics.Metrics
}

func NewAccessInterceptor(logger *zap.SugaredLogger, metrics *metrics.Metrics) *AccessInterceptor {
	return &AccessInterceptor{logger: logger, metrics: metrics}
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
			i.metrics.GRPCTotalNumberOfRequests.WithLabelValues(info.FullMethod, strconv.Itoa(int(st.Code()))).Inc()
			i.metrics.GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())
		} else {
			ctxLogger.Infow("gRPC request completed",
				"method", info.FullMethod,
				"duration", duration,
			)
			i.metrics.GRPCTotalNumberOfRequests.WithLabelValues(info.FullMethod, strconv.Itoa(200)).Inc()
			i.metrics.GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())
		}

		return resp, err
	}
}
