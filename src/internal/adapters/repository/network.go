package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"pi-monitor-api/internal/core/domain"
	"strconv"
	"strings"
)

/* ******************************************** AUX ******************************************** */

func parseUint(s string) uint64 {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func isToolInstalled(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

func speedUnitMultiplier(unit string) uint64 {
	switch unit {
	case "Kb/s":
		return 1000
	case "Mb/s":
		return 1000000
	case "Gb/s":
		return 1000000000
	default:
		return 1
	}
}

func parseIwconfigOutput(output []byte) (uint64, error) {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Bit Rate=") {
			parts := strings.Split(line, "Bit Rate=")
			if len(parts) < 2 {
				continue
			}
			speedPart := strings.Fields(parts[1])[0]
			speedUnit := strings.Fields(parts[1])[1]
			speed := parseUint(speedPart)
			return speed * speedUnitMultiplier(speedUnit), nil
		}
	}
	return 0, nil
}

/* ************************************* MOCKING SCAFFOLDING ************************************* */

type ToolInstalled interface {
	isToolInstalled(tool string) bool
}

type RealToolInstalled struct{}

func (r *RealToolInstalled) isToolInstalled(tool string) bool {
	return isToolInstalled(tool)
}

type CmdExecutor interface {
	Command(name string, arg ...string) CmdExecutor
	Output() ([]byte, error)
}

type RealCmdExecutor struct {
	Cmd *exec.Cmd
}

func (c *RealCmdExecutor) Command(name string, arg ...string) CmdExecutor {
	c.Cmd = exec.Command(name, arg...)
	return c
}

func (c *RealCmdExecutor) Output() ([]byte, error) {
	return c.Cmd.Output()
}

/* ******************************************** NETWORK ******************************************** */

type NetworkRepository struct {
	fileReader  FileReader
	toolChecker ToolInstalled
	cmdExec     CmdExecutor
}

func NewNetworkRepository(fr FileReader, ti ToolInstalled, cmd CmdExecutor) *NetworkRepository {
	return &NetworkRepository{
		fileReader:  fr,
		toolChecker: ti,
		cmdExec:     cmd,
	}
}

func (r *NetworkRepository) GetNetworkInterfaces() ([]domain.NetworkInterface, error) {
	file, err := r.fileReader.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var interfaces []domain.NetworkInterface
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, ":") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}

		interfaceName := strings.TrimSuffix(fields[0], ":")
		rxStats := domain.NetworkStats{
			Bytes:   parseUint(fields[1]),
			Packets: parseUint(fields[2]),
			Errors:  parseUint(fields[3]),
			Drops:   parseUint(fields[4]),
		}
		txStats := domain.NetworkStats{
			Bytes:   parseUint(fields[9]),
			Packets: parseUint(fields[10]),
			Errors:  parseUint(fields[11]),
			Drops:   parseUint(fields[12]),
		}

		bitRate, _ := r.getLinkSpeed(interfaceName)

		iface := domain.NetworkInterface{
			InterfaceName: interfaceName,
			BitRate:       bitRate,
			Rx:            rxStats,
			Tx:            txStats,
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

func (r *NetworkRepository) getLinkSpeed(interfaceName string) (uint64, error) {
	if speed, err := r.getWirelessSpeed(interfaceName); err == nil {
		return speed, nil
	}

	return r.getWiredSpeed(interfaceName)
}

func (r *NetworkRepository) getWirelessSpeed(interfaceName string) (uint64, error) {

	if !r.toolChecker.isToolInstalled("iwconfig") {
		return 0, errors.New("iwconfig not installed")
	}

	cmd := r.cmdExec.Command("iwconfig", interfaceName)
	output, err := cmd.Output() // Run the command and capture the output
	if err != nil {
		return 0, err
	}

	// Parse the output to find the speed
	return parseIwconfigOutput(output)
}

func (r *NetworkRepository) getWiredSpeed(interfaceName string) (uint64, error) {
	filePath := "/sys/class/net/" + interfaceName + "/speed"
	file, err := r.fileReader.Open(filePath)
	if err != nil {
		fmt.Println("Error opening wired iface speed file: ", err)
		return 0, errors.New("error opening iface speed file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		speedStr := scanner.Text()
		speed := parseUint(speedStr)
		return speed * 1000000, nil // Convert to bits per second
	}

	return 0, nil
}
