package cloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	utmcommon "github.com/naveenrajm7/packer-plugin-utm/builder/utm/common"
)

// This step configures the VM to send cloud init seed file.
//
// Uses:
//
//	config *config
//	ui     packersdk.Ui
type stepConfigureCloudSeed struct{}

func (s *stepConfigureCloudSeed) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(utmcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)
	vmId := state.Get("vmId").(string)

	// Get the host IP and HTTP port
	// hostIP := state.Get("http_ip").(string)
	// TODO: Get Host IP automatically from Mac's Shared_Network_Address
	hostIP := "192.168.75.1"
	httpPort := state.Get("http_port").(int)

	// Add Qemu args to send cloud init seed file
	ui.Say("Configuring VM to send cloud init seed file...")
	cloudQemuArg := fmt.Sprintf("-smbios type=1,serial=ds=nocloud-net;seedfrom=http://%s:%d/", hostIP, httpPort)
	addQemuArgsCommand := []string{
		"add_qemu_additional_args.applescript", vmId,
		"--args", cloudQemuArg,
	}

	ui.Say("Adding QEMU additional arguments...")
	_, err := driver.ExecuteOsaScript(addQemuArgsCommand...)
	if err != nil {
		err := fmt.Errorf("error adding QEMU additional arguments: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	// Save the Cloud QEMU argument for later cleanup
	state.Put("qemuAdditionalArg", cloudQemuArg)

	return multistep.ActionContinue
}

func (s *stepConfigureCloudSeed) Cleanup(multistep.StateBag) {}
