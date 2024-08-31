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

type MockRAMPort struct {
	mock.Mock
}

func (m *MockRAMPort) GetRAMStats() (domain.RAM, error) {
	args := m.Called()
	return args.Get(0).(domain.RAM), args.Error(1)
}

func TestGetRAMInfo_Success(t *testing.T) {

	mockRAMPort := new(MockRAMPort)
	ramData := domain.RAM{
		Total:     4096,
		Used:      2048,
		Free:      1024,
		Available: 1024,
	}
	mockRAMPort.On("GetRAMStats").Return(ramData, nil)

	ramHandler := handler.NewRAMHandler(mockRAMPort)

	req, err := http.NewRequest("GET", "/ram/info", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	ramHandler.GetRAMInfo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseRAM domain.RAM
	err = json.NewDecoder(rr.Body).Decode(&responseRAM)
	assert.NoError(t, err)

	assert.Equal(t, ramData, responseRAM)
	mockRAMPort.AssertExpectations(t)
}

func TestGetRAMInfo_Error(t *testing.T) {

	mockRAMPort := new(MockRAMPort)
	mockRAMPort.On("GetRAMStats").Return(domain.RAM{}, assert.AnError)

	ramHandler := handler.NewRAMHandler(mockRAMPort)

	req, err := http.NewRequest("GET", "/ram/info", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	ramHandler.GetRAMInfo(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to retrieve RAM info")
	mockRAMPort.AssertExpectations(t)
}
