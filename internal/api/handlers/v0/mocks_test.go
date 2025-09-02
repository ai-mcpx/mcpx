package v0_test

import (
	"github.com/stretchr/testify/mock"
	apiv0 "github.com/modelcontextprotocol/registry/pkg/api/v0"
)

// MockRegistryService is a mock implementation of the RegistryService interface
type MockRegistryService struct {
	mock.Mock
}

func (m *MockRegistryService) List(cursor string, limit int) ([]apiv0.ServerJSON, string, error) {
	args := m.Called(cursor, limit)
	return args.Get(0).([]apiv0.ServerJSON), args.String(1), args.Error(2)
}

func (m *MockRegistryService) GetByID(id string) (*apiv0.ServerJSON, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiv0.ServerJSON), args.Error(1)
}

func (m *MockRegistryService) Publish(request apiv0.ServerJSON) (*apiv0.ServerJSON, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiv0.ServerJSON), args.Error(1)
}

func (m *MockRegistryService) EditServer(id string, request apiv0.ServerJSON) (*apiv0.ServerJSON, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiv0.ServerJSON), args.Error(1)
}

func (m *MockRegistryService) Update(id string, request apiv0.ServerJSON) (*apiv0.ServerJSON, error) {
	args := m.Called(id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apiv0.ServerJSON), args.Error(1)
}

func (m *MockRegistryService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
