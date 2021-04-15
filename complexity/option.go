package complexity

import (
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/grpclog"
)

type Option func(*Server)

// WithMaxWait
// The max duration to wait, when qps reached the upper limit.
func WithMaxWait(maxWait time.Duration) Option {
	return func(s *Server) { s.maxWait = maxWait }
}

// WithGlobalLimiter
// The default limiter for any unknown token
func WithGlobalLimiter(limiter *rate.Limiter) Option {
	return func(s *Server) { s.globalLimiter = limiter }
}

// WithGlobalEvery
// refer: AddEvery
func WithGlobalEvery(d time.Duration, b int) Option {
	return func(s *Server) { s.globalLimiter = rate.NewLimiter(rate.Every(d), b) }
}

func WithLogger(logger Logger) Option {
	return func(s *Server) { s.logger = logger }
}

func WithGrpcLogger(logger grpclog.LoggerV2) Option {
	return func(s *Server) { s.logger = WrapGrpcLogger(logger) }
}

func AddLimiter(token string, limiter *rate.Limiter) Option {
	return func(s *Server) { s.limiters[token] = limiter }
}

// AddEvery
// Every d allow one request with 1 complexity
// Every request's complexity should less b
func AddEvery(token string, d time.Duration, b int) Option {
	return func(s *Server) { s.limiters[token] = rate.NewLimiter(rate.Every(d), b) }
}
