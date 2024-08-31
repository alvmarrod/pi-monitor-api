package ports

// RAMPort defines the interface for interacting with RAM-related
// operations.

import "pi-monitor-api/internal/core/domain"

type RAMPort interface {
	GetRAMStats() (domain.RAM, error)
}
