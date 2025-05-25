package main

import (
	"context"
	"log"
	"net"
	"time"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/discovery"
	"github.com/shimkek/omd-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr    = common.EnvGetString("GRPC_ADDR", "localhost:2000")
	consulAddr  = common.EnvGetString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		log.Fatal("failed to create registry:", err)
	}

	ctx := context.Background()
	instanceID := discovery.GenreateInstanceID(serviceName)
	if err := registry.RegisterService(ctx, instanceID, serviceName, "localhost", 2000); err != nil {
		log.Fatal("failed to register ordera service:", err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID); err != nil {
				log.Printf("Health check failed for service %s: %v", instanceID, err)
			}

			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.DeregisterService(ctx, instanceID)

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
