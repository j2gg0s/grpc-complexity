package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/j2gg0s/grpc-complexity/example/helloworld/helloworld"
	"google.golang.org/grpc"

	"github.com/j2gg0s/grpc-complexity/complexity"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

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
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
