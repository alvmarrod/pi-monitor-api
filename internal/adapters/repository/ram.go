package repository

import (
	"bufio"
	"errors"
	"strconv"
	"strings"

	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"
)

/* ************************************* MOCKING SCAFFOLDING ************************************* */

type RAMRepository struct {
	fileReader FileReader
}

func NewRAMRepository(fr FileReader) *RAMRepository {
	return &RAMRepository{fileReader: fr}
}

/* ******************************************** RAM ******************************************** */

func (r *RAMRepository) GetRAMStats() (domain.RAM, error) {
	file, err := r.fileReader.Open("/proc/meminfo")
	if err != nil {
		return domain.RAM{}, err
	}
	defer file.Close()

	stats := map[string]uint64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// log.Println("line:", line)
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return domain.RAM{}, errors.New("unexpected file format")
		}
		// Remove trailing colon
		key := fields[0][:len(fields[0])-1]
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return domain.RAM{}, errors.New("unexpected file format")
		}
		stats[key] = value
	}

	if len(stats) == 0 {
		return domain.RAM{}, errors.New("unexpected file format")
	}

	// Buffers and cache are not included in the used memory
	return domain.RAM{
		Total:     stats["MemTotal"] * 1024,
		Available: stats["MemAvailable"] * 1024,
		Free:      stats["MemFree"] * 1024,
		Used:      (stats["MemTotal"] - stats["MemFree"]) * 1024,
	}, nil
}
