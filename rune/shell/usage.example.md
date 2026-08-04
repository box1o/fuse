# Shell library usage

Source `shell/lib/init.sh` from a Bash script, then call its namespaced functions.

## Basic example

Create a script at the repository root:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

source "${SCRIPT_DIR}/shell/lib/init.sh"

log::info "Starting application"

work_dir="$(paths::get 'system|tmp')"
output_file="${work_dir}/cracker/example.txt"

fs::atomic_write "${output_file}" <<'EOF'
Hello from the shell library
EOF

log::info "Created: ${output_file}"
```

Run it with:

```sh
bash example.sh
```

## Common helpers

### Argument parsing

```bash
args::reset
args::flag --name force --short f
args::option --name output --short o --required
args::option --name retries --default 3

args::parse "$@" || return

force=false
if args::has force; then
    force=true
fi

output="$(args::get output)" || return
retries="$(args::get retries)" || return
```

Positional values remain available in their original order:

```bash
count="$(args::count)"
first="$(args::position 0)"
mapfile -t all_positionals < <(args::positionals)
```

### Logging

```bash
log::debug "Debug message"
log::info "Information"
log::warn "Warning"
log::error "Error"
```

### Filesystem

```bash
fs::mkdir "/tmp/example"
fs::copy_file "source.txt" "/tmp/example/target.txt"
fs::remove_tree "/tmp/example"
```

### Paths

```bash
full_path="$(path::join "/tmp" "example" "file.txt")"

paths::register 'app|output' "$PWD/output"
output_dir="$(paths::get 'app|output')"
```

### Environment variables

```bash
env::require API_TOKEN
env::set_default APP_PORT 8080
port="$(env::get APP_PORT)"
```

### Commands

```bash
proc::require curl tar
proc::run echo "Hello"
```

### Retries

```bash
retry::run --attempts 3 --delay 2 -- curl https://example.com
```

### Downloads

```bash
net::download \
    "https://example.com/archive.tar.gz" \
    "/tmp/archive.tar.gz"
```

### YAML

The YAML helpers require mikefarah `yq` v4.

```bash
port="$(yaml::get config.yaml '.server.port' '8080')"
yaml::set config.yaml '.server.host' 'localhost'
yaml::set_raw config.yaml '.server.port' '9090'
```

### Versioned configuration documents

Load a YAML document once and address it by name:

```bash
config::load application ./config.yaml 1

host="$(config::get application '.server.host' '127.0.0.1')"
config::require application '.server.port'
config::keys application '.workers'
config::values application '.server.allowed_hosts'
```

`config::load` validates YAML and can enforce the top-level `version`. Multiple
documents may be loaded under different names in the same process.

## Using the library from a nested script

Resolve the repository root before sourcing the library:

```bash
#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd -P)"

source "${PROJECT_ROOT}/shell/lib/init.sh"
```

The library requires Bash. Run scripts with `bash script.sh`, not `sh script.sh`.

## Rune command runtime

After `make link`, invoke the runtime without a relative path prefix:

```bash
rune --help
rune module list
rune runtime doctor
rune paths list
```

Inspect Firecracker module help:

```bash
rune help vm
rune vm download --help
```

Create and validate a module:

```bash
rune module scaffold example \
    --command run \
    --summary "Example module"

rune module validate example
rune example run --help
```

Rune loads implementation files automatically. A module declaration only names
the module and its public commands:

```bash
module::register \
    --name example \
    --description "Example module"

command::register \
    --cmd run \
    --handler example::run \
    --description "Run the example"
```

Store module values in `vars.sh` with a nameref:

```bash
example::_init()
{
    local -n _example_=${1}

    local __source__="${HOME}/src/example"
    local __build__="${HOME}/.cache/rune/example/build"

    var::update _example_ source "${__source__}"
    var::update _example_ build "${__build__}"
    var::update _example_ binary "${__build__}/example"
}
```

Read a value from any command file:

```bash
example::run()
{
    local __binary__
    __binary__=$(module::var example binary) || return

    proc::run "${__binary__}" "$@"
}
```

`example::_init` is detected automatically. `rune example run --help` is
generated from the registration description; no `_usage` function is needed.

Additional trusted module roots can be supplied with:

```bash
RUNE_MODULE_PATH=/path/one:/path/two rune module list
```
