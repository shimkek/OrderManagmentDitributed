package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	RegisterService(ctx context.Context, instanceID, serverName, host string, port int) error
	DeregisterService(ctx context.Context, instanceID string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceID string) error
}

func GenreateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
