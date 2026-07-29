# Project configuration

Rune reads a versioned YAML document. The default is `rune.yaml`; select another
file with `RUNE_CONFIG_FILE` or a command's `--file` option.

```bash
rune config path
rune config validate
rune config get '.images | keys'
rune image catalog
```

A complete three-image catalog is available at `examples/rune-images.yaml`.

## Kernel definition

```yaml
kernels:
  linux-6.1:
    version: "6.1"
    architecture: x86_64
    source: firecracker-ci
```

Supported sources are `firecracker-ci` and `url`. URL sources also require a
valid `url` field.

## Image definition

```yaml
images:
  cpp-worker:
    version: v1
    architecture: x86_64
    size: 6G
    system:
      distribution: debian
      release: bookworm
    kernel: linux-6.1
    packages:
      - build-essential
      - clang
      - cmake
      - ninja-build
    capabilities:
      - ci
      - cpp
    environment:
      RUNE_WORKER_TYPE: cpp
    files:
      - source: files/company-ca.crt
        destination: /usr/local/share/ca-certificates/company-ca.crt
        mode: "0644"
```

File sources are resolved relative to the YAML file. Guest destinations must be
absolute and cannot traverse parent directories.

Every image automatically includes the baseline boot/network packages,
BusyBox, the Rune shell library, `rune-guest`, and the boot-time job queue.

## Cache workflow

```bash
sudo RUNE_CONFIG_FILE=examples/rune-images.yaml rune image build-all
rune image catalog
rune vm launch builder --image cpp-worker:v1 --cpu 4 --memory 4096
```

`build-all` builds every declared image. `ensure-all` builds only missing
artifacts. VM creation clones the cached image into a writable managed volume;
the immutable image remains unchanged.

Changing an image definition should use a new version. Rune deliberately does
not infer content hashes as human-facing versions.
