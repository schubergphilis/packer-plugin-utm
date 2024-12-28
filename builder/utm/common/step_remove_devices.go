// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step removes any devices (floppy disks, ISOs, etc.) from the
// machine that we may have added.
//
// Uses:
//
//	driver Driver
//	ui packersdk.Ui
//	vmName string
//
// Produces:
type StepRemoveDevices struct {
	Bundling UtmBundleConfig
}

func (s *StepRemoveDevices) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)

	// TODO: Remove the attached floppy disk, if it exists

	var isoUnmountCommands map[string][]string
	isoUnmountCommandsRaw, ok := state.GetOk("disk_unmount_commands")
	if !ok {
		// No disks to unmount
		return multistep.ActionContinue
	} else {
		isoUnmountCommands = isoUnmountCommandsRaw.(map[string][]string)
	}

	for diskCategory, unmountCommand := range isoUnmountCommands {
		if diskCategory == "boot_iso" && s.Bundling.BundleISO {
			// skip the unmount if user wants to bundle the iso
			continue
		}

		if _, err := driver.ExecuteOsaScript(unmountCommand...); err != nil {
			err := fmt.Errorf("Error detaching ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	// log that we removed the isos, so we don't waste time trying to do it
	// in the step_attach_isos cleanup.
	state.Put("detached_isos", true)

	return multistep.ActionContinue
}

func (s *StepRemoveDevices) Cleanup(state multistep.StateBag) {
}
