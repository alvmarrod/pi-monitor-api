package repository

import (
	"testing"

	"github.com/alvmarrod/pi-monitor-api/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

/* ******************************************** MOCKING ******************************************** */

// All the mock structures and functions are defined already in other repositories files

/* ******************************************** AUX TEST ******************************************** */

func TestIsPartition(t *testing.T) {
	assert.True(t, isPartition("sda"))
	assert.True(t, isPartition("sda1"))
	assert.True(t, isPartition("sdb1"))
	assert.True(t, isPartition("sdc1"))
	assert.True(t, isPartition("nvme0n1"))
	assert.True(t, isPartition("nvme0n1p1"))
	assert.False(t, isPartition("some random string"))
	assert.False(t, isPartition("lvm1"))
}

func TestGetDeviceName(t *testing.T) {
	assert.Equal(t, "sda", getDeviceName("sda1"))
	assert.Equal(t, "sdb", getDeviceName("sdb2"))
	assert.Equal(t, "sdc", getDeviceName("sdc3"))
	assert.Equal(t, "nvme0", getDeviceName("nvme0n1p1"))
	assert.Equal(t, "", getDeviceName("mmc"))
}

func TestGroupDevices(t *testing.T) {
	partitions := []domain.Partition{
		{Name: "sda1", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
		{Name: "sda2", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
		{Name: "sdb1", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
		{Name: "sdc1", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
		{Name: "nvme0n1p1", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
		{Name: "nvme0n1p2", Filesystem: "ext4", Total: 100, Used: 50, Free: 50},
	}

	expectedDevices := []domain.Device{
		{
			Name: "sda",
			Partitions: map[string]domain.Partition{
				"sda1": {},
				"sda2": domain.Partition{},
			},
		},
		{
			Name: "sdb",
			Partitions: map[string]domain.Partition{
				"sdb1": domain.Partition{},
			},
		},
		{
			Name: "sdc",
			Partitions: map[string]domain.Partition{
				"sdc1": domain.Partition{},
			},
		},
		{
			Name: "nvme0",
			Partitions: map[string]domain.Partition{
				"nvme0n1p1": domain.Partition{},
				"nvme0n1p2": domain.Partition{},
			},
		},
	}

	foundDevices := groupDevices(partitions)
	t.Log(foundDevices)

	assert.Len(t, foundDevices, 4)
	for _, exp_device := range expectedDevices {
		found := false
		for _, found_device := range foundDevices {
			if found_device.Name == exp_device.Name {
				found = true
				assert.Len(t, found_device.Partitions, len(exp_device.Partitions))
				break
			}
		}
		assert.True(t, found, "Device %s not found in the list", exp_device.Name)
	}

}

func TestParseDfOutput(t *testing.T) {

	testBattery := map[string]map[string]any{
		"Case 1 - Empty file": {
			"input":    []byte(""),
			"expected": map[string]domain.Partition{},
		},
		"Case 2 - Incorrect file": {
			"input":    []byte("some incorrect file data"),
			"expected": map[string]domain.Partition{},
		},
		"Case 3 - Correct file": {
			"input": []byte(`Filesystem         1B-blocks         Used    Available Use% Mounted on
			none              3991056384         4096   3991052288   1% /mnt/wsl
			drivers        1023117619200 982047334400  41070284800  96% /usr/lib/wsl/drivers
			none              3991056384            0   3991056384   0% /usr/lib/modules
			none              3991056384            0   3991056384   0% /usr/lib/modules/5.15.153.1-microsoft-standard-WSL2
			/dev/sdc        269427478528  37333770240 218333036544  15% /
			none              3991056384       131072   3990925312   1% /mnt/wslg
			none              3991056384            0   3991056384   0% /usr/lib/wsl/lib
			rootfs            3987623936      2129920   3985494016   1% /init
			none              3987623936            0   3987623936   0% /dev
			none              3991056384        20480   3991035904   1% /run
			none              3991056384            0   3991056384   0% /run/lock
			none              3991056384            0   3991056384   0% /run/shm
			none              3991056384            0   3991056384   0% /run/user
			tmpfs             3991056384            0   3991056384   0% /sys/fs/cgroup
			none              3991056384       577536   3990478848   1% /mnt/wslg/versions.txt
			none              3991056384       577536   3990478848   1% /mnt/wslg/doc`),
			"expected": map[string]domain.Partition{
				"sdc": {
					Name:       "sdc",
					Filesystem: "/dev/sdc",
					MountPoint: "/",
					Total:      269427478528,
					Used:       37333770240,
					Free:       218333036544,
				},
			},
		},
	}

	for caseName, caseData := range testBattery {
		t.Log(caseName)
		resultPartitions := parseDfOutput(caseData["input"].([]byte))

		for _, expectedPartition := range caseData["expected"].(map[string]domain.Partition) {

			found := false
			for _, resultPartition := range resultPartitions {
				if resultPartition.Name == expectedPartition.Name {
					found = true
					assert.Equal(t, expectedPartition.Filesystem, resultPartition.Filesystem)
					assert.Equal(t, expectedPartition.Total, resultPartition.Total)
					assert.Equal(t, expectedPartition.Used, resultPartition.Used)
					assert.Equal(t, expectedPartition.Free, resultPartition.Free)
					break
				}
			}
			assert.True(t, found, "Partition %s not found in the list", expectedPartition.Name)

		}

	}

}

/* ***************************************** STORAGE TESTS ***************************************** */

func TestGetFilesystemInfo(t *testing.T) {

	testBattery := map[string]map[string]any{
		"Case 1 - Empty file": {
			"input":    []byte(""),
			"expected": []domain.Partition{},
		},
		"Case 2 - Incorrect file": {
			"input":    []byte("some incorrect file data"),
			"expected": []domain.Partition{},
		},
		"Case 3 - Correct file": {
			"input": []byte(`Filesystem         1B-blocks         Used    Available Use% Mounted on
none              3991056384         4096   3991052288   1% /mnt/wsl
drivers        1023117619200 982592319488  40525299712  97% /usr/lib/wsl/drivers
none              3991056384            0   3991056384   0% /usr/lib/modules
none              3991056384            0   3991056384   0% /usr/lib/modules/5.15.153.1-microsoft-standard-WSL2
/dev/sdc        269427478528  37230718976 218436087808  15% /
none              3991056384       151552   3990904832   1% /mnt/wslg
none              3991056384            0   3991056384   0% /usr/lib/wsl/lib
rootfs            3987623936      2129920   3985494016   1% /init
none              3987623936            0   3987623936   0% /dev
none              3991056384        20480   3991035904   1% /run
none              3991056384            0   3991056384   0% /run/lock
none              3991056384            0   3991056384   0% /run/shm
none              3991056384            0   3991056384   0% /run/user
tmpfs             3991056384            0   3991056384   0% /sys/fs/cgroup
none              3991056384       577536   3990478848   1% /mnt/wslg/versions.txt
none              3991056384       577536   3990478848   1% /mnt/wslg/doc`),
			"expected": []domain.Partition{
				{
					Name:       "sdc",
					Filesystem: "/dev/sdc",
					MountPoint: "/",
					Total:      269427478528,
					Used:       37230718976,
					Free:       218436087808,
				},
			},
		},
	}

	for caseName, caseData := range testBattery {

		t.Log(caseName)

		ti := &MockToolInstalled{
			Installed: map[string]bool{
				"df": true,
			},
		}
		cmd := &MockCmdExecutor{output: string(caseData["input"].([]byte))}

		// Other interfaces are not needed, so we go with the real ones
		fr := &RealFileReader{}

		repo := NewStorageRepository(fr, ti, cmd)

		// Execution
		resultPartitions, err := repo.getFilesystemInfo()

		// Validation
		assert.NoError(t, err)
		assert.Len(t, resultPartitions, len(caseData["expected"].([]domain.Partition)))

		for _, expectedPartition := range caseData["expected"].([]domain.Partition) {

			found := false
			for _, resultPartition := range resultPartitions {
				if resultPartition.Name == expectedPartition.Name {
					found = true
					assert.Equal(t, expectedPartition.Filesystem, resultPartition.Filesystem)
					assert.Equal(t, expectedPartition.Total, resultPartition.Total)
					assert.Equal(t, expectedPartition.Used, resultPartition.Used)
					assert.Equal(t, expectedPartition.Free, resultPartition.Free)
					break
				}
			}
			assert.True(t, found, "Partition %s not found in the list", expectedPartition.Name)

		}

	}

}

func TestReadMounts(t *testing.T) {

	testBattery := map[string]map[string]any{
		"Case 1 - Empty file": {
			"input":    []byte(""),
			"expected": map[string]string{},
		},
		"Case 2 - Incorrect file": {
			"input":    []byte("some incorrect data"),
			"expected": map[string]string{},
		},
		"Case 3 - Correct file": {
			"input": []byte(`none /mnt/wsl tmpfs rw,relatime 0 0
drivers /usr/lib/wsl/drivers 9p ro,dirsync,nosuid,nodev,noatime,aname=drivers;fmask=222;dmask=222,mmap,access=client,msize=65536,trans=fd,rfd=7,wfd=7 0 0
none /usr/lib/modules tmpfs rw,relatime 0 0
none /usr/lib/modules/5.15.153.1-microsoft-standard-WSL2 overlay rw,nosuid,nodev,noatime,lowerdir=/modules,upperdir=/modules_overlay/rw/upper,workdir=/modules_overlay/rw/work 0 0
/dev/sdc / ext4 rw,relatime,discard,errors=remount-ro,data=ordered 0 0
none /mnt/wslg tmpfs rw,relatime 0 0
/dev/sdc /mnt/wslg/distro ext4 ro,relatime,discard,errors=remount-ro,data=ordered 0 0
none /usr/lib/wsl/lib overlay rw,nosuid,nodev,noatime,lowerdir=/gpu_lib_packaged:/gpu_lib_inbox,upperdir=/gpu_lib/rw/upper,workdir=/gpu_lib/rw/work 0 0
rootfs /init rootfs ro,size=3894164k,nr_inodes=973541 0 0
none /dev devtmpfs rw,nosuid,relatime,size=3894164k,nr_inodes=973541,mode=755 0 0
sysfs /sys sysfs rw,nosuid,nodev,noexec,noatime 0 0
proc /proc proc rw,nosuid,nodev,noexec,noatime 0 0
devpts /dev/pts devpts rw,nosuid,noexec,noatime,gid=5,mode=620,ptmxmode=000 0 0
none /run tmpfs rw,nosuid,nodev,mode=755 0 0
none /run/lock tmpfs rw,nosuid,nodev,noexec,noatime 0 0
none /run/shm tmpfs rw,nosuid,nodev,noatime 0 0
none /dev/shm tmpfs rw,nosuid,nodev,noatime 0 0
none /run/user tmpfs rw,nosuid,nodev,noexec,noatime,mode=755 0 0
binfmt_misc /proc/sys/fs/binfmt_misc binfmt_misc rw,relatime 0 0
tmpfs /sys/fs/cgroup tmpfs rw,nosuid,nodev,noexec,relatime,mode=755 0 0
cgroup2 /sys/fs/cgroup/unified cgroup2 rw,nosuid,nodev,noexec,relatime,nsdelegate 0 0
cgroup /sys/fs/cgroup/cpuset cgroup rw,nosuid,nodev,noexec,relatime,cpuset 0 0
cgroup /sys/fs/cgroup/cpu cgroup rw,nosuid,nodev,noexec,relatime,cpu 0 0
cgroup /sys/fs/cgroup/cpuacct cgroup rw,nosuid,nodev,noexec,relatime,cpuacct 0 0
cgroup /sys/fs/cgroup/blkio cgroup rw,nosuid,nodev,noexec,relatime,blkio 0 0
cgroup /sys/fs/cgroup/memory cgroup rw,nosuid,nodev,noexec,relatime,memory 0 0
cgroup /sys/fs/cgroup/devices cgroup rw,nosuid,nodev,noexec,relatime,devices 0 0
cgroup /sys/fs/cgroup/freezer cgroup rw,nosuid,nodev,noexec,relatime,freezer 0 0
cgroup /sys/fs/cgroup/net_cls cgroup rw,nosuid,nodev,noexec,relatime,net_cls 0 0
cgroup /sys/fs/cgroup/perf_event cgroup rw,nosuid,nodev,noexec,relatime,perf_event 0 0
cgroup /sys/fs/cgroup/net_prio cgroup rw,nosuid,nodev,noexec,relatime,net_prio 0 0
cgroup /sys/fs/cgroup/hugetlb cgroup rw,nosuid,nodev,noexec,relatime,hugetlb 0 0
cgroup /sys/fs/cgroup/pids cgroup rw,nosuid,nodev,noexec,relatime,pids 0 0
cgroup /sys/fs/cgroup/rdma cgroup rw,nosuid,nodev,noexec,relatime,rdma 0 0
cgroup /sys/fs/cgroup/misc cgroup rw,nosuid,nodev,noexec,relatime,misc 0 0
none /mnt/wslg/versions.txt overlay rw,relatime,lowerdir=/systemvhd,upperdir=/system/rw/upper,workdir=/system/rw/work 0 0
none /mnt/wslg/doc overlay rw,relatime,lowerdir=/systemvhd,upperdir=/system/rw/upper,workdir=/system/rw/work 0 0
none /mnt/wslg/.X11-unix tmpfs ro,relatime 0 0
C:\134 /mnt/c 9p rw,dirsync,noatime,aname=drvfs;path=C:\;uid=1000;gid=1000;symlinkroot=/mnt/,mmap,access=client,msize=65536,trans=fd,rfd=4,wfd=4 0 0
/dev/sdc /var/lib/docker ext4 rw,relatime,discard,errors=remount-ro,data=ordered 0 0`),
			"expected": map[string]string{
				"sdc": "/",
			},
		},
	}

	for caseName, caseData := range testBattery {

		t.Log(caseName)

		fr := &MockFileReader{Data: string(caseData["input"].([]byte))}

		// Other interfaces are not needed, so we go with the real ones
		cmd := &RealCmdExecutor{}
		ti := &RealToolInstalled{}

		repo := NewStorageRepository(fr, ti, cmd)

		// Execution
		resultMounts, err := repo.readMounts()

		// Validation
		assert.NoError(t, err)
		assert.Len(t, resultMounts, len(caseData["expected"].(map[string]string)))

		for expectedDevice, expectedMount := range caseData["expected"].(map[string]string) {

			found := false
			for resultDevice, resultMount := range resultMounts {
				if resultDevice == expectedDevice {
					found = true
					assert.Equal(t, expectedMount, resultMount)
					break
				}
			}

			assert.True(t, found, "Device %s not found in the list", expectedDevice)
		}

	}

}

func TestReadBasicPartitions(t *testing.T) {

	testBattery := map[string]map[string]any{
		"Case 1 - Empty file": {
			"input":    []byte(""),
			"expected": []domain.Partition{},
		},
		"Case 2 - Incorrect file": {
			"input":    []byte("some incorrect data"),
			"expected": []domain.Partition{},
		},
		"Case 3 - Correct file": {
			"input": []byte(`major minor  #blocks  name

   1        0      65536 ram0
   1        1      65536 ram1
   1        2      65536 ram2
   1        3      65536 ram3
   1        4      65536 ram4
   8        0     397928 sda
   8       16    2097156 sdb
   8       32  268435456 sdc`),
			"expected": []domain.Partition{
				{Name: "sda", Filesystem: "", Total: 0, Used: 0, Free: 0},
				{Name: "sdb", Filesystem: "", Total: 0, Used: 0, Free: 0},
				{Name: "sdc", Filesystem: "", Total: 0, Used: 0, Free: 0},
			},
		},
	}

	for caseName, caseData := range testBattery {

		t.Log(caseName)

		fr := &MockFileReader{Data: string(caseData["input"].([]byte))}

		// Other interfaces are not needed, so we go with the real ones
		cmd := &RealCmdExecutor{}
		ti := &RealToolInstalled{}

		repo := NewStorageRepository(fr, ti, cmd)

		// Execution
		resultPartitions, err := repo.readBasicPartitions()

		// Validation
		assert.NoError(t, err)
		assert.Len(t, resultPartitions, len(caseData["expected"].([]domain.Partition)))

		for _, expectedPartitions := range caseData["expected"].([]domain.Partition) {
			assert.Contains(t, resultPartitions, expectedPartitions)
		}

	}

}

func TestReadPartitions(t *testing.T) {

}
