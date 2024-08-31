package repository

import (
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ******************************************** RAM TEST ******************************************** */

func TestGetRAMStats(t *testing.T) {

	mockData := `MemTotal:        7795036 kB
MemFree:         5805016 kB
MemAvailable:    6075188 kB
Buffers:           18872 kB
Cached:           427272 kB`
	mockFileReader := &MockFileReader{Data: mockData}

	repo := NewRAMRepository(mockFileReader)

	ram, err := repo.GetRAMStats()
	assert.NoError(t, err)
	assert.NotZero(t, ram.Total)
	assert.NotZero(t, ram.Available)
	assert.NotZero(t, ram.Free)
	assert.NotZero(t, ram.Used)
}

func TestGetRAMStats_FileError(t *testing.T) {

	mockData := "some incorrect file data"
	mockFileReader := &MockFileReader{Data: mockData}

	repo := NewRAMRepository(mockFileReader)

	ram, err := repo.GetRAMStats()
	assert.Error(t, err)
	assert.Equal(t, domain.RAM{}, ram)
}
