package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

// Utm46Driver are inherited from Utm45Driver.
type Utm46Driver struct {
	Utm45Driver
}

// UTM 4.6 : We import a VM by utm file using UTM import command.
func (d *Utm46Driver) Import(name string, path string) (string, error) {
	var stdout bytes.Buffer

	// Import VM
	cmd := exec.Command(
		"osascript", "-e",
		fmt.Sprintf(`tell application "UTM" to import new virtual machine from POSIX file "%s"`, path),
	)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Get the output of the command
	output := stdout.String()

	// Regular expression to capture the VM ID
	re := regexp.MustCompile(`virtual machine id ([A-F0-9-]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		vmId := matches[1] // Capture the VM ID
		return vmId, nil
	}

	return "", fmt.Errorf("failed to import VM: %s", output)
}
