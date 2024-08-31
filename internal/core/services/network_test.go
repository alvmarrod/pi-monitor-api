package services

import (
	"errors"
	"pi-monitor-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockNetworkPort struct {
	mockResult []domain.NetworkInterface
	mockError  error
}

func (m *mockNetworkPort) GetNetworkInterfaces() ([]domain.NetworkInterface, error) {
	return m.mockResult, m.mockError
}

func TestGetNetworkInterfacesValues(t *testing.T) {

	mockPort := &mockNetworkPort{
		mockResult: []domain.NetworkInterface{
			{
				InterfaceName: "eth0",
				BitRate:       1000,
				Rx: domain.NetworkStats{
					Packets: 100,
					Bytes:   10000,
					Errors:  0,
					Drops:   0,
				},
				Tx: domain.NetworkStats{
					Packets: 200,
					Bytes:   20000,
					Errors:  0,
					Drops:   0,
				},
			},
			{
				InterfaceName: "wlan0",
				BitRate:       150,
				Rx: domain.NetworkStats{
					Packets: 50,
					Bytes:   5000,
					Errors:  1,
					Drops:   1,
				},
				Tx: domain.NetworkStats{
					Packets: 70,
					Bytes:   7000,
					Errors:  0,
					Drops:   0,
				},
			},
		},
		mockError: nil,
	}
	svc := NewNetworkService(mockPort)

	result, err := svc.GetNetworkInterfaces()

	assert.NoError(t, err)
	assert.Len(t, result, 2, "Expected two network interfaces")

	eth0 := result[0]
	assert.Equal(t, "eth0", eth0.InterfaceName)
	assert.GreaterOrEqual(t, eth0.BitRate, uint64(0), "BitRate should be >= 0")
	assert.GreaterOrEqual(t, eth0.Rx.Packets, uint64(0), "RxPackets should be >= 0")
	assert.GreaterOrEqual(t, eth0.Rx.Bytes, uint64(0), "RxBytes should be >= 0")
	assert.GreaterOrEqual(t, eth0.Tx.Packets, uint64(0), "TxPackets should be >= 0")
	assert.GreaterOrEqual(t, eth0.Tx.Bytes, uint64(0), "TxBytes should be >= 0")

	wlan0 := result[1]
	assert.Equal(t, "wlan0", wlan0.InterfaceName)
	assert.GreaterOrEqual(t, wlan0.BitRate, uint64(0), "BitRate should be >= 0")
	assert.GreaterOrEqual(t, wlan0.Rx.Packets, uint64(0), "RxPackets should be >= 0")
	assert.GreaterOrEqual(t, wlan0.Rx.Bytes, uint64(0), "RxBytes should be >= 0")
	assert.GreaterOrEqual(t, wlan0.Tx.Packets, uint64(0), "TxPackets should be >= 0")
	assert.GreaterOrEqual(t, wlan0.Tx.Bytes, uint64(0), "TxBytes should be >= 0")
}

func TestGetNetworkInfoSimulateError(t *testing.T) {

	mockPort := &mockNetworkPort{
		mockResult: []domain.NetworkInterface{},
		mockError:  errors.New("unable to read Network Interfaces info"),
	}

	svc := NewNetworkService(mockPort)

	_, err := svc.GetNetworkInterfaces()

	assert.Error(t, err)
	assert.Equal(t, "unable to read Network Interfaces info", err.Error())
}
