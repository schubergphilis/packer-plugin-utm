// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// UTM version to guest additions version map
var additionsVersionMap = map[string]string{
	"4.6.4": "0.229.2",
	"4.7.4": "0.1.271",
}

type guestAdditionsUrlTemplate struct {
	Version string
}

// This step uploads a file containing the UTM version, which
// can be useful for various provisioning reasons.
//
// Produces:
//
//	guest_additions_path string - Path to the guest additions.
type StepDownloadGuestAdditions struct {
	GuestAdditionsMode   string
	GuestAdditionsURL    string
	GuestAdditionsSHA256 string
	Ctx                  interpolate.Context
}

func (s *StepDownloadGuestAdditions) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)

	if s.skipGuestAdditionsDownload() {
		log.Println("Not downloading guest additions since it is disabled.")
		return multistep.ActionContinue
	}

	version, halt := s.fetchUTMVersion(state, driver)
	if halt != nil {
		return multistep.ActionHalt
	}

	url, checksumType, halt := s.prepareGuestAdditionsURL(state, driver, ui, version)
	if halt != nil {
		return multistep.ActionHalt
	}

	checksum, haltAction := s.prepareGuestAdditionsChecksum(ctx, state, checksumType, version)
	if haltAction != multistep.ActionContinue {
		return haltAction
	}

	return s.downloadGuestAdditions(ctx, state, url, checksum)
}

func (s *StepDownloadGuestAdditions) Cleanup(state multistep.StateBag) {}

func (s *StepDownloadGuestAdditions) downloadAdditionsSHA256(
	ctx context.Context,
	state multistep.StateBag,
	additionsVersion string,
	additionsName string,
) (string, multistep.StepAction) {
	// UTM does not provide a SHA256 checksum for the guest additions
	// See https://github.com/utmapp/qemu/releases/tag/v10.0.2-utm
	// The checksum of latest version has been hardcoded.
	// This is a temporary solution until UTM provides the checksum
	// This is the SHA256 checksum of the guest additions for UTM 4.7.4
	checksum := "65b6a69b392ee01dd314c10f3dad9ebbf9c4160be43f5f0dd6bb715944d9095b"

	return checksum, multistep.ActionContinue
}

func (s *StepDownloadGuestAdditions) skipGuestAdditionsDownload() bool {
	return s.GuestAdditionsMode == GuestAdditionsModeDisable
}

func (s *StepDownloadGuestAdditions) fetchUTMVersion(state multistep.StateBag, driver Driver) (string, error) {
	version, err := driver.Version()
	if err != nil {
		state.Put("error", fmt.Errorf("error reading version for guest additions download: %s", err))
		return "", err
	}
	if newVersion, ok := additionsVersionMap[version]; ok {
		log.Printf("Rewriting guest additions version: %s to %s", version, newVersion)
		version = newVersion
	}
	return version, nil
}

func (s *StepDownloadGuestAdditions) prepareGuestAdditionsURL(
	state multistep.StateBag,
	driver Driver,
	ui packersdk.Ui,
	version string,
) (string, string, error) {
	// Prepare the template context for interpolation
	s.Ctx.Data = &guestAdditionsUrlTemplate{
		Version: version,
	}

	url, err := interpolate.Render(s.GuestAdditionsURL, &s.Ctx)
	if err != nil {
		prepErr := fmt.Errorf("error preparing guest additions url: %s", err)
		state.Put("error", prepErr)
		ui.Error(prepErr.Error())
		return "", "", prepErr
	}
	checksumType := "sha256"

	// Fallback to driver or default remote URL if necessary
	if url == "" {
		url, err = driver.GuestToolsIsoPath()

		if err == nil {
			checksumType = "none"
		} else {
			ui.Error(err.Error())
			url = fmt.Sprintf(
				"https://getutm.app/downloads/%s", fmt.Sprintf("utm-guest-tools-%s.iso", "latest"))
		}
	}

	if url == "" {
		err := fmt.Errorf("could not detect guest additions URL.\n" +
			"Please specify `guest_additions_url` manually")
		state.Put("error", err)
		ui.Error(err.Error())
		return "", "", err
	}

	return url, checksumType, nil
}

func (s *StepDownloadGuestAdditions) prepareGuestAdditionsChecksum(
	ctx context.Context,
	state multistep.StateBag,
	checksumType string,
	version string,
) (string, multistep.StepAction) {
	var checksum string

	if checksumType != "none" {
		if s.GuestAdditionsSHA256 != "" {
			checksum = s.GuestAdditionsSHA256
		} else {
			additionsName := fmt.Sprintf("utm-guest-tools-%s.iso", "latest")
			var action multistep.StepAction
			checksum, action = s.downloadAdditionsSHA256(ctx, state, version, additionsName)
			if action != multistep.ActionContinue {
				return "", action
			}
		}
	}
	return checksum, multistep.ActionContinue
}

func (s *StepDownloadGuestAdditions) downloadGuestAdditions(
	ctx context.Context,
	state multistep.StateBag,
	url string,
	checksum string,
) multistep.StepAction {
	downStep := &commonsteps.StepDownload{
		Checksum:    checksum,
		Description: "Guest additions",
		ResultKey:   "guest_additions_path",
		Url:         []string{url},
		Extension:   "iso",
	}
	return downStep.Run(ctx, state)
}
