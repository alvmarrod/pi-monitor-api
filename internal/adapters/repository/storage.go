package repository

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"pi-monitor-api/internal/core/domain"
	"strconv"
	"strings"
)

/* ******************************************** AUX ******************************************** */

func isPartition(name string) bool {
	return strings.HasPrefix(name, "sd") ||
		strings.HasPrefix(name, "nvme") ||
		strings.HasPrefix(name, "hd")
}

func getDeviceName(partitionName string) string {
	// For NVME, namespaces are considered part of the partition name
	if strings.HasPrefix(partitionName, "nvme") {
		return partitionName[:5]
	} else if strings.HasPrefix(partitionName, "sd") {
		return partitionName[:3]
	} else if strings.HasPrefix(partitionName, "hd") {
		return partitionName[:3]
	} else {
		fmt.Printf("Unsupported device type for partition: %s\n", partitionName)
		return ""
	}
}

func groupDevices(partitions []domain.Partition) []domain.Device {
	devicesMap := make(map[string]domain.Device)

	for _, partition := range partitions {
		deviceName := getDeviceName(partition.Name)
		device, exists := devicesMap[deviceName]
		if !exists {
			device = domain.Device{Name: deviceName, Partitions: make(map[string]domain.Partition)}
		}
		device.Partitions[partition.Name] = partition
		devicesMap[deviceName] = device
	}

	var devices []domain.Device
	for _, device := range devicesMap {
		devices = append(devices, device)
	}
	return devices
}

func parseDfOutput(output []byte) map[string]domain.Partition {
	fsInfo := make(map[string]domain.Partition)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			fmt.Printf("Invalid df output line: %s\n", line)
			continue
		}
		device := strings.TrimPrefix(fields[0], "/dev/")
		total, _ := strconv.ParseUint(fields[1], 10, 64)
		used, _ := strconv.ParseUint(fields[2], 10, 64)
		free, _ := strconv.ParseUint(fields[3], 10, 64)
		filesystem := fields[0] // Filesystem type

		if isPartition(device) {
			fsInfo[device] = domain.Partition{
				Name:       device,
				Filesystem: filesystem,
				Total:      total,
				Used:       used,
				Free:       free,
			}
		}
	}

	return fsInfo
}

func pathProximityToRoot(path string) int {
	return len(strings.Split(path, "/"))
}

func newPathIsCloserToRoot(base string, newPath string) bool {
	oldProximity := pathProximityToRoot(base)
	newProximity := pathProximityToRoot(newPath)
	if newProximity == oldProximity {
		fmt.Println("Two paths are at the same level!")
	}
	return pathProximityToRoot(newPath) <= pathProximityToRoot(base)
}

/* ************************************* MOCKING SCAFFOLDING ************************************* */

/* ******************************************** STORAGE ******************************************** */

type StorageRepository struct {
	fileReader  FileReader
	toolChecker ToolInstalled
	cmdExec     CmdExecutor
}

func NewStorageRepository(fr FileReader, ti ToolInstalled, cmd CmdExecutor) *StorageRepository {
	return &StorageRepository{
		fileReader:  fr,
		toolChecker: ti,
		cmdExec:     cmd,
	}
}

func (r *StorageRepository) GetDevices() ([]domain.Device, error) {
	partitions, err := r.readPartitions()
	if err != nil {
		return []domain.Device{}, err
	}

	return groupDevices(partitions), nil
}

func (r *StorageRepository) readPartitions() ([]domain.Partition, error) {
	partitions, err := r.readBasicPartitions()
	if err != nil {
		return []domain.Partition{}, err
	}

	mountPoints, err := r.readMounts()
	if err != nil {
		return []domain.Partition{}, err
	}

	// Retrieve detailed partition info using df command
	fsInfo, err := r.getFilesystemInfo()
	if err != nil {
		return []domain.Partition{}, err
	}

	for i := range partitions {
		partition := &partitions[i]
		if mountPoint, exists := mountPoints[partition.Name]; exists {
			partition.MountPoint = mountPoint
		}
		if info, exists := fsInfo[partition.Name]; exists {
			partition.Filesystem = info.Filesystem
			partition.Total = info.Total
			partition.Used = info.Used
			partition.Free = info.Free
		}
	}

	return partitions, nil
}

func (r *StorageRepository) readBasicPartitions() ([]domain.Partition, error) {
	file, err := r.fileReader.Open("/proc/partitions")
	if err != nil {
		return []domain.Partition{}, err
	}
	defer file.Close()

	var partitions []domain.Partition
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 20 || line[:3] == "major" {
			continue // Skip header or short lines
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := fields[3]
		if isPartition(name) {
			partitions = append(partitions, domain.Partition{
				Name:       name,
				MountPoint: "",
				Filesystem: "",
				Total:      0,
				Used:       0,
				Free:       0,
			})
		}
	}

	return partitions, nil
}

func (r *StorageRepository) readMounts() (map[string]string, error) {
	file, err := r.fileReader.Open("/proc/mounts")
	if err != nil {
		return map[string]string{}, err
	}
	defer file.Close()

	mountPoints := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		device := strings.TrimPrefix(fields[0], "/dev/")

		if isPartition(device) {
			newMountPoint := fields[1]
			alreadyRegistered := false
			for registeredDevice, oldMountPoint := range mountPoints {
				if registeredDevice == device {
					alreadyRegistered = true
					if newPathIsCloserToRoot(oldMountPoint, newMountPoint) {
						mountPoints[device] = newMountPoint
						break
					}
				}
			}

			if !alreadyRegistered {
				mountPoints[device] = newMountPoint
			}

		}
	}

	return mountPoints, nil
}

func (r *StorageRepository) getFilesystemInfo() (map[string]domain.Partition, error) {

	if !r.toolChecker.isToolInstalled("df") {
		return map[string]domain.Partition{}, errors.New("df not installed")
	}

	var output []byte
	var err error

	cmd := r.cmdExec.Command("df", "-l", "--block-size=1") // Get sizes in bytes
	output, err = cmd.Output()
	if err != nil {
		// For busybox df command
		cmd := r.cmdExec.Command("df", "-B", "1")
		output, err = cmd.Output()
		if err != nil {
			log.Printf("Failed to run df command: %v", err)
			return map[string]domain.Partition{}, err
		}
	}

	return parseDfOutput(output), err

}
