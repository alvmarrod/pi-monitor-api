package ports

// NetworkPort defines the interface for interacting with network-related
// operations.

import "github.com/alvmarrod/pi-monitor-api/internal/core/domain"

type NetworkPort interface {
	GetNetworkInterfaces() ([]domain.NetworkInterface, error)
}
