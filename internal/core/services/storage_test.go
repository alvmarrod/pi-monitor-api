package services

import (
	"errors"
	"testing"

	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

type mockStoragePort struct {
	mockResult []domain.Device
	mockError  error
}

func (m *mockStoragePort) GetDevices() ([]domain.Device, error) {
	return m.mockResult, m.mockError
}

func TestGetStorageDevices(t *testing.T) {
	// Initialize the mock port and service
	mockPort := &mockStoragePort{
		mockResult: []domain.Device{
			{
				Name: "sda",
				Partitions: map[string]domain.Partition{
					"/": {
						Name:       "sda1",
						MountPoint: "/",
						Filesystem: "ext4",
						Total:      500000,
						Used:       250000,
						Free:       250000,
					},
					"/home": {
						Name:       "sda2",
						MountPoint: "/home",
						Filesystem: "ext4",
						Total:      1000000,
						Used:       400000,
						Free:       600000,
					},
				},
			},
			{
				Name: "sdb",
				Partitions: map[string]domain.Partition{
					"/mnt/data": {
						Name:       "sdb1",
						MountPoint: "/mnt/data",
						Filesystem: "ext4",
						Total:      2000000,
						Used:       500000,
						Free:       1500000,
					},
				},
			},
		},
		mockError: nil,
	}
	svc := NewStorageService(mockPort)

	// Call the service method
	result, err := svc.GetDevices()

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert that the returned devices are correct
	assert.Len(t, result, 2, "Expected two devices")

	// Check the first device
	sda := result[0]
	assert.Equal(t, "sda", sda.Name)
	assert.Len(t, sda.Partitions, 2, "Expected two partitions on sda")

	// Check the first partition of sda
	rootPartition := sda.Partitions["/"]
	assert.Equal(t, "sda1", rootPartition.Name)
	assert.Equal(t, "/", rootPartition.MountPoint)
	assert.Equal(t, "ext4", rootPartition.Filesystem)
	assert.GreaterOrEqual(t, rootPartition.Total, uint64(0), "Total should be >= 0")
	assert.GreaterOrEqual(t, rootPartition.Used, uint64(0), "Used should be >= 0")
	assert.GreaterOrEqual(t, rootPartition.Free, uint64(0), "Free should be >= 0")

	// Check the second device
	sdb := result[1]
	assert.Equal(t, "sdb", sdb.Name)
	assert.Len(t, sdb.Partitions, 1, "Expected one partition on sdb")

	// Check the first partition of sdb
	dataPartition := sdb.Partitions["/mnt/data"]
	assert.Equal(t, "sdb1", dataPartition.Name)
	assert.Equal(t, "/mnt/data", dataPartition.MountPoint)
	assert.Equal(t, "ext4", dataPartition.Filesystem)
	assert.GreaterOrEqual(t, dataPartition.Total, uint64(0), "Total should be >= 0")
	assert.GreaterOrEqual(t, dataPartition.Used, uint64(0), "Used should be >= 0")
	assert.GreaterOrEqual(t, dataPartition.Free, uint64(0), "Free should be >= 0")
}

func TestGetStorageInfoSimulateError(t *testing.T) {

	mockPort := &mockStoragePort{
		mockResult: []domain.Device{},
		mockError:  errors.New("unable to read Storage info"),
	}

	svc := NewStorageService(mockPort)

	_, err := svc.GetDevices()

	assert.Error(t, err)
	assert.Equal(t, "unable to read Storage info", err.Error())
}
