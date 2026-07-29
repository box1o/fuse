# Guest jobs and logs

Every newly built image contains:

```text
/usr/local/lib/rune/shell/lib/       Reusable Rune shell library
/usr/local/bin/rune-guest            Job runner
/usr/local/bin/rune-guest-queue      Boot-time queue worker
/var/lib/rune/queue                  Jobs waiting for the next boot
/var/lib/rune/completed              Completed job scripts and statuses
/var/log/rune/jobs/<id>/output.log   Durable combined stdout and stderr
```

## Host workflow

Queue jobs while the VM is stopped:

```bash
rune vm enqueue worker --id configure -- bash -lc './configure'
rune vm enqueue worker --id build --workdir /workspace/project -- \
    bash -lc 'cmake -S . -B build && cmake --build build'
rune vm start worker
rune vm wait-job worker build --timeout 1800
```

The queue runs jobs sequentially during boot. Each job command is stored using
shell-safe argument quoting. Output is written inside the guest and mirrored to
the Firecracker serial console.

Read logs from the serial console while running:

```bash
rune vm job-logs worker build --source console --lines 200
rune --json vm job-status worker build
```

After stopping the VM, read the durable log directly from ext4:

```bash
rune vm stop worker
rune vm job-logs worker build --source disk
```

`--source auto` selects console for a running VM and disk for a stopped VM.

Applications that need the same capture behavior should be launched through
`rune-guest run`. Arbitrary daemon logs are not automatically collected. A
future online agent may stream journal and application logs over vsock without
changing this durable on-disk format.
