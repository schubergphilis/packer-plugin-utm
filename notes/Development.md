



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