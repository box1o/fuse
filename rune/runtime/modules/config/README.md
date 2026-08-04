# Configuration module

The configuration module exposes the active versioned YAML document through a
small command API:

```bash
rune config path
rune config validate
rune config get '.images | keys'
rune config show
```

Use `RUNE_CONFIG_FILE` or `--file` to select another document. The reusable
shell API is provided by `shell/lib/config.sh`.
