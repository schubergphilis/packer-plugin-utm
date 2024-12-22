package cloud

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
type stepCreateCloudVM struct {
	vmId string
}

func (s *stepCreateCloudVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(utmcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)

	vmName := config.VMName
	isoPath := state.Get("iso_path").(string)

	// Create VM command
	createCommand := []string{
		"create_vm_from_source.applescript", "--name", vmName,
		"--backend", config.VMBackend,
		"--arch", config.VMArch,
		"--source", isoPath,
	}

	ui.Say("Creating virtual machine...")
	output, err := driver.ExecuteOsaScript(createCommand...)
	if err != nil {
		err := fmt.Errorf("error creating VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Regular expression to capture the VM ID
	re := regexp.MustCompile(`virtual machine id ([A-F0-9-]+)`)
	matches := re.FindStringSubmatch(output)
	var vmId string
	if len(matches) > 1 {
		vmId = matches[1] // Capture the VM ID
		s.vmId = vmId
		state.Put("vmId", s.vmId)
		// save the vm name, used in export step
		state.Put("vmName", vmName)
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

func (s *stepCreateCloudVM) Cleanup(state multistep.StateBag) {
	if s.vmId == "" {
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
	if err := driver.Delete(s.vmId); err != nil {
		ui.Error(fmt.Sprintf("Error deleting VM: %s", err))
	}
}
