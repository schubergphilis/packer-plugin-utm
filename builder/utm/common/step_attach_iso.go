package common

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step attaches the boot ISO, cd_files iso, and guest additions to the
// virtual machine, if present.
type StepAttachISOs struct {
	AttachBootISO           bool
	ISOInterface            string
	GuestAdditionsMode      string
	GuestAdditionsInterface string
	diskUnmountCommands     map[string][]string
}

func (s *StepAttachISOs) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// Check whether there is anything to attach
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Mounting ISOs...")
	diskMountMap := map[string]string{}
	// Track the bootable iso (only used in utm-iso builder. )
	if s.AttachBootISO {
		isoPath := state.Get("iso_path").(string)
		diskMountMap["boot_iso"] = isoPath
	}

	// Determine if we even have a cd_files disk to attach
	if cdPathRaw, ok := state.GetOk("cd_path"); ok {
		cdFilesPath := cdPathRaw.(string)
		diskMountMap["cd_files"] = cdFilesPath
	}

	// Determine if we have guest additions to attach
	if s.GuestAdditionsMode != GuestAdditionsModeAttach {
		log.Println("Not attaching guest additions since we're uploading.")
	} else {
		// Get the guest additions path since we're doing it
		guestAdditionsPath := state.Get("guest_additions_path").(string)
		diskMountMap["guest_additions"] = guestAdditionsPath
	}

	if len(diskMountMap) == 0 {
		ui.Message("No ISOs to mount; continuing...")
		return multistep.ActionContinue
	}

	driver := state.Get("driver").(Driver)
	vmId := state.Get("vmId").(string)

	// Iterate over the ISOs to attach
	// Attach one after the other
	// if you need to order the ISOs, this will not work. (Use a slice of structs instead)
	for diskCategory, isoPath := range diskMountMap {
		// If it's a symlink, resolve it to its target.
		resolvedIsoPath, err := filepath.EvalSymlinks(isoPath)
		if err != nil {
			err := fmt.Errorf("error resolving symlink for ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		isoPath = resolvedIsoPath

		// We may have different potential iso we can attach.
		var controllerName string
		switch diskCategory {
		case "boot_iso":
			controllerName = s.ISOInterface
			ui.Message("Mounting boot ISO...")
		case "guest_additions":
			controllerName = s.GuestAdditionsInterface
			ui.Message("Mounting guest additions ISO...")
		case "cd_files":
			controllerName = "usb"
			ui.Message("Mounting cd_files ISO...")
		}

		// Convert controllerName to the corresponding enum code
		controllerEnumCode, err := GetControllerEnumCode(controllerName)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		// Attach the ISO
		command := []string{
			"attach_iso.applescript", vmId,
			"--interface", controllerEnumCode,
			"--source", isoPath,
		}

		output, err := driver.ExecuteOsaScript(command...)
		if err != nil {
			err := fmt.Errorf("error attaching ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Track the disks we've mounted so we can remove them without having
		// to re-derive what was mounted where

		// Regular expression to capture the UUID from the output
		re := regexp.MustCompile(`[0-9a-fA-F-]{36}`)
		matches := re.FindStringSubmatch(output)
		if len(matches) > 0 {
			uuid := matches[0] // Capture the UUID
			unmountCommand := []string{
				"remove_drive.applescript", vmId, uuid,
			}
			s.diskUnmountCommands[diskCategory] = unmountCommand
		} else {
			err := fmt.Errorf("error extracting UUID from output: %s", output)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	state.Put("disk_unmount_commands", s.diskUnmountCommands)
	return multistep.ActionContinue
}

func (s *StepAttachISOs) Cleanup(state multistep.StateBag) {
	if len(s.diskUnmountCommands) == 0 {
		return
	}

	driver := state.Get("driver").(Driver)
	_, ok := state.GetOk("detached_isos")

	if !ok {
		for _, command := range s.diskUnmountCommands {
			_, err := driver.ExecuteOsaScript(command...)
			if err != nil {
				log.Printf("error detaching iso: %s", err)
			}
		}
	}
}
