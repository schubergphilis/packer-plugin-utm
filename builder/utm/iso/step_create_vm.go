package iso

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	utmcommon "github.com/naveenrajm7/packer-plugin-utm/builder/utm/common"
)

// This step creates the actual virtual machine.
//
// Produces:
//
//	vmId string - The UUID of the VM
type stepCreateVM struct {
	vmName string
	vmId   string
}

func (s *stepCreateVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(utmcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)

	vmName := config.VMName

	// Create VM command
	createCommand := []string{
		"create_vm.applescript", "--name", vmName,
		"--backend", config.VMBackend,
		"--arch", config.VMArch,
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
		s.vmName = vmName
		s.vmId = vmId
		state.Put("vmName", s.vmName)
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
		"--cpus", strconv.Itoa(config.HWConfig.CpuCount),
		"--memory", strconv.Itoa(config.HWConfig.MemorySize),
		"--name", vmName,
		"--uefi-boot", strconv.FormatBool(config.UEFIBoot),
		"--hypervisor", strconv.FormatBool(config.Hypervisor),
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

func (s *stepCreateVM) Cleanup(state multistep.StateBag) {
	if s.vmName == "" {
		return
	}

	driver := state.Get("driver").(utmcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if (config.KeepRegistered) && (!cancelled && !halted) {
		ui.Say("Keeping virtual machine registered with UTM host (keep_registered = true)")
		return
	}

	ui.Say("Deregistering and deleting VM...")
	if err := driver.Delete(s.vmName); err != nil {
		ui.Error(fmt.Sprintf("Error deleting VM: %s", err))
	}
}
