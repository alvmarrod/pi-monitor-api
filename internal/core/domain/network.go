package domain

type NetworkStats struct {
	Packets uint64
	Bytes   uint64
	Errors  uint64
	Drops   uint64
}

type NetworkInterface struct {
	InterfaceName string
	BitRate       uint64
	Rx            NetworkStats
	Tx            NetworkStats
}
