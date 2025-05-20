package main

import (
	"log"
	"net/http"

	common "github.com/shimkek/omd-common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/shimkek/omd-common/api"
)

var (
	httpAddr         = common.EnvGetString("HTTP_ADDR", ":8080")
	orderServiceAddr = "localhost:2000"
)

func main() {
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to order service at %s", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)

	if err := handler.registerRoutes(mux); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server starting on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(err)
	}
}
