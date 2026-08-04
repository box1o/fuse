# Rune

Rune is a modular Bash runtime for building cached Linux images and running
isolated jobs in Firecracker microVMs. It is designed as an execution primitive
for CI, build, rendering, and other workflow platforms.

The Go program under `cmd/` is separate. Rune does not require Go for module
discovery, image construction, VM lifecycle, or guest jobs.

## Quick start

Install development and runtime tools on Arch Linux:

```bash
sudo pacman -S bats shellcheck shfmt go-yq curl debootstrap e2fsprogs
make link
```

Validate a configuration and inspect its image catalog:

```bash
RUNE_CONFIG_FILE=examples/rune-images.yaml rune config validate
RUNE_CONFIG_FILE=examples/rune-images.yaml rune image catalog
```

Build one image and launch a VM:

```bash
sudo RUNE_CONFIG_FILE=examples/rune-images.yaml rune image build cpp-worker
RUNE_CONFIG_FILE=examples/rune-images.yaml \
    rune vm launch builder --image cpp-worker:v1 --cpu 4 --memory 4096
```

Queue and observe a guest job:

```bash
rune vm stop builder
rune vm enqueue builder --id compile -- bash -lc 'cmake -S . -B build && cmake --build build'
rune vm start builder
rune vm wait-job builder compile --timeout 1800
rune vm job-logs builder compile
```

## Documentation

- [Architecture](docs/architecture.md)
- [Project configuration](docs/configuration.md)
- [Guest jobs and logs](docs/jobs.md)

