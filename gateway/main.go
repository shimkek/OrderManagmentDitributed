package main

import (
	"context"
	"log"
	"net/http"
	"time"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/discovery"
	"github.com/shimkek/omd-common/discovery/consul"
	"github.com/shimkek/omd-gateway/gateway" // Import the gateway package
)

var (
	httpAddr   = common.EnvGetString("HTTP_ADDR", ":8080")
	consulAddr = common.EnvGetString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		log.Fatal("failed to create registry:", err)
	}

	ctx := context.Background()
	instanceID := discovery.GenreateInstanceID("gateway")
	if err := registry.RegisterService(ctx, instanceID, "gateway", "localhost", 8080); err != nil {
		log.Fatal("failed to register gateway service:", err)
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

	mux := http.NewServeMux()
	ordersGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGateway)

	if err := handler.registerRoutes(mux); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server starting on %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal(err)
	}
}
