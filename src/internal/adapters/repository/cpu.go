package repository

import (
	"bufio"
	"errors"
	"os"
	"pi-monitor-api/internal/core/domain"
	"strconv"
	"strings"
)

/* ************************************* MOCKING SCAFFOLDING ************************************* */

type FileReader interface {
	Open(name string) (*os.File, error)
}

type RealFileReader struct{}

func (r *RealFileReader) Open(name string) (*os.File, error) {
	return os.Open(name)
}

/* ******************************************** CPU ******************************************** */

type CPURepository struct {
	fileReader FileReader
}

func NewCPURepository(fr FileReader) *CPURepository {
	return &CPURepository{fileReader: fr}
}

func (r *CPURepository) GetCPULoad() (domain.CPU, error) {
	file, err := r.fileReader.Open("/proc/loadavg")
	if err != nil {
		return domain.CPU{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return domain.CPU{}, scanner.Err()
	}

	parts := strings.Fields(scanner.Text())
	if len(parts) != 5 {
		// Handle case where file format is unexpected
		return domain.CPU{}, errors.New("unexpected file format")
	}

	oneMin, _ := strconv.ParseFloat(parts[0], 64)
	fiveMin, _ := strconv.ParseFloat(parts[1], 64)
	fifteenMin, _ := strconv.ParseFloat(parts[2], 64)

	return domain.CPU{
		LoadAvg1Min:  oneMin,
		LoadAvg5Min:  fiveMin,
		LoadAvg15Min: fifteenMin,
	}, nil
}
