package main

import (
	"log"
	"net"

	common "github.com/shimkek/omd-common"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvGetString("GRPC_ADDR", "localhost:2000")
)

func main() {

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)

	NewGrpcHandler(grpcServer, service)

	log.Printf("gRPC Server starting on %s", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
