package services

import (
	"errors"
	"testing"

	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

type mockRAMPort struct {
	mockResult domain.RAM
	mockError  error
}

func (m *mockRAMPort) GetRAMStats() (domain.RAM, error) {
	return m.mockResult, m.mockError
}

func TestGetRAMInfoValues(t *testing.T) {

	mockPort := &mockRAMPort{
		mockResult: domain.RAM{
			Total:     8000,
			Available: 4000,
			Free:      2000,
			Used:      4000,
		},
		mockError: nil,
	}

	svc := NewRAMService(mockPort)

	result, err := svc.GetRAMStats()

	assert.NoError(t, err)

	assert.GreaterOrEqual(t, result.Total, uint64(0), "Total RAM should be >= 0")
	assert.GreaterOrEqual(t, result.Available, uint64(0), "Available RAM should be >= 0")
	assert.GreaterOrEqual(t, result.Free, uint64(0), "Free RAM should be >= 0")
	assert.GreaterOrEqual(t, result.Used, uint64(0), "Used RAM should be >= 0")

	assert.LessOrEqual(t, result.Used+result.Available, result.Total, "Used + Available should be <= Total RAM")
}

func TestGetRAMInfoSimulateError(t *testing.T) {

	mockPort := &mockRAMPort{
		mockResult: domain.RAM{},
		mockError:  errors.New("unable to read RAM stats"),
	}

	svc := NewRAMService(mockPort)

	_, err := svc.GetRAMStats()

	assert.Error(t, err)
	assert.Equal(t, "unable to read RAM stats", err.Error())
}
