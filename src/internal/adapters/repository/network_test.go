package repository

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ******************************************** MOCKING ******************************************** */

type MockToolInstalled struct {
	Installed map[string]bool
}

func (m *MockToolInstalled) isToolInstalled(tool string) bool {
	if m.Installed != nil {
		if m.Installed[tool] {
			return m.Installed[tool]
		}
	}
	return false
}

type MockCmdExecutor struct {
	Cmd    *exec.Cmd
	output string
}

func (c *MockCmdExecutor) Command(name string, arg ...string) CmdExecutor {
	c.Cmd = exec.Command(name, arg...)
	return c
}

func (c *MockCmdExecutor) Output() ([]byte, error) {
	// fmt.Printf("Output: %s\n", c.output)
	return []byte(c.output), nil
}

/* ******************************************** NETWORK TEST ******************************************** */

func TestParseUint(t *testing.T) {
	testBattery := map[string]uint64{
		"100": 100,
		"0":   0,
		"-1":  0,
		"abc": 0,
		"":    0,
	}

	for input, expected := range testBattery {
		actual := parseUint(input)
		if actual != expected {
			t.Errorf("parseUint(%s) = %d; want %d", input, actual, expected)
		}
	}

}

func TestSpeedUnitMultiplier(t *testing.T) {
	testBattery := map[string]uint64{
		"b/s":            1,
		"Kb/s":           1000,
		"Mb/s":           1000000,
		"Gb/s":           1000000000,
		"incorrect unit": 1,
	}

	for input, expected := range testBattery {
		actual := speedUnitMultiplier(input)
		if actual != expected {
			t.Errorf("speedUnitMultiplier(%s) = %d; want %d", input, actual, expected)
		}
	}

}

func TestParseIwconfigOutput(t *testing.T) {

	testBattery := map[string]uint64{
		"Bit Rate=1 Gb/s":     1 * 1000 * 1000 * 1000,
		"Bit Rate=15 Mb/s":    15 * 1000 * 1000,
		"Bit Rate=100 Kb/s":   100 * 1000,
		"Bit Rate=100 b/s":    100,
		"incorrect file data": 0,
		"0":                   0,
		"-1":                  0,
		"abc":                 0,
		"":                    0,
	}

	for input, expected := range testBattery {
		actual, err := parseIwconfigOutput([]byte(input))
		if actual != expected {
			t.Errorf("parseIwconfigOutput(%s) = %d; want %d", input, actual, expected)
		}
		assert.NoError(t, err)
	}

}

func TestGetWiredSpeed(t *testing.T) {
	testBattery := map[string]uint64{
		"100":                100 * 1000 * 1000,
		"incorect file data": 0,
		"0":                  0,
		"-1":                 0,
		"abc":                0,
		"":                   0,
	}

	for input, expected := range testBattery {
		// Create a mock file reader with mock data
		mockFileReader := &MockFileReader{Data: input}

		// Other interfaces are not needed, so we go with the real ones
		execFinder := &RealToolInstalled{}
		cmd := &RealCmdExecutor{}

		// Instantiate the repository with the mock file reader
		repo := NewNetworkRepository(mockFileReader, execFinder, cmd)

		speed, err := repo.getWiredSpeed("eth0")

		assert.Equal(t, expected, speed)
		assert.NoError(t, err)
	}

}

func TestGetWirelessSpeed(t *testing.T) {

	var testBattery map[string]uint64
	iwconfigInstalled := isToolInstalled("iwconfig")
	var ti ToolInstalled
	var cmd CmdExecutor

	// Other interfaces are not needed, so we go with the real ones
	fileReader := &RealFileReader{}

	var repo *NetworkRepository

	testBattery = map[string]uint64{
		"Bit Rate=100 Mb/s":   100 * 1000 * 1000,
		"incorrect file data": 0,
		"0":                   0,
		"-1":                  0,
		"abc":                 0,
		"":                    0,
	}

	if !iwconfigInstalled {
		t.Log("iwconfig is not installed. Running tests mocking its behavior")

		ti = &MockToolInstalled{
			Installed: map[string]bool{
				"iwconfig": true,
			},
		}

	} else {
		t.Log("iwconfig is installed, running test")
		ti = &RealToolInstalled{}

	}

	for cmdOutput, expected := range testBattery {

		cmd = &MockCmdExecutor{output: cmdOutput}
		repo = NewNetworkRepository(fileReader, ti, cmd)

		speed, err := repo.getWirelessSpeed("wlan0")

		if expected != speed {
			t.Errorf("getWirelessSpeed(%s) = %d; want %d", cmdOutput, speed, expected)
		}
		assert.NoError(t, err)

	}

}

func TestGetNetworkInterfaces(t *testing.T) {

	var repo *NetworkRepository

	ti := &MockToolInstalled{
		Installed: map[string]bool{
			"iwconfig": true,
		},
	}

	// We don't test with wired device because current mock structure doesn't support mocking
	// file reading one by one, and wired iface link speed is retrieved from a file
	testBattery := map[string]map[string]any{
		"Case 1 - No devices": {
			"filereader":  &MockFileReader{Data: ""},
			"cmdExecutor": &MockCmdExecutor{},
			"expected": map[string]any{
				"ifaces": []string{},
			},
		},
		"Case 2 - One device wireless": {
			"filereader": &MockFileReader{Data: `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
 wlan0: 203957342  161674    0    0    0     0          0      3658 32173524   51534    0    0    0     0       0          0
`},
			"cmdExecutor": &MockCmdExecutor{output: "Bit Rate=100 Mb/s"},
			"expected": map[string]any{
				"ifaces": []string{"wlan0"},
				"wlan0": map[string]uint64{
					"rx_bytes":   203957342,
					"tx_bytes":   32173524,
					"rx_packets": 161674,
					"tx_packets": 51534,
					"rx_errors":  0,
					"tx_errors":  0,
					"rx_dropped": 0,
					"tx_dropped": 0,
				},
			},
		},
	}

	for caseName, caseData := range testBattery {
		t.Log(caseName)
		repo = NewNetworkRepository(caseData["filereader"].(FileReader), ti, caseData["cmdExecutor"].(CmdExecutor))

		resultIfaces, err := repo.GetNetworkInterfaces()

		assert.NoError(t, err)
		caseExpected := caseData["expected"].(map[string]any)

		assert.Len(t, resultIfaces, len(caseExpected["ifaces"].([]string)))

		for _, expectedIfaceName := range caseExpected["ifaces"].([]string) {

			found := false

			for _, resultIface := range resultIfaces {

				if resultIface.InterfaceName == expectedIfaceName {
					expectedIface := caseExpected[expectedIfaceName].(map[string]uint64)
					assert.Equal(t, expectedIface["rx_bytes"], resultIface.Rx.Bytes)
					t.Log(resultIface.Tx.Bytes)
					t.Log(expectedIface["tx_bytes"])
					assert.Equal(t, expectedIface["tx_bytes"], resultIface.Tx.Bytes)
					assert.Equal(t, expectedIface["rx_packets"], resultIface.Rx.Packets)
					assert.Equal(t, expectedIface["tx_packets"], resultIface.Tx.Packets)
					assert.Equal(t, expectedIface["rx_errors"], resultIface.Rx.Errors)
					assert.Equal(t, expectedIface["tx_errors"], resultIface.Tx.Errors)
					assert.Equal(t, expectedIface["rx_dropped"], resultIface.Rx.Drops)
					assert.Equal(t, expectedIface["tx_dropped"], resultIface.Tx.Drops)
					found = true
					break
				}

			}

			assert.True(t, found, "Interface %s not found in the list", expectedIfaceName)

		}

	}

}
