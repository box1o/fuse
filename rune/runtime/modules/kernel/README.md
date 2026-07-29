# Kernel module

The kernel module resolves versioned guest kernels from `rune.yaml`. A kernel
may use an explicit URL or the latest matching Firecracker CI artifact.

```bash
rune kernel validate
rune kernel ensure linux-6.1 --dry-run
rune kernel ensure linux-6.1
rune kernel list
rune kernel show linux-6.1
```

Downloaded kernels are checksummed and shared by all images that reference the
same configured kernel.
