package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/broker"
	"github.com/shimkek/omd-common/discovery"
	"github.com/shimkek/omd-common/discovery/consul"
	"github.com/shimkek/omd-payments/gateway"
	StripeProcessor "github.com/shimkek/omd-payments/processor/stripe"
	"github.com/stripe/stripe-go/v82"
	"google.golang.org/grpc"
)

var (
	serviceName  = "payments"
	grpcAddr     = common.EnvGetString("GRPC_ADDR", "localhost:2001")
	consulAddr   = common.EnvGetString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvGetString("RABBITMQ_USER", "guest")
	amqpPassword = common.EnvGetString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = common.EnvGetString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvGetString("RABBITMQ_PORT", "5672")
	stripeKey    = common.EnvGetString("STRIPE_KEY", "")
	httpAddr     = common.EnvGetString("PAYMENTS_HTTP_ADDR", "localhost:8081")
	stripeSecret = common.EnvGetString("STRIPE_WEBHOOK_SECRET", "")
	jaegerAddr   = common.EnvGetString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	ctx := context.Background()
	if err := common.SetGlobalTracer(ctx, serviceName, jaegerAddr); err != nil {
		log.Fatal("failed to set global tracer")
	}

	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		log.Fatal("failed to create registry:", err)
	}

	instanceID := discovery.GenreateInstanceID(serviceName)
	if err := registry.RegisterService(ctx, instanceID, serviceName, "localhost", 2000); err != nil {
		log.Fatal("failed to register payments service:", err)
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

	stripe.Key = stripeKey

	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	gateway := gateway.NewGRPCGateway(registry)

	stripeProcessor := StripeProcessor.NewProcessor()
	service := NewService(stripeProcessor, gateway)
	svcWithTelemetry := NewTelemetryMiddleware(service)

	amqpConsumer := NewConsumer(svcWithTelemetry)
	go amqpConsumer.Listen(ch)

	//http server
	mux := http.NewServeMux()

	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Starting http server on %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("failed to serve http server")
		}
	}()
	//grpc server
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Printf("Starting gRPC server on %s", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
