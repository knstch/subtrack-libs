package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"

	metrics "github.com/knstch/subtrack-libs/prometeus"
)

func WithTrackingRequests(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		start := time.Now()
		resp, err := next(ctx, request)

		method, _ := ctx.Value(httptransport.ContextKeyRequestMethod).(string)
		path, _ := ctx.Value(httptransport.ContextKeyRequestPath).(string)
		errLabel := "false"
		if err != nil {
			errLabel = "true"
		}

		metrics.RequestCount.With("method", method, "path", path, "error", errLabel).Add(1)
		metrics.RequestDuration.With("method", method).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
