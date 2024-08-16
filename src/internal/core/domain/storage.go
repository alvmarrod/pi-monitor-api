package domain

type Partition struct {
	Name       string
	MountPoint string
	Filesystem string
	Total      uint64
	Used       uint64
	Free       uint64
}

type Device struct {
	Name       string
	Partitions map[string]Partition // Keyed by mount point
}
