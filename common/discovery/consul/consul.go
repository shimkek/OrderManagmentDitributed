package consul

import (
	"context"
	"fmt"
	"log"

	consul "github.com/hashicorp/consul/api"
)

type Registry struct {
	Client *consul.Client
}

func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{
		Client: client,
	}, nil
}

func (r *Registry) RegisterService(ctx context.Context, instanceID, serviceName, host string, port int) error {
	log.Printf("Registering service %s with ID %s at %s:%d", serviceName, instanceID, host, port)

	agent := r.Client.Agent()

	service := &consul.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Address: host,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err := agent.ServiceRegister(service)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) DeregisterService(ctx context.Context, instanceID string) error {
	log.Printf("Deregistering service with ID %s", instanceID)

	return r.Client.Agent().CheckDeregister(instanceID)
}
func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.Client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	instances := make([]string, 0)
	for _, entry := range entries {
		address := fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
		instances = append(instances, address)
	}
	return instances, nil
}
func (r *Registry) HealthCheck(instanceID string) error {
	return r.Client.Agent().UpdateTTL(instanceID, "online", consul.HealthPassing)
}
