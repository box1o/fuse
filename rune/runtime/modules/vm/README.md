# VM module

The VM module uses Firecracker to manage local microVMs under the Rune cache. It is a
development lifecycle today, not yet a production multi-tenant jailer or
networking implementation.

## Quick start

```bash
rune vm setup
rune vm create builder --image custom-worker:v1 --cpu 2 --memory 1024
rune vm start builder
rune vm status builder
rune vm list
rune vm logs builder
rune vm enqueue builder --id build -- bash -lc 'make -j4'
rune vm job-logs builder build
rune vm stop builder
rune vm resources builder
rune vm resize builder --memory 4096 --cpu 4
rune vm resize builder --disk 20G
rune vm delete builder
```

The combined workflow is:

```bash
rune vm run builder --image custom-worker:v1 --cpu 2 --memory 1024
```

`run` performs setup, ensures the configured image, creates the named VM when
it does not exist, and starts it. Existing VMs keep their original CPU, memory,
kernel, and rootfs settings.

## Storage layout

```text
~/.cache/rune/firecracker/
├── downloads/               Release archives
├── runtime/                 firecracker and jailer binaries
├── assets/                  Legacy shared assets
└── vms/<name>/
    ├── config.json
    ├── metadata.yaml        Provider configuration and root-volume reference
    ├── kernel.path
    ├── console.log
    ├── firecracker.pid
    └── firecracker.sock

~/.cache/rune/volumes/<name>/
├── volume.ext4              Per-VM writable root filesystem
└── metadata.yaml            Source, owner, size, and filesystem
```

Set `FIRECRACKER_CACHE_DIR` to move the complete layout. Set
`FIRECRACKER_KERNEL_URL` and `FIRECRACKER_ROOTFS_URL` before running `setup` to
use different base assets.

## Custom images

Create a VM from existing images:

```bash
rune vm create builder \
    --kernel /absolute/path/to/vmlinux \
    --rootfs /absolute/path/to/rootfs.ext4
```

The rootfs is copied into a managed volume and becomes writable. The kernel is
referenced read-only from its original location.

## Resources and resizing

Inspect a VM and change resources while it is stopped:

```bash
rune vm stop builder
rune vm resources builder
rune vm set-memory builder 4096
rune vm resize builder --cpu 4 --disk 20G
rune vm start builder
```

CPU and memory changes update the next boot configuration. Disk growth is
delegated to the generic `volume` module, which checks ext4 before and after
growth. Disk shrinking is not supported. Older VMs with an embedded
`rootfs.ext4` can still boot and use chroot, but disk resize requires recreating
them as managed-volume VMs.

Cleanup removes stale runtime files from stopped VMs. Logs are preserved unless
requested explicitly:

```bash
rune vm cleanup builder
rune vm cleanup builder --logs
rune vm cleanup --all
```

## Chroot

Only stopped VMs may be mounted:

```bash
sudo rune vm chroot builder
sudo rune vm chroot builder -- /bin/sh -c 'cat /etc/os-release'
sudo rune vm chroot --read-only builder
```

This operation requires root because it loop-mounts the ext4 image. Chroot is a
filesystem-root change and is not a security sandbox.

## Host requirements

Starting a VM requires read/write access to `/dev/kvm`. The current lifecycle
does not configure TAP networking. Production use additionally requires the
Firecracker jailer, unique identities, cgroups, network namespaces, firewall
rules, resource limits, and external log retention.

See `docs/jobs.md` for the guest job and log contract. See
`docs/architecture.md` for the stable integration boundary.
