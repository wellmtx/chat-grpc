package main

import (
	"chat/framework/pb"
	"chat/framework/servers"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	host, err := net.Listen("tcp", "localhost:5400")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	chatServer := servers.NewChatServer()

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	grpcServer.Serve(host)
}
