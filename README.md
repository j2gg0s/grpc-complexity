# go-complexity

gRPC's rate limiter with reuqest complexity.

## Usage

### Install protoc plugin
`go get github.com/j2gg0s/grpc-complexity/cmd/protoc-gen-go-complexity`

### Generate Code
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
cserver, err := complexity.New(
    complexity.WithGlobalEvery(time.Second, 3),
    complexity.WithMaxWait(5*time.Second),
)
if err != nil {
    log.Fatalf("new complexity server: %v", err)
}

pb.RegisterGreeterComplexityServer(
    cserver,
    &pb.DefaultGreeterComplexityServer{},
)
s := grpc.NewServer(grpc.ChainUnaryInterceptor(cserver.UnaryServerInterceptor()))
```

See detail in `example/helloworld`.
