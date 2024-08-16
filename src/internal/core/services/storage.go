package services

import (
	"pi-monitor-api/internal/core/domain"
	"pi-monitor-api/internal/core/ports"
)

// StorageService provides business logic related to storage operations.
// Acts as a middleman between the core domain model (Storage) and the outside
type StorageService struct {
	storagePort ports.StoragePort
}

// Service constructor
func NewStorageService(storagePort ports.StoragePort) *StorageService {
	return &StorageService{storagePort: storagePort}
}

// Business logic to get the storage devices
func (s *StorageService) GetDevices() ([]domain.Device, error) {
	return s.storagePort.GetDevices()
}
