package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pi-monitor-api/internal/adapters/handler"
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) GetDevices() ([]domain.Device, error) {
	args := m.Called()
	return args.Get(0).([]domain.Device), args.Error(1)
}

func TestGetStorageInfo(t *testing.T) {
	mockService := new(MockStorageService)
	handler := handler.NewStorageHandler(mockService)

	mockDevices := []domain.Device{
		{
			Name: "sda",
			Partitions: map[string]domain.Partition{
				"/": {
					Name:       "sda1",
					MountPoint: "/",
					Filesystem: "ext4",
					Total:      1000000000,
					Used:       500000000,
					Free:       500000000,
				},
			},
		},
	}

	mockService.On("GetDevices").Return(mockDevices, nil)

	req, err := http.NewRequest("GET", "/storage", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetStorageInfo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse, _ := json.Marshal(mockDevices)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())

	mockService.AssertExpectations(t)
}

func TestGetStorageInfo_Error(t *testing.T) {

	mockService := new(MockStorageService)
	mockService.On("GetDevices").Return([]domain.Device{}, errors.New("some error"))

	handler := handler.NewStorageHandler(mockService)

	req, err := http.NewRequest("GET", "/storage", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler.GetStorageInfo(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	expectedBody := "Failed to retrieve storage info\n"
	assert.Equal(t, expectedBody, rr.Body.String())
	mockService.AssertExpectations(t)
}
