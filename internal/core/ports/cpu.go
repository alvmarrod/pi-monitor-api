package ports

// CPUPort defines the interface for interacting with CPU-related
// operations.

import "github.com/alvmarrod/pi-monitor-api/internal/core/domain"

type CPUPort interface {
	GetCPULoad() (domain.CPU, error)
}
