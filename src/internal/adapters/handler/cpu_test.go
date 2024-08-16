package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pi-monitor-api/internal/adapters/handler"
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCPUPort struct {
	mock.Mock
}

func (m *MockCPUPort) GetCPULoad() (domain.CPU, error) {
	args := m.Called()
	return args.Get(0).(domain.CPU), args.Error(1)
}

func TestGetCPULoad_Success(t *testing.T) {

	mockCPUPort := new(MockCPUPort)
	cpuData := domain.CPU{
		LoadAvg1Min:  0.10,
		LoadAvg5Min:  0.15,
		LoadAvg15Min: 0.20,
	}
	mockCPUPort.On("GetCPULoad").Return(cpuData, nil)

	cpuHandler := handler.NewCPUHandler(mockCPUPort)

	req, err := http.NewRequest("GET", "/cpu/load", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	cpuHandler.GetCPULoad(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseCPU domain.CPU
	err = json.NewDecoder(rr.Body).Decode(&responseCPU)
	assert.NoError(t, err)

	assert.Equal(t, cpuData, responseCPU)
	mockCPUPort.AssertExpectations(t)
}

func TestGetCPULoad_Error(t *testing.T) {

	mockCPUPort := new(MockCPUPort)
	mockCPUPort.On("GetCPULoad").Return(domain.CPU{}, assert.AnError)

	cpuHandler := handler.NewCPUHandler(mockCPUPort)

	req, err := http.NewRequest("GET", "/cpu/load", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	cpuHandler.GetCPULoad(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to retrieve CPU load")
	mockCPUPort.AssertExpectations(t)
}
