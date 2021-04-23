package complexity

import (
	"context"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var counter *prometheus.CounterVec = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_complexity",
		Help: "Sum of grpc's complexity",
	},
	[]string{"service", "method", "token"},
)

type ServiceRegistrar interface {
	RegisterService(*ServiceDesc, interface{})
}

type Server struct {
	logger   Logger
	services map[string]*serviceInfo

	limiters      map[string]*rate.Limiter
	maxWait       time.Duration
	globalLimiter *rate.Limiter

	mu sync.Mutex

	enableMetric bool
	registerer   prometheus.Registerer
	counter      *prometheus.CounterVec
}

func New(opts ...Option) (*Server, error) {
	s := &Server{
		mu: sync.Mutex{},

		logger: WrapGrpcLogger(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)),

		maxWait:       time.Second,
		globalLimiter: rate.NewLimiter(rate.Inf, 1<<31),

		services: map[string]*serviceInfo{},
		limiters: map[string]*rate.Limiter{},

		enableMetric: true,
		registerer:   prometheus.DefaultRegisterer,
		counter:      counter,
	}
	for _, opt := range opts {
		opt(s)
	}

	if s.enableMetric && s.registerer != nil && s.counter != nil {
		if err := s.registerer.Register(s.counter); err != nil {
			if _, ok := err.(*prometheus.AlreadyRegisteredError); !ok {
				return nil, err
			}
		}
	}

	return s, nil
}

func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {
	if ss != nil {
		ht := reflect.TypeOf(sd.HandlerType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			s.logger.Fatalf(context.Background(), "complexity: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
		}
	}
	s.register(sd, ss)
}

func (s *Server) register(sd *ServiceDesc, ss interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.services[sd.ServiceName]; ok {
		s.logger.Fatalf(context.Background(), "complexity: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	}
	info := &serviceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*MethodDesc),
		mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}
	s.services[sd.ServiceName] = info
}

func (s *Server) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var method *MethodDesc
		var service *serviceInfo
		var sName, mName string
		if d := strings.Split(info.FullMethod, "/"); len(d) != 3 {
			s.logger.Warningf(ctx, "complexity: invalid method %s", info.FullMethod)
		} else {
			sName, mName = d[1], d[2]
			if svc, ok := s.services[sName]; ok {
				service = svc
				if md, ok := svc.methods[mName]; ok {
					method = md
				}
			}
		}

		if method == nil {
			s.logger.Infof(ctx, "complexity: unknown method %s", info.FullMethod)
			return handler(ctx, req)
		}

		if weights := method.ComplexityHandler(service.serviceImpl, ctx, req); len(weights) > 0 {
			for token, weight := range weights {
				limiter, ok := s.limiters[token]
				if !ok {
					limiter = s.globalLimiter
				}
				nctx, cancel := context.WithTimeout(ctx, s.maxWait)
				defer cancel()
				if err := limiter.WaitN(nctx, weight); err != nil {
					s.logger.Infof(ctx, "complexity: limiter[%s] wait %d: %v", token, weight, err)
					return nil, err
				}
				s.counter.WithLabelValues(sName, mName, token).Add(float64(weight))
			}
		}

		resp, err := handler(ctx, req)
		return resp, err
	}
}

// serviceInfo wraps information about a service. It is very similar to
// ServiceDesc and is constructed from it for internal purposes.
type serviceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*MethodDesc
	mdata       interface{}
}

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct {
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     []MethodDesc
	Metadata    interface{}
}

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName        string
	ComplexityHandler func(srv interface{}, ctx context.Context, req interface{}) map[string]int
}
