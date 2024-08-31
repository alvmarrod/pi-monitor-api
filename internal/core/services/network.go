package services

import (
	"pi-monitor-api/internal/core/domain"
	"pi-monitor-api/internal/core/ports"
)

// NetworkService provides business logic related to network operations.
// Acts as a middleman between the core domain model (Network) and the outside
type NetworkService struct {
	networkPort ports.NetworkPort
}

// Service constructor
func NewNetworkService(networkPort ports.NetworkPort) *NetworkService {
	return &NetworkService{networkPort: networkPort}
}

// Business logic to get the network interfaces, N functions from here
func (s *NetworkService) GetNetworkInterfaces() ([]domain.NetworkInterface, error) {
	return s.networkPort.GetNetworkInterfaces()
}
