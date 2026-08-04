# Volume module

The volume module owns reusable file-backed ext4 volumes. Growth is supported;
shrinking is intentionally excluded.

```bash
rune volume create worker-root --from /path/to/rootfs.ext4
rune volume show worker-root
rune volume check worker-root
rune volume grow worker-root --size 20G
rune volume cleanup
```

`cleanup` only lists unattached volumes unless `--force` is provided.
