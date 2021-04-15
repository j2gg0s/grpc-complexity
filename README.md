# go-complexity

gRPC's rate limiter with reuqest complexity.

## Usage

### protoc
```
protoc \
    --proto_path=. \
    --go_out=${GOPATH}/src \
    --go-grpc_out=${GOPATH}/src \
    --go-complexity_out=${GOPATH}/src \
    example/helloworld/helloworld/helloworld.proto
```

### Add interceptor
```
cserver := complexity.New(
    complexity.WithGlobalEvery(time.Second, 3),
    complexity.WithMaxWait(5*time.Second),
)
pb.RegisterGreeterComplexityServer(
    cserver,
    &pb.DefaultGreeterComplexityServer{},
)
s := grpc.NewServer(grpc.ChainUnaryInterceptor(cserver.UnaryServerInterceptor()))
```

See detail in `example/helloworld`.
