// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package helloworld

import (
	context "context"
	grpc "google.golang.org/grpc"
)

type GreeterComplexity interface {
	SayHello(ctx context.Context, in *HelloRequest, out *HelloReply, err error, opts ...grpc.CallOption) int
}
