package repository

import (
	"os"
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ******************************************** MOCKING ******************************************** */

type MockFileReader struct {
	Data string
}

func (m *MockFileReader) Open(name string) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "mockfile")
	if err != nil {
		return nil, err
	}

	if _, err := tmpFile.WriteString(m.Data); err != nil {
		return nil, err
	}

	// Reset the read pointer to the beginning of the file
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

/* ******************************************** CPU TEST ******************************************** */

func TestGetCPULoad(t *testing.T) {

	mockData := "0.10 0.15 0.20 1/100 100/1000"
	mockFileReader := &MockFileReader{Data: mockData}

	repo := NewCPURepository(mockFileReader)

	cpu, err := repo.GetCPULoad()
	assert.NoError(t, err)
	assert.Equal(t, 0.10, cpu.LoadAvg1Min)
	assert.Equal(t, 0.15, cpu.LoadAvg5Min)
	assert.Equal(t, 0.20, cpu.LoadAvg15Min)
}

func TestGetCPULoad_FileError(t *testing.T) {

	mockData := "incorrect file data"
	mockFileReader := &MockFileReader{Data: mockData}

	repo := NewCPURepository(mockFileReader)

	cpu, err := repo.GetCPULoad()
	assert.Error(t, err, "unexpected file format")
	assert.Equal(t, domain.CPU{}, cpu)
}
