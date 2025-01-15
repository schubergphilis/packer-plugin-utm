package common

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step creates the actual virtual machine.
//
// Produces:
//
//	vmId string - The UUID of the VM
type StepCreateVM struct {
	// takes
	VMName         string
	VMBackend      string
	VMArch         string
	HWConfig       HWConfig
	UEFIBoot       bool
	Hypervisor     bool
	KeepRegistered bool
	// produces
	vmId string
}

func (s *StepCreateVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)

	// Create VM command
	createCommand := []string{
		"create_vm.applescript", "--name", s.VMName,
		"--backend", s.VMBackend,
		"--arch", s.VMArch,
	}

	ui.Say("Creating virtual machine...")
	output, err := driver.ExecuteOsaScript(createCommand...)
	if err != nil {
		err := fmt.Errorf("error creating VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Regular expression to capture the VM UUID
	re := regexp.MustCompile(`[0-9a-fA-F-]{36}`)
	matches := re.FindStringSubmatch(output)
	var vmId string
	if len(matches) > 0 {
		vmId = matches[0] // Capture the VM UUID
		s.vmId = vmId
		state.Put("vmName", s.VMName)
		state.Put("vmId", s.vmId)
	} else {
		err := fmt.Errorf("error extracting VM ID from output: %s", output)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("VM Id: %s", vmId)

	// Customize VM command
	customizeCommand := []string{
		"customize_vm.applescript", vmId,
		"--cpus", strconv.Itoa(s.HWConfig.CpuCount),
		"--memory", strconv.Itoa(s.HWConfig.MemorySize),
		"--name", s.VMName,
		"--uefi-boot", strconv.FormatBool(s.UEFIBoot),
		"--use-hypervisor", strconv.FormatBool(s.Hypervisor),
	}

	ui.Say("Customizing virtual machine...")
	_, err = driver.ExecuteOsaScript(customizeCommand...)
	if err != nil {
		err := fmt.Errorf("error customizing VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCreateVM) Cleanup(state multistep.StateBag) {
	if s.vmId == "" {
		return
	}

	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if (s.KeepRegistered) && (!cancelled && !halted) {
		ui.Say("Keeping virtual machine registered with UTM host (keep_registered = true)")
		return
	}

	ui.Say("Deregistering and deleting VM...")
	if err := driver.Delete(s.vmId); err != nil {
		ui.Error(fmt.Sprintf("Error deleting VM: %s", err))
	}
}
