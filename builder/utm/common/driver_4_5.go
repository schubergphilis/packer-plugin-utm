package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Utm45Driver is the base type for UTM drivers
type Utm45Driver struct {
	// This is the path to the utmctl binary
	UtmctlPath string
}

func (d *Utm45Driver) Delete(name string) error {
	_, err := d.Utmctl("delete", name)
	return err
}

// ExecuteOsaScript executes an AppleScript command with the given arguments.
func (d *Utm45Driver) ExecuteOsaScript(command ...string) (string, error) {
	if len(command) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	// log the command to be executed
	log.Printf("Executing OSA script command: %s", command)

	// Read the script content from the embedded files
	scriptPath := filepath.Join("scripts", command[0])
	scriptContent, err := osascripts.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read script %s: %v", scriptPath, err)
	}

	// Construct the command to execute
	cmd := exec.Command("osascript", "-")

	// Append additional arguments to the command
	if len(command) > 1 {
		cmd.Args = append(cmd.Args, command[1:]...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, string(scriptContent))
	}()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if stdoutString != "" {
		log.Printf("stdout: %s", stdoutString)
	}
	if stderrString != "" {
		log.Printf("stderr: %s", stderrString)
	}

	return stdoutString, err
}

// UTM 4.5 Doesn't support exporting VMs
func (d *Utm45Driver) Export(vmId string, path string) error {
	// just print a message to the user
	log.Printf("UTM API does not support exporting VMs yet.")
	log.Printf("Please manually export the VM using 'Share...' action in UTM VM menu.")
	log.Printf("Please make sure the VM is exported to the path %s ", path)
	log.Printf("The exported UTM file in the output directory will be passed as build Artifact.")

	// TODO: Pause and wait for the user to export the VM
	// // ask user to input the path of the exported file
	// confirmOption, err := ui.Ask(
	// 	fmt.Sprintf("Confirm you have exported the VM to path [%s] [Y/n]:", outputPath))

	// if err != nil {
	// 	err := fmt.Errorf("error during export step: %s", err)
	// 	state.Put("error", err)
	// 	ui.Error(err.Error())
	// 	return multistep.ActionHalt
	// }
	return nil
}

// UTM 4.5 : doesn't support adding support guest tools
func (d *Utm45Driver) GuestToolsIsoPath() (string, error) {
	return "", fmt.Errorf("UTM driver does not provide guest additions")
}

// UTM 4.5 : We just create a VM shortcut using UTM open command.
func (d *Utm45Driver) Import(path string) (string, error) {
	var stdout bytes.Buffer
	// TODO: While importing we should have ability to set the name of the VM
	// UTM does not support setting the name of the VM while importing
	// So we make sure VM name is same as the name in plist.config (previous name in UTM bundle)
	// This is a limitation of UTM
	cmd := exec.Command(
		"osascript", "-e",
		fmt.Sprintf(`tell application "UTM" to open POSIX file "%s"`, path),
	)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	// "missing value" in the output means AppleScript was successful
	// but not necessarily the VM was imported successfully
	// UTM does not provide a way to check if the VM was imported successfully
	// So we pray!
	// The error appears in UI, but not through script
	return "", nil
}

func (d *Utm45Driver) IsRunning(name string) (bool, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(d.UtmctlPath, "status", name)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return false, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "started" {
		return true, nil
	}

	// We consider "stopping" to still be running. We wait for it to
	// be completely stopped or some other state.
	if output == "stopping" {
		return true, nil
	}

	// We consider "paused" to still be running. We wait for it to
	// be completely stopped or some other state.
	if output == "paused" {
		return true, nil
	}

	// There might be other intermediate states that we consider
	// running, like "pausing", "resuming", "starting", etc.
	// but for now we just use these three.

	return false, nil
}

func (d *Utm45Driver) Stop(name string) error {
	if _, err := d.Utmctl("stop", name); err != nil {
		return err
	}
	return nil
}

func (d *Utm45Driver) Utmctl(args ...string) (string, error) {
	var stdout, stderr bytes.Buffer

	log.Printf("Executing utmctl: %#v", args)
	cmd := exec.Command(d.UtmctlPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("Utmctl error: %s", stderrString)
	}

	if stdoutString != "" {
		log.Printf("stdout: %s", stdoutString)
	}
	if stderrString != "" {
		log.Printf("stderr: %s", stderrString)
	}

	return stdoutString, err
}

func (d *Utm45Driver) Verify() error {
	return nil
}

// Version reads the version of UTM that is installed.
func (d *Utm45Driver) Version() (string, error) {
	var stdout bytes.Buffer

	cmd := exec.Command("osascript", "-e",
		`tell application "System Events" to return version of application "UTM"`)

	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	versionOutput := strings.TrimSpace(stdout.String())
	log.Printf("UTM version output: %s", versionOutput)

	// Check if the output contains the error message
	if strings.Contains(versionOutput, "get application") {
		return "", fmt.Errorf("UTM is not installed")
	}

	versionRe := regexp.MustCompile(`^(\d+\.\d+\.\d+)$`)
	matches := versionRe.FindStringSubmatch(versionOutput)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("no version found: %s", versionOutput)
	}

	log.Printf("UTM version: %s", matches[1])
	return matches[1], nil
}
