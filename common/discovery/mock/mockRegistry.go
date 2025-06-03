package mockRegistry

import "context"

type MockRegistry struct {
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{}
}
func (r *MockRegistry) RegisterService(ctx context.Context, instanceID, serverName, host string, port int) error {
	return nil
}
func (r *MockRegistry) DeregisterService(ctx context.Context, instanceID string) error {
	return nil
}
func (r *MockRegistry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	return nil, nil
}
func (r *MockRegistry) HealthCheck(instanceID string) error {
	return nil
}
