Dependencies are explicit and acyclic:

```text
config -> image -> kernel
vm -> resource
vm -> volume
vm -> image
vm -> kernel
```

## Sources of truth

- `rune.yaml` describes desired kernels, immutable images, and runtime policy.
- Image `metadata.yaml` describes a completed immutable artifact.
- Volume `metadata.yaml` describes a writable filesystem and its owner.
- VM `metadata.yaml` describes provider, kernel, root volume, CPU, and memory.
- Firecracker `config.json` is generated output, never authoritative state.

All metadata and generated configuration writes are atomic. Mutating image,
kernel, volume, and VM operations use per-object locks under the Rune state
directory.

## Integration boundary

An external workflow service should call CLI commands and consume their exit
codes. It should use stable identifiers for images, VMs, and jobs:

```bash
rune config validate
rune image ensure cpp-worker:v1
rune vm launch run-018f --image cpp-worker:v1 --cpu 4 --memory 4096
rune vm enqueue run-018f --id build -- bash -lc 'cmake -S . -B build && cmake --build build'
rune vm start run-018f
rune vm job-logs run-018f build
rune vm stop run-018f
rune vm delete run-018f
```

