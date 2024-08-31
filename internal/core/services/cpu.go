package services

import (
	"pi-monitor-api/internal/core/domain"
	"pi-monitor-api/internal/core/ports"
)

// CPUService provides business logic related to CPU operations.
// Acts as a middleman between the core domain model (CPU) and the outside
type CPUService struct {
	cpuPort ports.CPUPort
}

// Service constructor
func NewCPUService(cpuPort ports.CPUPort) *CPUService {
	return &CPUService{cpuPort: cpuPort}
}

// Business logic to get the CPU load, N functions from here
func (s *CPUService) GetCPULoad() (domain.CPU, error) {
	return s.cpuPort.GetCPULoad()
}
