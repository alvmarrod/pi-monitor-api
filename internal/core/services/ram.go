package services

import (
	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"
	"github.com/alvmarrod/pi-monitor-api/internal/core/ports"
)

// RAMService provides business logic related to RAM operations.
// Acts as a middleman between the core domain model (RAM) and the outside
type RAMService struct {
	ramPort ports.RAMPort
}

// Service constructor
func NewRAMService(ramPort ports.RAMPort) *RAMService {
	return &RAMService{ramPort: ramPort}
}

// Business logic to get the RAM stats
func (s *RAMService) GetRAMStats() (domain.RAM, error) {
	return s.ramPort.GetRAMStats()
}
