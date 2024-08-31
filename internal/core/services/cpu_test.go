package services

import (
	"errors"
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCPUPort struct {
	mockResult domain.CPU
	mockError  error
}

func (m *mockCPUPort) GetCPULoad() (domain.CPU, error) {
	return m.mockResult, m.mockError
}

func TestGetCPUInfoValues(t *testing.T) {

	mockPort := &mockCPUPort{
		mockResult: domain.CPU{
			LoadAvg1Min:  1.5,
			LoadAvg5Min:  0.9,
			LoadAvg15Min: 0.6,
		},
		mockError: nil,
	}

	svc := NewCPUService(mockPort)

	result, err := svc.GetCPULoad()

	assert.NoError(t, err)

	assert.GreaterOrEqual(t, result.LoadAvg1Min, 0.0, "LoadAvg1Min should be >= 0")
	assert.GreaterOrEqual(t, result.LoadAvg5Min, 0.0, "LoadAvg5Min should be >= 0")
	assert.GreaterOrEqual(t, result.LoadAvg15Min, 0.0, "LoadAvg15Min should be >= 0")
}

func TestGetCPUInfoSimulateError(t *testing.T) {

	mockPort := &mockCPUPort{
		mockResult: domain.CPU{},
		mockError:  errors.New("unable to read CPU load"),
	}

	svc := NewCPUService(mockPort)

	_, err := svc.GetCPULoad()

	assert.Error(t, err)
	assert.Equal(t, "unable to read CPU load", err.Error())
}
