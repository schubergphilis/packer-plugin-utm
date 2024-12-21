//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package cloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	utmcommon "github.com/naveenrajm7/packer-plugin-utm/builder/utm/common"
)

// Config is the configuration structure for the UTM Cloud builder.
type Config struct {
	common.PackerConfig        `mapstructure:",squash"`
	commonsteps.HTTPConfig     `mapstructure:",squash"`
	commonsteps.ISOConfig      `mapstructure:",squash"`
	utmcommon.ExportConfig     `mapstructure:",squash"`
	utmcommon.OutputConfig     `mapstructure:",squash"`
	utmcommon.ShutdownConfig   `mapstructure:",squash"`
	utmcommon.CommConfig       `mapstructure:",squash"`
	utmcommon.HWConfig         `mapstructure:",squash"`
	utmcommon.UtmVersionConfig `mapstructure:",squash"`
	utmcommon.UtmBundleConfig  `mapstructure:",squash"`

	// Set this to true if you would like to keep the VM registered with
	// UTM. Defaults to false.
	KeepRegistered bool `mapstructure:"keep_registered" required:"false"`
	// Defaults to false. When enabled, Packer will not export the VM. Useful
	// if the build output is not the resultant image, but created inside the
	// VM.
	SkipExport bool `mapstructure:"skip_export" required:"false"`
	// QEMU system architecture of the virtual machine.
	// For a QEMU virtual machine, you must specify the architecture
	// Which is required in confirguration. By default, this is aarch64.
	// You should use same architecture as the cloud image.
	VMArch string `mapstructure:"vm_arch" required:"false"`
	// Backend to use for the virtual machine.
	// Only qemu cloud images are supported.
	// By default, this is qemu.
	VMBackend string `mapstructure:"vm_backend" required:"false"`
	// This is the name of the utm file for the new virtual machine, without
	// the file extension. By default this is packer-BUILDNAME, where
	// "BUILDNAME" is the name of the build.
	VMName string `mapstructure:"vm_name" required:"false"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(c, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	// Accumulate any errors and warnings
	var errs *packersdk.MultiError
	warnings := make([]string, 0)

	isoWarnings, isoErrs := c.ISOConfig.Prepare(&c.ctx)
	warnings = append(warnings, isoWarnings...)
	errs = packersdk.MultiErrorAppend(errs, isoErrs...)

	errs = packersdk.MultiErrorAppend(errs, c.ExportConfig.Prepare(&c.ctx)...)
	errs = packersdk.MultiErrorAppend(
		errs, c.OutputConfig.Prepare(&c.ctx, &c.PackerConfig)...)
	errs = packersdk.MultiErrorAppend(errs, c.HTTPConfig.Prepare(&c.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, c.ShutdownConfig.Prepare(&c.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, c.CommConfig.Prepare(&c.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, c.UtmBundleConfig.Prepare(&c.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, c.UtmVersionConfig.Prepare(c.CommConfig.Comm.Type)...)

	if c.VMArch == "" {
		c.VMArch = "aarch64"
	}

	if c.VMBackend == "" {
		c.VMBackend = "qemu"
	}
	// Validate and use Enums for the VM backend
	// Only qemu cloud images are supported.
	switch c.VMBackend {
	case "qemu":
		c.VMBackend = "QeMu"
	default:
		errs = packersdk.MultiErrorAppend(
			errs, errors.New("vm_backend must be 'qemu'"))
	}

	if c.VMName == "" {
		c.VMName = fmt.Sprintf(
			"packer-%s-%d", c.PackerBuildName, interpolate.InitTime.Unix())
	}

	// Warnings
	if c.ShutdownCommand == "" {
		warnings = append(warnings,
			"A shutdown_command was not specified. Without a shutdown command, Packer\n"+
				"will forcibly halt the virtual machine, which may result in data loss.")
	}

	if errs != nil && len(errs.Errors) > 0 {
		return warnings, errs
	}

	return warnings, nil

}
