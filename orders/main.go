package main

import (
	"context"
	"net"
	"time"

	common "github.com/shimkek/omd-common"
	"github.com/shimkek/omd-common/broker"
	"github.com/shimkek/omd-common/discovery"
	"github.com/shimkek/omd-common/discovery/consul"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName  = "orders"
	grpcAddr     = common.EnvGetString("GRPC_ADDR", "localhost:2000")
	consulAddr   = common.EnvGetString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvGetString("RABBITMQ_USER", "guest")
	amqpPassword = common.EnvGetString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = common.EnvGetString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvGetString("RABBITMQ_PORT", "5672")
	jaegerAddr   = common.EnvGetString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	ctx := context.Background()
	if err := common.SetGlobalTracer(ctx, serviceName, jaegerAddr); err != nil {
		logger.Fatal("failed to set global tracer")
	}

	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		logger.Fatal("failed to create registry:", zap.Error(err))
	}

	instanceID := discovery.GenreateInstanceID(serviceName)
	if err := registry.RegisterService(ctx, instanceID, serviceName, "localhost", 2000); err != nil {
		logger.Fatal("failed to register ordera service:", zap.Error(err))
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID); err != nil {
				logger.Error("Health check failed : %v", zap.Error(err))
			}

			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.DeregisterService(ctx, instanceID)

	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}
	defer l.Close()

	store := NewStore()
	service := NewService(store)
	svcWithTelemetry := NewTelemetryMiddleware(service)
	svcWithLogging := NewLoggingMiddleware(svcWithTelemetry)

	NewGrpcHandler(grpcServer, svcWithLogging, ch)

	consumer := NewConsumer(svcWithLogging)
	go consumer.Listen(ch)

	logger.Info("gRPC Server starting", zap.String("address", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		logger.Fatal(err.Error())
	}
}
