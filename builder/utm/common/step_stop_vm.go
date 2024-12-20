// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step stops the machine.
//
// Uses:
//
//	driver Driver
//	ui     packersdk.Ui
//	vmId string
//
// Produces:
//
//	<nothing>
type StepStopVm struct{}

func (s *StepStopVm) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)
	vmId := state.Get("vmId").(string)

	ui.Say("Stopping virtual machine...")
	if err := driver.Stop(vmId); err != nil {
		err := fmt.Errorf("error stopping VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	} else {
		ui.Say("Automatic stop failed. Please stop the machine.")
	}

	log.Println("VM stopped.")
	return multistep.ActionContinue
}

func (s *StepStopVm) Cleanup(state multistep.StateBag) {}
