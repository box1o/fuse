# Image module

The image module builds immutable Debian root filesystems from `rune.yaml`.
Image names, packages, capabilities, files, and environment values come only
from user configuration.

Every new image includes the reusable Rune shell library, `rune-guest`, and the
boot-time job queue. See `docs/configuration.md` for the full schema and
multi-image workflow.

```bash
rune image validate
rune image build custom-worker --dry-run
sudo rune image build custom-worker
rune image build-all --dry-run
rune image ensure-all
rune image catalog
rune image list
rune image show custom-worker:v1
rune image path custom-worker:v1
```

Every image includes a minimal boot baseline: BusyBox, systemd init, udev, CA
certificates, curl, and basic network tools. Requested packages are added and
deduplicated. The result is an ext4 image with immutable metadata and SHA-256.

Building uses `debootstrap` and requires root. Validation and dry-run do not.
