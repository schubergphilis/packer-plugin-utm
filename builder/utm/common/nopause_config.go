// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown

package common

import (
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// Temporary configuration to disable the pause step
// Once the pause step is removed, this configuration will be removed as well
type NoPauseConfig struct {
	// If true, the build process will not pause to prompt the user to confirm
	DisplayNoPause bool `mapstructure:"display_nopause" required:"false"`
	BootNoPause    bool `mapstructure:"boot_nopause" required:"false"`
	ExportNoPause  bool `mapstructure:"export_nopause" required:"false"`
}

func (c *NoPauseConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	return errs
}
