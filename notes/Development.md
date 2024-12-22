



Generates hcl2spec configs
```Go
//go:generate packer-sdc mapstructure-to-hcl2 -type Config
```

## Docs

Generates docs 
```Go
//go:generate packer-sdc struct-markdown
```

Add comments before configs and properties

Update mdx files in /docs

Run
```make generate```

## Release

* Update plugin version at `version/version.go`
* `go mod tidy`
* make test
* `git add go.mod go.sum`
* `git commit -m "Update example.com/package to vX.Y.Z"`
* `git tag vX.Y.Z`
* push git tag
* GHA will create a release with binaries