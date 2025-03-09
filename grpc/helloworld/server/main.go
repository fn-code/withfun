package main

import (
	"context"
	"github.com/fn-code/withfun/grpc/helloworld/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hella " + in.GetName()}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":40023")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server running at %v:", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
