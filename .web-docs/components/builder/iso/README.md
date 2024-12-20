Type: `utm-iso`
Artifact BuilderId: `naveenrajm7.iso`

The UTM Packer builder is able to create
[UTM](https://mac.getutm.app/) virtual machines and export them in
the .utm format., starting from an existing ISO file.

The builder builds a virtual machine by creating 
a new virtual machine from scratch, booting it, installing an OS, 
provisioning software within the OS, then shutting it down. 
The result of the UTM builder is a directory containing 
all the files necessary to run the virtual machine portably.

<!--
  A basic example on the usage of the builder. Multiple examples
  can be provided to highlight various build configurations.
-->
### Basic Example

Here is a basic example. This example is not functional. 
It will start the OS installer but then fail because we don't provide the install config file for OpenBSD to self-install. 
Still, the example serves to show the basic configuration:

```hcl
source "utm-iso" "basic-example" {
  iso_url = "https://cdn.openbsd.org/pub/OpenBSD/7.6/arm64/install76.iso"
  iso_checksum = "file:https://cdn.openbsd.org/pub/OpenBSD/7.6/arm64/SHA256"
  vm_backend = "qemu"
  ssh_username = "vagrant"
  ssh_password = "vagrant"
  # Assumes ssh user (vagrant) exists and can doas without password
  shutdown_command = "doas -u root shutdown -p now"
}

build {
  sources = [ "source.utm-iso.basic-example" ]
}
```

It is important to add a `shutdown_command`. By default Packer halts the virtual
machine and the file system may not be sync'd. Thus, changes made in a
provisioner might not be saved.

<!-- Builder Configuration Fields -->
## UTM-ISO Builder Configuration Reference

There are many configuration options available for the builder. In addition to
the items listed here, you will want to look at the general configuration
references for [ISO](#iso-configuration),
[HTTP](#http-directory-configuration),
[Export](#export-configuration),
[Boot](#boot-configuration),
[Shutdown](#shutdown-configuration),
[Run](#run-configuration),
[Communicator](#communicator-configuration)
configuration references, which are
necessary for this build to succeed and can be found further down the page.


### Optional:

<!-- Code generated from the comments of the Config struct in builder/utm/iso/config.go; DO NOT EDIT MANUALLY -->

- `boot_steps` ([][]string) - This is an array of tuples of boot commands, to type when the virtual
  machine is booted. The first element of the tuple is the actual boot
  command. The second element of the tuple, which is optional, is a
  description of what the boot command does. This is intended to be used for
  interactive installers that requires many commands to complete the
  installation. Both the command and the description will be printed when
  logging is enabled. When debug mode is enabled Packer will pause after
  typing each boot command. This will make it easier to follow along the
  installation process and make sure the Packer and the installer are in
  sync. `boot_steps` and `boot_commands` are mutually exclusive.
  
  Example:
  
  In HCL:
  ```hcl
  boot_steps = [
    ["1<enter><wait5>", "Install NetBSD"],
    ["a<enter><wait5>", "Installation messages in English"],
    ["a<enter><wait5>", "Keyboard type: unchanged"],
  
    ["a<enter><wait5>", "Install NetBSD to hard disk"],
    ["b<enter><wait5>", "Yes"]
  ]
  ```
  
  In JSON:
  ```json
  {
    "boot_steps": [
      ["1<enter><wait5>", "Install NetBSD"],
      ["a<enter><wait5>", "Installation messages in English"],
      ["a<enter><wait5>", "Keyboard type: unchanged"],
  
      ["a<enter><wait5>", "Install NetBSD to hard disk"],
      ["b<enter><wait5>", "Yes"]
    ]
  }
  ```

- `disk_size` (uint) - The size, in megabytes, of the hard disk to create for the VM. By
  default, this is 40000 (about 40 GB).

- `keep_registered` (bool) - Set this to true if you would like to keep the VM registered with
  UTM. Defaults to false.

- `skip_export` (bool) - Defaults to false. When enabled, Packer will not export the VM. Useful
  if the build output is not the resultant image, but created inside the
  VM.

- `vnc_bind_address` (string) - The IP address that should be
  binded to for VNC. By default packer will use 127.0.0.1 for this. If you
  wish to bind to all interfaces use 0.0.0.0.

- `vnc_use_password` (bool) - Whether or not to set a password on the VNC server. This option
  automatically enables the QMP socket. See `qmp_socket_path`. Defaults to
  `false`.

- `vnc_port_min` (int) - The minimum and maximum port
  to use for VNC access to the virtual machine. The builder uses VNC to type
  the initial boot_command. Because Packer generally runs in parallel,
  Packer uses a randomly chosen port in this range that appears available. By
  default this is 5900 to 6000. The minimum and maximum ports are inclusive.
  The minimum port cannot be set below 5900 due to a quirk in how QEMU parses
  vnc display address.

- `vnc_port_max` (int) - VNC Port Max

- `vm_arch` (string) - QEMU system architecture of the virtual machine.
  If this is a QEMU virtual machine, you must specify the architecture
  Which is required in confirguration. By default, this is aarch64.

- `vm_backend` (string) - Backend to use for the virtual machine.
  apple : Apple Virtualization.framework backend.
  qemu : QEMU backend.
  By default, this is qemu.

- `vm_name` (string) - This is the name of the utm file for the new virtual machine, without
  the file extension. By default this is packer-BUILDNAME, where
  "BUILDNAME" is the name of the build.

<!-- End of code generated from the comments of the Config struct in builder/utm/iso/config.go; -->


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

<!-- End of code generated from the comments of the HTTPConfig struct in multistep/commonsteps/http_config.go; -->




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




### Boot Configuration

<!-- Code generated from the comments of the BootConfig struct in bootcommand/config.go; DO NOT EDIT MANUALLY -->

The boot configuration is very important: `boot_command` specifies the keys
to type when the virtual machine is first booted in order to start the OS
installer. This command is typed after boot_wait, which gives the virtual
machine some time to actually load.

The boot_command is an array of strings. The strings are all typed in
sequence. It is an array only to improve readability within the template.

There are a set of special keys available. If these are in your boot
command, they will be replaced by the proper key:

-   `<bs>` - Backspace

-   `<del>` - Delete

-   `<enter> <return>` - Simulates an actual "enter" or "return" keypress.

-   `<esc>` - Simulates pressing the escape key.

-   `<tab>` - Simulates pressing the tab key.

-   `<f1> - <f12>` - Simulates pressing a function key.

-   `<up> <down> <left> <right>` - Simulates pressing an arrow key.

-   `<spacebar>` - Simulates pressing the spacebar.

-   `<insert>` - Simulates pressing the insert key.

-   `<home> <end>` - Simulates pressing the home and end keys.

  - `<pageUp> <pageDown>` - Simulates pressing the page up and page down
    keys.

-   `<menu>` - Simulates pressing the Menu key.

-   `<leftAlt> <rightAlt>` - Simulates pressing the alt key.

-   `<leftCtrl> <rightCtrl>` - Simulates pressing the ctrl key.

-   `<leftShift> <rightShift>` - Simulates pressing the shift key.

-   `<leftSuper> <rightSuper>` - Simulates pressing the ⌘ or Windows key.

  - `<wait> <wait5> <wait10>` - Adds a 1, 5 or 10 second pause before
    sending any additional keys. This is useful if you have to generally
    wait for the UI to update before typing more.

  - `<waitXX>` - Add an arbitrary pause before sending any additional keys.
    The format of `XX` is a sequence of positive decimal numbers, each with
    optional fraction and a unit suffix, such as `300ms`, `1.5h` or `2h45m`.
    Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. For
    example `<wait10m>` or `<wait1m20s>`.

  - `<XXXOn> <XXXOff>` - Any printable keyboard character, and of these
    "special" expressions, with the exception of the `<wait>` types, can
    also be toggled on or off. For example, to simulate ctrl+c, use
    `<leftCtrlOn>c<leftCtrlOff>`. Be sure to release them, otherwise they
    will be held down until the machine reboots. To hold the `c` key down,
    you would use `<cOn>`. Likewise, `<cOff>` to release.

  - `{{ .HTTPIP }} {{ .HTTPPort }}` - The IP and port, respectively of an
    HTTP server that is started serving the directory specified by the
    `http_directory` configuration parameter. If `http_directory` isn't
    specified, these will be blank!

-   `{{ .Name }}` - The name of the VM.

Example boot command. This is actually a working boot command used to start an
CentOS 6.4 installer:

In JSON:

```json
"boot_command": [

	   "<tab><wait>",
	   " ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos6-ks.cfg<enter>"
	]

```

In HCL2:

```hcl
boot_command = [

	   "<tab><wait>",
	   " ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos6-ks.cfg<enter>"
	]

```

The example shown below is a working boot command used to start an Ubuntu
12.04 installer:

In JSON:

```json
"boot_command": [

	"<esc><esc><enter><wait>",
	"/install/vmlinuz noapic ",
	"preseed/url=http://{{ .HTTPIP }}:{{ .HTTPPort }}/preseed.cfg ",
	"debian-installer=en_US auto locale=en_US kbd-chooser/method=us ",
	"hostname={{ .Name }} ",
	"fb=false debconf/frontend=noninteractive ",
	"keyboard-configuration/modelcode=SKIP keyboard-configuration/layout=USA ",
	"keyboard-configuration/variant=USA console-setup/ask_detect=false ",
	"initrd=/install/initrd.gz -- <enter>"

]
```

In HCL2:

```hcl
boot_command = [

	"<esc><esc><enter><wait>",
	"/install/vmlinuz noapic ",
	"preseed/url=http://{{ .HTTPIP }}:{{ .HTTPPort }}/preseed.cfg ",
	"debian-installer=en_US auto locale=en_US kbd-chooser/method=us ",
	"hostname={{ .Name }} ",
	"fb=false debconf/frontend=noninteractive ",
	"keyboard-configuration/modelcode=SKIP keyboard-configuration/layout=USA ",
	"keyboard-configuration/variant=USA console-setup/ask_detect=false ",
	"initrd=/install/initrd.gz -- <enter>"

]
```

For more examples of various boot commands, see the sample projects from our
[community templates page](https://packer.io/community-tools#templates).

<!-- End of code generated from the comments of the BootConfig struct in bootcommand/config.go; -->


#### Optional:

<!-- Code generated from the comments of the BootConfig struct in bootcommand/config.go; DO NOT EDIT MANUALLY -->

- `boot_keygroup_interval` (duration string | ex: "1h5m2s") - Time to wait after sending a group of key pressses. The value of this
  should be a duration. Examples are `5s` and `1m30s` which will cause
  Packer to wait five seconds and one minute 30 seconds, respectively. If
  this isn't specified, a sensible default value is picked depending on
  the builder type.

- `boot_wait` (duration string | ex: "1h5m2s") - The time to wait after booting the initial virtual machine before typing
  the `boot_command`. The value of this should be a duration. Examples are
  `5s` and `1m30s` which will cause Packer to wait five seconds and one
  minute 30 seconds, respectively. If this isn't specified, the default is
  `10s` or 10 seconds. To set boot_wait to 0s, use a negative number, such
  as "-1s"

- `boot_command` ([]string) - This is an array of commands to type when the virtual machine is first
  booted. The goal of these commands should be to type just enough to
  initialize the operating system installer. Special keys can be typed as
  well, and are covered in the section below on the boot command. If this
  is not specified, it is assumed the installer will start itself.

<!-- End of code generated from the comments of the BootConfig struct in bootcommand/config.go; -->
