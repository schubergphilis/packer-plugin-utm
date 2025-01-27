Type: `utm-cloud`
Artifact BuilderId: `naveenrajm.cloud`

The builder builds a virtual machine by importing an existing cloud image, qcow2 file.
It then boots this image, runs provisioners on this new VM, and exports that VM
to create the image. The imported machine is deleted prior to finishing the
build.

## Basic Example

Here is a basic example. This example is functional if you have an qcow2  matching
the settings here.

**HCL2**

```hcl
source "utm-cloud" "basic-example" {
  iso_url = "cloud.qcow2"
  iso_checksum = "sha256:1234567890abcdef"
  // Required to launch http server to serve cloud-init files
  http_directory =  "/path-to-cloud-file/"
  ssh_username = "vagrant"
  ssh_password = "vagrant"
  shutdown_command = "echo 'vagrant' | sudo -S /sbin/halt -h -p"
}

build {
  sources = ["sources.utm-cloud.basic-example"]
}
```

It is important to add a `shutdown_command`. By default Packer halts the virtual
machine and the file system may not be sync'd. Thus, changes made in a
provisioner might not be saved.

## Configuration Reference

There are many configuration options available for the builder. In addition to
the items listed here, you will want to look at the general configuration
references for [ISO](#iso-configuration),
[HTTP](#http-directory-configuration),
[Export](#export-configuration),
[Shutdown](#shutdown-configuration),
[Run](#run-configuration),
[Communicator](#communicator-configuration)
configuration references, which are
necessary for this build to succeed and can be found further down the page.

### Optional:

<!-- Code generated from the comments of the Config struct in builder/utm/cloud/config.go; DO NOT EDIT MANUALLY -->

- `hypervisor` (bool) - Set this to true if you would like to use Hypervisor
  Defaults to false.

- `uefi_boot` (bool) - Set this to true if you would like to use UEFI firmware to boot with
  UTM. Defaults to false.

- `disk_size` (uint) - The size, in megabytes, of the hard disk to create for the VM. By
  default, this is 40000 (about 40 GB).

- `hard_drive_interface` (string) - The type of controller that the primary hard drive is attached to,
  defaults to VirtIO. When set to usb, the drive is attached to an USB
  controller. When set to scsi, the drive is attached to an  SCSI
  controller. When set to nvme, the drive is attached to an NVMe
  controller. When set to virtio, the drive is attached to a VirtIO
  controller. Please note that when you use "nvme",
  and you may need to enable EFI mode for nvme to work (this note is from VirtualBox)

- `iso_interface` (string) - The type of controller that the ISO is attached to, defaults to usb.
  When set to nvme, the drive is attached to an NVMe controller.
  When set to virtio, the drive is attached to a VirtIO controller.

- `disk_additional_size` ([]uint) - Additional disks to create. Attachment starts at 1 since 0
  is the default disk. Each value represents the disk image size in MiB.
  Each additional disk uses the same disk parameters as the default disk.
  Unset by default.

- `resize_cloud_image` (bool) - Wheather to resize the cloud image to the disk size. Defaults to false.
  If set to true, the cloud image will be resized to the disk size.
  Required qemu-img to be installed in the system.

- `use_cd` (bool) - Pass cloud-init data to the VM using a CD-ROM. Defaults to false.
  If set to true, you must provide cd_files with the path to the cloud-init
  data and cd_label with value "cidata".
  If set to false, you must provide http_directory with the cloud-init data.

- `keep_registered` (bool) - Set this to true if you would like to keep the VM registered with
  UTM. Defaults to false.

- `skip_export` (bool) - Defaults to false. When enabled, Packer will not export the VM. Useful
  if the build output is not the resultant image, but created inside the
  VM.

- `vm_arch` (string) - QEMU system architecture of the virtual machine.
  For a QEMU virtual machine, you must specify the architecture
  Which is required in confirguration. By default, this is aarch64.
  You should use same architecture as the cloud image.

- `vm_backend` (string) - Backend to use for the virtual machine.
  Only qemu cloud images are supported.
  By default, this is qemu.

- `vm_name` (string) - This is the name of the utm file for the new virtual machine, without
  the file extension. By default this is packer-BUILDNAME, where
  "BUILDNAME" is the name of the build.

<!-- End of code generated from the comments of the Config struct in builder/utm/cloud/config.go; -->


<!-- Code generated from the comments of the UtmVersionConfig struct in builder/utm/common/utm_version_config.go; DO NOT EDIT MANUALLY -->

- `utm_version_file` (\*string) - The path within the virtual machine to
  upload a file that contains the UTM version that was used to create
  the machine. This information can be useful for provisioning. By default
  this is .utm_version, which will generally be upload it into the
  home directory. Set to an empty string to skip uploading this file, which
  can be useful when using the none communicator.

<!-- End of code generated from the comments of the UtmVersionConfig struct in builder/utm/common/utm_version_config.go; -->


### ISO Configuration

<!-- Code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; DO NOT EDIT MANUALLY -->

By default, Packer will symlink, download or copy image files to the Packer
cache into a "`hash($iso_url+$iso_checksum).$iso_target_extension`" file.
Packer uses [hashicorp/go-getter](https://github.com/hashicorp/go-getter) in
file mode in order to perform a download.

go-getter supports the following protocols:

* Local files
* Git
* Mercurial
* HTTP
* Amazon S3

Examples:
go-getter can guess the checksum type based on `iso_checksum` length, and it is
also possible to specify the checksum type.

In JSON:

```json

	"iso_checksum": "946a6077af6f5f95a51f82fdc44051c7aa19f9cfc5f737954845a6050543d7c2",
	"iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```json

	"iso_checksum": "file:ubuntu.org/..../ubuntu-14.04.1-server-amd64.iso.sum",
	"iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```json

	"iso_checksum": "file://./shasums.txt",
	"iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```json

	"iso_checksum": "file:./shasums.txt",
	"iso_url": "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

In HCL2:

```hcl

	iso_checksum = "946a6077af6f5f95a51f82fdc44051c7aa19f9cfc5f737954845a6050543d7c2"
	iso_url = "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```hcl

	iso_checksum = "file:ubuntu.org/..../ubuntu-14.04.1-server-amd64.iso.sum"
	iso_url = "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```hcl

	iso_checksum = "file://./shasums.txt"
	iso_url = "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

```hcl

	iso_checksum = "file:./shasums.txt",
	iso_url = "ubuntu.org/.../ubuntu-14.04.1-server-amd64.iso"

```

<!-- End of code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; -->


#### Required:

<!-- Code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; DO NOT EDIT MANUALLY -->

- `iso_checksum` (string) - The checksum for the ISO file or virtual hard drive file. The type of
  the checksum is specified within the checksum field as a prefix, ex:
  "md5:{$checksum}". The type of the checksum can also be omitted and
  Packer will try to infer it based on string length. Valid values are
  "none", "{$checksum}", "md5:{$checksum}", "sha1:{$checksum}",
  "sha256:{$checksum}", "sha512:{$checksum}" or "file:{$path}". Here is a
  list of valid checksum values:
   * md5:090992ba9fd140077b0661cb75f7ce13
   * 090992ba9fd140077b0661cb75f7ce13
   * sha1:ebfb681885ddf1234c18094a45bbeafd91467911
   * ebfb681885ddf1234c18094a45bbeafd91467911
   * sha256:ed363350696a726b7932db864dda019bd2017365c9e299627830f06954643f93
   * ed363350696a726b7932db864dda019bd2017365c9e299627830f06954643f93
   * file:http://releases.ubuntu.com/20.04/SHA256SUMS
   * file:file://./local/path/file.sum
   * file:./local/path/file.sum
   * none
  Although the checksum will not be verified when it is set to "none",
  this is not recommended since these files can be very large and
  corruption does happen from time to time.

- `iso_url` (string) - A URL to the ISO containing the installation image or virtual hard drive
  (VHD or VHDX) file to clone.

<!-- End of code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; -->


#### Optional:

<!-- Code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; DO NOT EDIT MANUALLY -->

- `iso_urls` ([]string) - Multiple URLs for the ISO to download. Packer will try these in order.
  If anything goes wrong attempting to download or while downloading a
  single URL, it will move on to the next. All URLs must point to the same
  file (same checksum). By default this is empty and `iso_url` is used.
  Only one of `iso_url` or `iso_urls` can be specified.

- `iso_target_path` (string) - The path where the iso should be saved after download. By default will
  go in the packer cache, with a hash of the original filename and
  checksum as its name.

- `iso_target_extension` (string) - The extension of the iso file after download. This defaults to `iso`.

<!-- End of code generated from the comments of the ISOConfig struct in multistep/commonsteps/iso_config.go; -->



### Http directory configuration

<!-- Code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; DO NOT EDIT MANUALLY -->

Packer will create an http server serving `http_directory` when it is set, a
random free port will be selected and the architecture of the directory
referenced will be available in your builder.

Example usage from a builder:

```
wget http://{{ .HTTPIP }}:{{ .HTTPPort }}/foo/bar/preseed.cfg
```

<!-- End of code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; -->


#### Optional:

<!-- Code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; DO NOT EDIT MANUALLY -->

- `http_directory` (string) - Path to a directory to serve using an HTTP server. The files in this
  directory will be available over HTTP that will be requestable from the
  virtual machine. This is useful for hosting kickstart files and so on.
  By default this is an empty string, which means no HTTP server will be
  started. The address and port of the HTTP server will be available as
  variables in `boot_command`. This is covered in more detail below.

- `http_content` (map[string]string) - Key/Values to serve using an HTTP server. `http_content` works like and
  conflicts with `http_directory`. The keys represent the paths and the
  values contents, the keys must start with a slash, ex: `/path/to/file`.
  `http_content` is useful for hosting kickstart files and so on. By
  default this is empty, which means no HTTP server will be started. The
  address and port of the HTTP server will be available as variables in
  `boot_command`. This is covered in more detail below.
  Example:
  ```hcl
    http_content = {
      "/a/b"     = file("http/b")
      "/foo/bar" = templatefile("${path.root}/preseed.cfg", { packages = ["nginx"] })
    }
  ```

- `http_port_min` (int) - These are the minimum and maximum port to use for the HTTP server
  started to serve the `http_directory`. Because Packer often runs in
  parallel, Packer will choose a randomly available port in this range to
  run the HTTP server. If you want to force the HTTP server to be on one
  port, make this minimum and maximum port the same. By default the values
  are `8000` and `9000`, respectively.

- `http_port_max` (int) - HTTP Port Max

- `http_bind_address` (string) - This is the bind address for the HTTP server. Defaults to 0.0.0.0 so that
  it will work with any network interface.

- `http_network_protocol` (string) - Defines the HTTP Network protocol. Valid options are `tcp`, `tcp4`, `tcp6`,
  `unix`, and `unixpacket`. This value defaults to `tcp`.

<!-- End of code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; -->


### CD configuration

<!-- Code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; DO NOT EDIT MANUALLY -->

An iso (CD) containing custom files can be made available for your build.

By default, no extra CD will be attached. All files listed in this setting
get placed into the root directory of the CD and the CD is attached as the
second CD device.

This config exists to work around modern operating systems that have no
way to mount floppy disks, which was our previous go-to for adding files at
boot time.

<!-- End of code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; -->


#### Optional:

<!-- Code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; DO NOT EDIT MANUALLY -->

- `cd_files` ([]string) - A list of files to place onto a CD that is attached when the VM is
  booted. This can include either files or directories; any directories
  will be copied onto the CD recursively, preserving directory structure
  hierarchy. Symlinks will have the link's target copied into the directory
  tree on the CD where the symlink was. File globbing is allowed.
  
  Usage example (JSON):
  
  ```json
  "cd_files": ["./somedirectory/meta-data", "./somedirectory/user-data"],
  "cd_label": "cidata",
  ```
  
  Usage example (HCL):
  
  ```hcl
  cd_files = ["./somedirectory/meta-data", "./somedirectory/user-data"]
  cd_label = "cidata"
  ```
  
  The above will create a CD with two files, user-data and meta-data in the
  CD root. This specific example is how you would create a CD that can be
  used for an Ubuntu 20.04 autoinstall.
  
  Since globbing is also supported,
  
  ```hcl
  cd_files = ["./somedirectory/*"]
  cd_label = "cidata"
  ```
  
  Would also be an acceptable way to define the above cd. The difference
  between providing the directory with or without the glob is whether the
  directory itself or its contents will be at the CD root.
  
  Use of this option assumes that you have a command line tool installed
  that can handle the iso creation. Packer will use one of the following
  tools:
  
    * xorriso
    * mkisofs
    * hdiutil (normally found in macOS)
    * oscdimg (normally found in Windows as part of the Windows ADK)

- `cd_content` (map[string]string) - Key/Values to add to the CD. The keys represent the paths, and the values
  contents. It can be used alongside `cd_files`, which is useful to add large
  files without loading them into memory. If any paths are specified by both,
  the contents in `cd_content` will take precedence.
  
  Usage example (HCL):
  
  ```hcl
  cd_files = ["vendor-data"]
  cd_content = {
    "meta-data" = jsonencode(local.instance_data)
    "user-data" = templatefile("user-data", { packages = ["nginx"] })
  }
  cd_label = "cidata"
  ```

- `cd_label` (string) - CD Label

<!-- End of code generated from the comments of the CDConfig struct in multistep/commonsteps/extra_iso_config.go; -->


### Export configuration

#### Optional:

<!-- Code generated from the comments of the ExportConfig struct in builder/utm/common/export_config.go; DO NOT EDIT MANUALLY -->

- `format` (string) - Only UTM, this specifies the output format
  of the exported virtual machine. This defaults to utm.

<!-- End of code generated from the comments of the ExportConfig struct in builder/utm/common/export_config.go; -->




### Shutdown configuration

#### Optional:

<!-- Code generated from the comments of the ShutdownConfig struct in builder/utm/common/shutdown_config.go; DO NOT EDIT MANUALLY -->

- `shutdown_command` (string) - The command to use to gracefully shut down the
  machine once all the provisioning is done. By default this is an empty
  string, which tells Packer to just forcefully shut down the machine unless a
  shutdown command takes place inside script so this may safely be omitted. If
  one or more scripts require a reboot it is suggested to leave this blank
  since reboots may fail and specify the final shutdown command in your
  last script.

- `shutdown_timeout` (duration string | ex: "1h5m2s") - The amount of time to wait after executing the
  shutdown_command for the virtual machine to actually shut down. If it
  doesn't shut down in this time, it is an error. By default, the timeout is
  5m or five minutes.

- `post_shutdown_delay` (duration string | ex: "1h5m2s") - The amount of time to wait after shutting
  down the virtual machine. If you get the error
  Error removing floppy controller, you might need to set this to 5m
  or so. By default, the delay is 0s or disabled.

- `disable_shutdown` (bool) - Packer normally halts the virtual machine after all provisioners have
  run when no `shutdown_command` is defined.  If this is set to `true`, Packer
  *will not* halt the virtual machine but will assume that you will send the stop
  signal yourself through the preseed.cfg or your final provisioner.
  Packer will wait for a default of 5 minutes until the virtual machine is shutdown.
  The timeout can be changed using `shutdown_timeout` option.

<!-- End of code generated from the comments of the ShutdownConfig struct in builder/utm/common/shutdown_config.go; -->


### Hardware configuration

#### Optional:

<!-- Code generated from the comments of the HWConfig struct in builder/utm/common/hw_config.go; DO NOT EDIT MANUALLY -->

- `cpus` (int) - The number of cpus to use for building the VM.
  Defaults to 1.

- `memory` (int) - The amount of memory to use for building the VM
  in megabytes. Defaults to 512 megabytes.

<!-- End of code generated from the comments of the HWConfig struct in builder/utm/common/hw_config.go; -->


### Communicator configuration

#### Optional common fields:

<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

- `communicator` (string) - Packer currently supports three kinds of communicators:
  
  -   `none` - No communicator will be used. If this is set, most
      provisioners also can't be used.
  
  -   `ssh` - An SSH connection will be established to the machine. This
      is usually the default.
  
  -   `winrm` - A WinRM connection will be established.
  
  In addition to the above, some builders have custom communicators they
  can use. For example, the Docker builder has a "docker" communicator
  that uses `docker exec` and `docker cp` to execute scripts and copy
  files.

- `pause_before_connecting` (duration string | ex: "1h5m2s") - We recommend that you enable SSH or WinRM as the very last step in your
  guest's bootstrap script, but sometimes you may have a race condition
  where you need Packer to wait before attempting to connect to your
  guest.
  
  If you end up in this situation, you can use the template option
  `pause_before_connecting`. By default, there is no pause. For example if
  you set `pause_before_connecting` to `10m` Packer will check whether it
  can connect, as normal. But once a connection attempt is successful, it
  will disconnect and then wait 10 minutes before connecting to the guest
  and beginning provisioning.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


<!-- Code generated from the comments of the CommConfig struct in builder/utm/common/comm_config.go; DO NOT EDIT MANUALLY -->

- `host_port_min` (int) - The minimum port to use for the Communicator port on the host machine which is forwarded
  to the SSH or WinRM port on the guest machine. By default this is 2222.

- `host_port_max` (int) - The maximum port to use for the Communicator port on the host machine which is forwarded
  to the SSH or WinRM port on the guest machine. Because Packer often runs in parallel,
  Packer will choose a randomly available port in this range to use as the
  host port. By default this is 4444.

- `skip_nat_mapping` (bool) - Defaults to false. When enabled, Packer
  does not setup forwarded port mapping for communicator (SSH or WinRM) requests and uses ssh_port or winrm_port
  on the host to communicate to the virtual machine.

<!-- End of code generated from the comments of the CommConfig struct in builder/utm/common/comm_config.go; -->


### No Pause configuration

<!-- Code generated from the comments of the NoPauseConfig struct in builder/utm/common/nopause_config.go; DO NOT EDIT MANUALLY -->

Temporary configuration to disable the pause step
Once the pause step is removed, this configuration will be removed as well

<!-- End of code generated from the comments of the NoPauseConfig struct in builder/utm/common/nopause_config.go; -->


#### Optional:

<!-- Code generated from the comments of the NoPauseConfig struct in builder/utm/common/nopause_config.go; DO NOT EDIT MANUALLY -->

- `display_nopause` (bool) - If true, the build process will not pause to add display.
  false by default

- `boot_nopause` (bool) - If true, the build process will not pause to confirm successful boot.
  false by default

- `export_nopause` (bool) - If true, the build process will not pause to allow pre-export steps.
  false by default

<!-- End of code generated from the comments of the NoPauseConfig struct in builder/utm/common/nopause_config.go; -->
