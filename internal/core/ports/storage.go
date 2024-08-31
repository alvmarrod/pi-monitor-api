package ports

// StoragePort defines the interface for interacting with storage-related
// operations.

import "github.com/alvmarrod/pi-monitor-api/internal/core/domain"

type StoragePort interface {
	GetDevices() ([]domain.Device, error)
}
