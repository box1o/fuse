# Fuse CLI

The Fuse CLI supports both direct commands and an interactive terminal menu.
Both entry points call the same authentication service.

```bash
fuse
fuse auth login
fuse auth status
fuse auth logout
```

The hosted API defaults to `https://mback.teckstate.com`. Override it with
`FUSE_API_URL` or the global `--api-url` flag.

## Authentication

Login uses the OAuth device authorization pattern:

1. The CLI requests a one-time device code.
2. The user approves that code in the browser.
3. The CLI polls until authorization completes.
4. The returned credential is stored in the system keyring when available.
5. Systems without a compatible keyring use a `0600` file under the user
   configuration directory.

```bash
fuse auth login
fuse auth login --no-browser
fuse auth status --json
fuse auth logout --yes
```

Build the CLI from the repository root:

```bash
make cli
./bin/cli/fuse --version
```

Release builds can inject a version with `make cli CLI_VERSION=v1.2.3`.

## Automation

Use `--json` for machine-readable output, `--quiet` to suppress normal output,
and `--no-interactive` to guarantee that a command never prompts. JSON and
quiet modes cannot be combined.

When invoked without a command, Fuse opens the Huh menu only when stdin and
stdout are terminals. Otherwise, it prints command help instead of blocking.

## Package boundaries

- `internal/auth`: authentication workflows and ports
- `internal/api`: hosted HTTP API adapter
- `internal/credentials`: keyring and protected-file credential storage
- `internal/commands`: direct command interface
- `internal/interactive`: Huh menu interface
- `internal/ui`: human, JSON, quiet, prompt, and progress presentation
- `internal/platform`: browser, terminal, and timing adapters
- `internal/config`: CLI-only endpoint configuration

The CLI intentionally does not load the server configuration because that
configuration requires backend database, Redis, Stripe, mail, and OAuth
secrets.
