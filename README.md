# Jsonnet Tool

`jsonnet-tool` is a simple binary, primarily used to build YAML configuration files from Jsonnet source configurations.

## Commands

### `jsonnet-tool yaml`

```console
$ cat file.jsonnet
{
  'file.yaml': std.manifestYamlDoc({
    hello: true,
    there: 1,
    moo: {
      there: 1,
      hello: true,
    },
  }),
}
$ jsonnet-tool yaml \
    --multi "./output" \ #       - Directory to emit the YAML file to
    --header "# DO NOT EDIT" \ # - Header to prefix to output YAML
    -J "./libsonnet/" \ #        - Jsonnet Import Search Path
    -J "./vendor/" \ #           - .. supports multiple
    -P name \ #                  - Keys to appear at the top of YAML
    -P alert \ #                 - .. supports multiple
    --prefix "autogenerated-" \  - Prefix added to file names
    file.jsonnet
```

### `jsonnet-tool render`

Render is a generic rendering utility for jsonnet. In the case of JSON and YAML, the output does not need to be manifested, the tool will use
the extension of the file to appropriately manifest the output.

```console
$ cat file.jsonnet
{
  // File will contain YAML output
  'file.yaml': {
    hello: true,
    there: 1,
    moo: {
      there: 1,
      hello: true,
    },
  },

  // Subdirectories are automatically created
  'x/y/z/file.json': {
    hello: 1,
    x: [1, 2, 3],
  },
}
$ jsonnet-tool render \
    --multi "./output" \ #       - Directory to emit the YAML file to
    -J "./libsonnet/" \ #        - Jsonnet Import Search Path
    -J "./vendor/" \ #           - .. supports multiple
    --prefix "autogenerated-" \  - Prefix added to file names
    file.jsonnet
```

### `jsonnet-tool checksum`

Checksum will recursively parse imports in a jsonnet file and generate a full import tree, then output sha256 hashes for the entire tree
in a format that can be consumed by the sha256sum utility.

This is useful for caching jsonnet files.

```console
$ cat file.jsonnet
local stageGroupMapping = (import 'gitlab-metrics-config.libsonnet').stageGroupMapping;
local serviceCatalog = import 'service-catalog/service-catalog.libsonnet';
{}

$ jsonnet-tool checksum \
    -J . -J ../libsonnet -J ../metrics-catalog/ -J ../vendor -J ../services \
    file.jsonnet >file.jsonnet.sha256sum

$ cat file.jsonnet.sha256sum
0007482aebcdf09ed0798a499118f31d40ca244c71bd89a9ff9c947e26ed71a0  ../libsonnet/saturation-monitoring/shard_cpu.libsonnet
047b0b518845b0c38f6775a4e3040c08c3ce2bbfcdb870d9f30e3b190d1917e5  ../libsonnet/elasticlinkbuilder/index_catalog.libsonnet
0599206409b75d4bd1c896ee2b865166eb516ab23bcf45d0fa60a398690b0692  ../metrics-catalog/services/woodhouse.jsonnet
...

$ sha256sum --check --status <file.jsonnet.sha256sum
```

## Examples

Check the [`examples/`](examples/) directory for examples of files suitable for `jsonnet-tool`.

## Issue creation

Please create new issues in the [public infrastructure project](https://gitlab.com/gitlab-com/gl-infra/infrastructure/-/issues/) and not in this private issue tracker.

## Project Workflow

**[CONTRIBUTING.md](CONTRIBUTING.md) is important, do not skip reading this.**
