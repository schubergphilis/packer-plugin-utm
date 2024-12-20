Artifact BuilderId: `naveenrajm7.utm.post-processor.vagrant`

The Packer UTM vagrant post-processor takes an artifact with .utm directory and creates a Vagrant box.

## Basic Example

```hcl
  # Convert machines to vagrant boxes
  post-processor "utm-vagrant" {
    compression_level = 9
    output            = "${path.root}/${var.os_name}-${var.os_version}.box"
  }
```

<!-- Post-Processor Configuration Fields -->
## Configuration Reference

<!--
  Optional Configuration Fields

  Configuration options that are not required or have reasonable defaults
  should be listed under the optionals section. Defaults values should be
  noted in the description of the field
-->

### Optional:

<!-- Code generated from the comments of the Config struct in post-processor/vagrant/post-processor.go; DO NOT EDIT MANUALLY -->

- `compression_level` (int) - Compression Level

- `include` ([]string) - Include

- `output` (string) - Output Path

- `vagrantfile_template` (string) - Vagrantfile Template

- `vagrantfile_template_generated` (bool) - Vagrantfile Template Generated

- `provider_override` (string) - Provider Override

- `architecture` (string) - Architecture

<!-- End of code generated from the comments of the Config struct in post-processor/vagrant/post-processor.go; -->
