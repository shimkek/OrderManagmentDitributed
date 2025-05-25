package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addr, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.NewClient(addr[rand.Intn(len(addr))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
