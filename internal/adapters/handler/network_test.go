package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvmarrod/pi-monitor-api/internal/adapters/handler"
	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNetworkService struct {
	mock.Mock
}

func (m *MockNetworkService) GetNetworkInterfaces() ([]domain.NetworkInterface, error) {
	args := m.Called()
	return args.Get(0).([]domain.NetworkInterface), args.Error(1)
}

func TestGetNetworkInfo_Success(t *testing.T) {
	mockService := new(MockNetworkService)
	handler := handler.NewNetworkHandler(mockService)

	mockNetworkInfo := []domain.NetworkInterface{
		{
			InterfaceName: "eth0",
			BitRate:       1000000000, // 1 Gbps
			Rx: domain.NetworkStats{
				Packets: 1000,
				Bytes:   1000000,
				Errors:  0,
				Drops:   0,
			},
			Tx: domain.NetworkStats{
				Packets: 900,
				Bytes:   900000,
				Errors:  0,
				Drops:   0,
			},
		},
	}

	mockService.On("GetNetworkInterfaces").Return(mockNetworkInfo, nil)

	req, err := http.NewRequest("GET", "/network", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetNetworkInfo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse, _ := json.Marshal(mockNetworkInfo)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestGetNetworkInfo_Error(t *testing.T) {

	mockService := new(MockNetworkService)
	mockService.On("GetNetworkInterfaces").Return([]domain.NetworkInterface{}, errors.New("some error"))

	handler := handler.NewNetworkHandler(mockService)

	req, err := http.NewRequest("GET", "/network", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler.GetNetworkInfo(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expectedBody := "Failed to retrieve network info\n"
	assert.Equal(t, expectedBody, rr.Body.String())
	mockService.AssertExpectations(t)
}
