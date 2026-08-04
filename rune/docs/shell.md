# Bash Coding Convention

## 1. Purpose

This convention defines how Bash scripts should be written across the project.

It applies to:

* reusable library files;
* runtime code;
* module implementations;
* command handlers;
* tests;
* CLI entrypoints;
* helper scripts.

The goals are:

* predictable code structure;
* minimal variable collisions;
* safe filesystem and process handling;
* consistent error behavior;
* readable and maintainable scripts;
* easy testing;
* compatibility with ShellCheck and `shfmt`.

---

# 2. Bash version

All scripts target Bash, not POSIX `sh`.

Use this shebang:

```bash
#!/usr/bin/env bash
```

Do not use:

```bash
#!/bin/sh
```

unless the script is intentionally POSIX-compatible.

The codebase may use Bash-specific features such as:

```bash
[[ ... ]]
declare -A
local -n
arrays
process substitution
BASH_SOURCE
```

Minimum supported version:

```text
Bash 4.3+
```

When a feature requires a newer version, document it explicitly.

---

# 3. Shell options

Executable scripts should normally enable strict mode:

```bash
#!/usr/bin/env bash

set -Eeuo pipefail
```

Meaning:

```text
-E           Preserve ERR traps in functions and subshells.
-e           Exit executable scripts on unhandled command failure.
-u           Treat unset variables as errors.
-o pipefail  A pipeline fails when any command in it fails.
```

Do not enable strict mode inside reusable sourced library files.

Library files must not unexpectedly modify the caller's shell state.

Correct:

```bash
# shell/runtime/cli.sh (invoked through the rune symlink)
set -Eeuo pipefail

source "${PROJECT_ROOT}/shell/lib/init.sh"
```

Avoid:

```bash
# shell/lib/fs.sh
set -Eeuo pipefail
```

A sourced library changing shell options can break unrelated callers.

---

# 4. File naming

Use lowercase names with hyphens or simple words.

Preferred:

```text
init.sh
fs.sh
net.sh
command-registry.sh
module-loader.sh
```

For namespace files, a short namespace name is preferred:

```text
fs.sh
proc.sh
path.sh
yaml.sh
```

Avoid:

```text
FileSystemUtils.sh
my_script.SH
helperFunctions.sh
```

Module command implementation files should match command names:

```text
download.sh
verify.sh
install.sh
list.sh
```

---

# 5. Function naming

Use namespace-style function names:

```bash
namespace::function()
{
    # ...
}
```

Examples:

```bash
fs::mkdir()
net::download()
paths::register()
module::load()
command::dispatch()
```

Public functions:

```bash
namespace::function
```

Private helpers:

```bash
namespace::_function
```

Examples:

```bash
net::_require_curl()
path::_normalize()
command::_make_id()
```

Module commands:

```bash
module_name::command()
```

Examples:

```bash
firecracker::download()
network::verify()
artifact::package()
```

Command-specific helpers:

```bash
firecracker::download::_build_url()
firecracker::download::_usage()
```

Do not use generic global function names:

```bash
run()
init()
check()
download()
```

These names collide easily.

---

# 6. Function declaration style

Use this form:

```bash
fs::mkdir()
{
    local __directory__=${1:-}

    # ...
}
```

Do not use:

```bash
function fs::mkdir {
}
```

The selected style is easier to scan and consistent with the existing library.

Use one blank line between the function name and logical sections.

Example:

```bash
fs::copy_file()
{
    local __source__=${1:-}
    local __destination__=${2:-}

    validate::not_empty "source" "${__source__}" || return
    validate::not_empty "destination" "${__destination__}" || return

    fs::mkdir_parent "${__destination__}" || return

    cp --preserve=mode,timestamps -- \
        "${__source__}" \
        "${__destination__}"
}
```

---

# 7. Local variable naming

Use double underscores before and after local variable names:

```bash
local __name__
local __source__
local __destination__
local __status__
```

For multiple words, use lowercase snake case:

```bash
local __module_name__
local __project_root__
local __destination_directory__
local __expected_checksum__
```

Arrays:

```bash
local -a __arguments__=()
local -a __headers__=()
```

Associative arrays:

```bash
local -A __metadata__=()
local -A __options__=()
```

The double-underscore convention is important because Bash uses dynamic scoping.

A function can see local variables from its caller. Distinctive names reduce accidental collisions.

Preferred:

```bash
local __path__
```

Avoid:

```bash
local path
local _path
local Path
```

---

# 8. Nameref variables

Use `local -n` when a function must modify a caller-owned array or variable.

Nameref variables use one leading and trailing underscore:

```bash
local -n _registry_=${1}
```

Example:

```bash
array::append()
{
    local -n _array_=${1}
    local __value__=${2:-}

    _array_+=("${__value__}")
}
```

Convention:

```text
__name__   Normal local variable
_name_     Nameref variable
NAME       Global constant or registry
```

Always validate the nameref variable name when it comes from external input.

Do not use nameref unless it materially improves the API.

---

# 9. Global variables

Global variables should be rare.

Use uppercase snake case:

```bash
RUNTIME_VERSION="0.1.0"
LOG_LEVEL=20
DRY_RUN=false
```

Readonly constants:

```bash
readonly NET_DEFAULT_RETRIES=5
readonly NET_DEFAULT_RETRY_DELAY=2
readonly RUNTIME_DEFAULT_MODULE_DIR="runtime/modules"
```

Global associative registries:

```bash
declare -gA PATH_REGISTRY=()
declare -gA MODULE_REGISTRY=()
declare -gA COMMAND_HANDLER_REGISTRY=()
```

Only the owning namespace should mutate a global registry directly.

For example:

```text
paths::*     Owns PATH_REGISTRY
module::*    Owns MODULE_REGISTRY
command::*   Owns COMMAND_* registries
```

Module code must use public APIs instead of modifying registries directly.

Avoid broad globals:

```bash
SOURCE="/tmp/file"
NAME="test"
STATUS=1
```

These are too generic and collision-prone.

---

# 10. Constants

Use `readonly`:

```bash
readonly ARCHIVE_DEFAULT_FORMAT="tar.xz"
readonly RETRY_DEFAULT_ATTEMPTS=5
```

For local constants:

```bash
local -r __default_timeout__=30
```

Do not declare a value readonly if it is expected to be overridden by environment configuration.

Use:

```bash
NET_DEFAULT_RETRIES=${NET_DEFAULT_RETRIES:-5}
```

instead of:

```bash
readonly NET_DEFAULT_RETRIES=5
```

when callers are allowed to configure it before sourcing the library.

A good distinction is:

```text
*_DEFAULT_*   Configurable default
*_VERSION     Immutable library metadata
```

---

# 11. Boolean values

Represent booleans as strings:

```bash
local __force__=false
local __quiet__=false
local __dry_run__=true
```

Check them explicitly:

```bash
if [[ "${__force__}" == true ]]; then
    # ...
fi
```

Do not rely on:

```bash
if ${__force__}; then
```

because it executes the string as a command.

Accepted values should normally be:

```text
true
false
```

Validation:

```bash
validate::boolean "force" "${__force__}" || return
```

---

# 12. Function arguments

Read positional arguments into named local variables immediately.

```bash
fs::move()
{
    local __source__=${1:-}
    local __destination__=${2:-}

    # ...
}
```

Do not use `$1`, `$2`, and `$3` repeatedly throughout a long function.

Preferred:

```bash
local __module_name__=${1:-}
local __command_name__=${2:-}
```

Avoid:

```bash
if [[ -z "$1" ]]; then
    ...
fi

echo "$2"
```

This improves readability and protects against unset positional parameters under `set -u`.

---

# 13. Optional arguments

Use safe default expansion:

```bash
local __name__=${1:-}
local __timeout__=${2:-30}
```

For environment defaults:

```bash
local __cache_dir__=${CACHE_DIR:-"${HOME}/.cache/app"}
```

Do not use:

```bash
local __name__=$1
```

under `set -u`, because missing arguments produce an immediate failure.

---

# 14. Option parsing

For simple command handlers, use a `while` loop and `case`.

```bash
example::run()
{
    local __output__=
    local __force__=false
    local __retries__=5

    while [[ $# -gt 0 ]]; do
        case "${1}" in
            --output)
                [[ $# -ge 2 ]] || {
                    log::error "example run: --output requires a value"
                    return 2
                }

                __output__=${2}
                shift 2
                ;;

            --retries)
                [[ $# -ge 2 ]] || {
                    log::error "example run: --retries requires a value"
                    return 2
                }

                __retries__=${2}
                shift 2
                ;;

            --force)
                __force__=true
                shift
                ;;

            --help|-h)
                example::run::_usage
                return 0
                ;;

            --)
                shift
                break
                ;;

            -*)
                log::error "example run: unknown option: ${1}"
                return 2
                ;;

            *)
                break
                ;;
        esac
    done

    # Remaining positional arguments are still in "$@".
}
```

Rules:

* Validate values immediately.
* Return `2` for invalid usage.
* Support `--`.
* Keep command-specific option parsing in the command handler.
* Do not use `eval`.
* Avoid `getopts` when long options are required.

For handlers with more than one or two options, prefer the shared `args` API:

```bash
args::reset
args::flag --name verbose --short v
args::option --name output --short o --required
args::option --name format --default text

args::parse "$@" || return

if args::has verbose; then
    log::debug "Verbose output enabled"
fi

local __output__
__output__=$(args::get output) || return
```

The parser supports long options, `--name=value`, single short options,
defaults, required values, ordered positional arguments, and `--`. Call
`args::reset` before defining a command's schema. Combined short flags are
intentionally unsupported so parsing remains explicit.

---

# 15. Arrays

Use arrays for command arguments.

Correct:

```bash
local -a __arguments__=(
    --fail
    --location
    --retry "${__retries__}"
)

if [[ "${__quiet__}" == true ]]; then
    __arguments__+=(--silent)
fi

curl "${__arguments__[@]}" -- "${__url__}"
```

Avoid building command strings:

```bash
local __args__="--fail --retry ${__retries__}"
curl ${__args__} "${__url__}"
```

Command strings introduce word splitting and quoting bugs.

Append to arrays using:

```bash
__arguments__+=("--header" "${__header__}")
```

Use:

```bash
"${__arguments__[@]}"
```

not:

```bash
"${__arguments__[*]}"
```

when passing arguments to a command.

---

# 16. Associative arrays

Declare:

```bash
declare -gA PATH_REGISTRY=()
```

Assign:

```bash
PATH_REGISTRY["runtime|modules"]="/project/runtime/modules"
```

Read:

```bash
local __modules_dir__=${PATH_REGISTRY["runtime|modules"]}
```

Check existence:

```bash
[[ -v 'PATH_REGISTRY[runtime|modules]' ]]
```

When keys are stored in variables:

```bash
local __key__="runtime|modules"

[[ -v "PATH_REGISTRY[${__key__}]" ]]
```

Avoid ambiguous direct mutation outside the owning namespace.

---

# 17. Quoting

Quote variable expansions by default:

```bash
"${__path__}"
"${HOME}"
"${__arguments__[@]}"
```

Correct:

```bash
rm -- "${__file__}"
cp -- "${__source__}" "${__destination__}"
```

Avoid:

```bash
rm $__file__
cp $__source__ $__destination__
```

Unquoted values break on spaces, wildcard characters, and empty values.

Intentional unquoted expansion should be rare and documented.

---

# 18. Use `--` before path arguments

When invoking external commands with paths, use `--` when supported:

```bash
rm -- "${__file__}"
mkdir -p -- "${__directory__}"
cp -a -- "${__source__}" "${__destination__}"
mv -- "${__source__}" "${__destination__}"
```

This prevents filenames beginning with `-` from being interpreted as options.

---

# 19. Tests and conditions

Use `[[ ... ]]` instead of `[ ... ]`.

Preferred:

```bash
if [[ -f "${__file__}" ]]; then
    # ...
fi
```

Avoid:

```bash
if [ -f "$file" ]; then
fi
```

Common tests:

```bash
[[ -e "${__path__}" ]]
[[ -f "${__file__}" ]]
[[ -d "${__directory__}" ]]
[[ -L "${__path__}" ]]
[[ -r "${__file__}" ]]
[[ -w "${__path__}" ]]
[[ -x "${__file__}" ]]
[[ -n "${__value__}" ]]
[[ -z "${__value__}" ]]
```

For broken symlinks:

```bash
[[ -e "${__path__}" || -L "${__path__}" ]]
```

String comparison:

```bash
[[ "${__mode__}" == "production" ]]
```

Pattern matching:

```bash
[[ "${__name__}" == module-* ]]
```

Regex:

```bash
[[ "${__name__}" =~ ^[a-z][a-z0-9-]*$ ]]
```

---

# 20. Arithmetic

Use arithmetic contexts:

```bash
if (( __attempt__ >= __maximum_attempts__ )); then
    # ...
fi
```

Increment:

```bash
((__attempt__ += 1))
```

Use decimal validation before arithmetic when values come from users.

```bash
validate::positive_integer "attempts" "${__attempts__}" || return
```

Avoid:

```bash
[ "$count" -gt 5 ]
```

unless compatibility requires it.

---

# 21. Command substitution

Use:

```bash
local __output__
__output__=$(command)
```

Do not use backticks:

```bash
output=`command`
```

Preserve command status:

```bash
local __output__
local __status__

__output__=$(command)
__status__=$?

if (( __status__ != 0 )); then
    return "${__status__}"
fi
```

Or:

```bash
if ! __output__=$(command); then
    log::error "Command failed"
    return 1
fi
```

Do not use the negated form when you must preserve the exact original status without additional handling.

---

# 22. Pipelines

Because executables use `set -o pipefail`, failures inside pipelines are visible.

Preferred:

```bash
printf '%s\n' "${__items__[@]}" |
    sort |
    uniq
```

For complex pipeline status handling, avoid hiding the result:

```bash
local __status__

some_command |
    another_command

__status__=$?
```

Use temporary files or process substitution when pipeline behavior becomes unclear.

---

# 23. Return versus exit

Library functions and module handlers must use `return`.

Correct:

```bash
fs::copy_file()
{
    cp -- "${__source__}" "${__destination__}" || return
}
```

Incorrect:

```bash
fs::copy_file()
{
    cp -- "${__source__}" "${__destination__}" || exit 1
}
```

Only executable entrypoints may use `exit`.

Example:

```bash
runtime::main "$@"
exit $?
```

`log::fatal` may call `exit`, but it should be reserved for top-level unrecoverable conditions.

---

# 24. Exit codes

Use consistent exit codes:

```text
0     Success
1     Operational failure
2     Invalid arguments or usage
126   Command exists but is not executable
127   Required command not found
```

Examples:

```bash
log::error "Missing --output value"
return 2
```

```bash
log::error "Download failed"
return 1
```

Preserve external command statuses when they are meaningful.

---

# 25. Error handling

Use this pattern:

```bash
operation || return
```

When context is needed:

```bash
if ! fs::copy_file "${__source__}" "${__destination__}"; then
    log::error "Failed to copy configuration file"
    return 1
fi
```

Do not duplicate errors unnecessarily.

Bad:

```bash
fs::copy_file "${__source__}" "${__destination__}" || {
    log::error "Copy failed"
    return 1
}
```

when `fs::copy_file` already emitted a complete actionable error.

Add context only when the caller knows something useful that the callee does not.

---

# 26. Logging

Use log levels consistently.

```bash
log::debug "Resolved cache path: ${__cache__}"
log::info "Downloading release ${__version__}"
log::warn "Checksum was not provided"
log::error "Unable to open configuration file"
```

Guidance:

```text
DEBUG   Internal implementation detail
INFO    Normal user-visible progress
WARN    Recoverable or optional problem
ERROR   Operation failed
FATAL   Application cannot continue
```

Do not log passwords, tokens, authorization headers, private keys, or sensitive environment values.

Use concise messages with relevant context.

Preferred:

```bash
log::error "Configuration file does not exist: ${__file__}"
```

Avoid:

```bash
log::error "Something went wrong"
```

---

# 27. Path handling

Use `path::` functions for path manipulation and `fs::` functions for filesystem mutation.

Correct:

```bash
local __output__
__output__=$(path::join "${__root__}" "output/archive.tar.xz") || return

fs::mkdir_parent "${__output__}" || return
```

Do not manually concatenate paths everywhere:

```bash
local __output__="${__root__}/output/archive.tar.xz"
```

Manual concatenation is acceptable for simple controlled paths, but shared path behavior should use `path::join`.

Before destructive operations:

```bash
path::assert_within "${__target__}" "${__allowed_root__}" || return
fs::remove_tree "${__target__}"
```

All registered paths should be absolute.

---

# 28. Filesystem safety

Never recursively remove an unchecked variable.

Unsafe:

```bash
rm -rf "${__directory__}"
```

Required pattern:

```bash
validate::absolute_path "${__directory__}" || return
path::assert_within "${__directory__}" "${__allowed_root__}" || return
fs::remove_tree "${__directory__}"
```

Refuse dangerous values:

```text
empty
/
.
..
```

Use atomic writes for configuration and generated files:

```bash
fs::atomic_write "${__config_file__}" <<'EOF'
enabled=true
EOF
```

---

# 29. External process execution

Use `proc::` wrappers when available.

```bash
proc::require curl
proc::run curl --version
```

Capture output:

```bash
local __architecture__
__architecture__=$(proc::capture uname -m) || return
```

Allow failure deliberately:

```bash
if proc::run_allow_fail test -e "${__optional_file__}"; then
    # ...
fi
```

Do not execute generated strings.

Unsafe:

```bash
eval "${__command__}"
```

Safe:

```bash
local -a __command__=(
    curl
    --fail
    "${__url__}"
)

"${__command__[@]}"
```

---

# 30. Temporary files

Use `mktemp` through the `temp::` API.

```bash
local __temporary_file__
__temporary_file__=$(temp::file "download") || return
```

```bash
local __temporary_directory__
__temporary_directory__=$(temp::dir "extract") || return
```

Register cleanup:

```bash
temp::register "${__temporary_directory__}"
```

Do not use predictable names:

```bash
/tmp/my-temp-file
```

---

# 31. Traps and cleanup

Use the trap helper so multiple cleanup handlers can coexist.

```bash
trap::add EXIT 'temp::cleanup'
trap::add EXIT 'mount::cleanup'
```

Cleanup functions must preserve the original exit status when used directly.

Example:

```bash
app::cleanup()
{
    local __status__=$?

    temp::cleanup || true

    return "${__status__}"
}
```

Avoid setting traps inside low-level functions unless the function owns the full lifecycle.

---

# 32. Environment variables

Environment variables use uppercase snake case:

```bash
APP_CACHE_DIR
NET_DEFAULT_RETRIES
LOG_LEVEL
LOG_COLOR
```

Read with defaults:

```bash
local __cache__=${APP_CACHE_DIR:-"${HOME}/.cache/app"}
```

Require values through the environment API:

```bash
env::require HOME
```

Do not silently export variables unless exporting is the function's stated purpose.

Avoid sourcing untrusted `.env` files. A sourced environment file is executable shell code.

---

# 33. YAML usage

Use `yaml::` APIs and `yq`.

Correct:

```bash
yaml::validate "${__config__}" || return

local __port__
__port__=$(yaml::get "${__config__}" '.server.port' '8080') || return
```

Typed values:

```bash
yaml::set_raw "${__config__}" '.server.enabled' 'true'
yaml::set_raw "${__config__}" '.server.port' '8080'
```

Strings:

```bash
yaml::set "${__config__}" '.server.host' 'localhost'
```

Do not parse YAML with:

```bash
grep
sed
awk
cut
```

YAML syntax is too complex for reliable line-based parsing.

---

# 34. Source paths

Resolve a file's own directory using `BASH_SOURCE`.

```bash
readonly __SCRIPT_DIR__="$(
    cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &&
        pwd -P
)"
```

Do not rely on the current working directory.

Unsafe:

```bash
source ./lib/init.sh
```

Preferred:

```bash
source "${__SCRIPT_DIR__}/../lib/init.sh"
```

Use uppercase readonly names for file-scope resolved directories:

```bash
readonly __MODULE_DIR__=...
```

For module-specific constants:

```bash
readonly __FIRECRACKER_MODULE_DIR__=...
```

---

# 35. Source guards

Reusable initialization files should be idempotent.

```bash
if [[ ${SHELL_LIB_INITIALIZED:-0} == 1 ]]; then
    return 0
fi

SHELL_LIB_INITIALIZED=1
```

Individual files may also use guards when they can be sourced independently.

```bash
if [[ ${SHELL_LIB_FS_LOADED:-0} == 1 ]]; then
    return 0
fi

SHELL_LIB_FS_LOADED=1
```

Do not use readonly guards if later sourcing may attempt to assign the same value.

---

# 36. Dependency checks

Check tools near the operation that requires them.

```bash
yaml::_require_yq()
{
    proc::require yq
}
```

Do not require every optional tool during global library initialization.

For example, a missing `yq` should not prevent a user from running a filesystem-only command.

Built-in diagnostics may report optional dependencies:

```text
curl
yq
tar
flock
bats
shellcheck
shfmt
```

---

# 37. Comments

Comments should explain why, not restate what the code already says.

Useful:

```bash
# Keep the partial file so the next invocation can resume the download.
```

Unnecessary:

```bash
# Set retries to 5.
local __retries__=5
```

Use comments for:

* non-obvious safety constraints;
* external tool limitations;
* compatibility workarounds;
* lifecycle decisions;
* unusual error handling.

Do not leave commented-out dead code in production files.

---

# 38. Function size

Functions should be focused and normally remain below approximately 40–60 lines.

Split when a function handles multiple responsibilities.

Instead of one large function:

```bash
net::download()
```

extract:

```bash
net::download::_parse_options
net::download::_build_arguments
net::download::_verify
net::download::_finalize
```

Do not split trivial logic purely to reduce line count.

---

# 39. Public API documentation

Every public function should have a short documentation comment when its behavior is not obvious.

Example:

```bash
# Register a path under a unique key.
#
# Usage:
#   paths::register [options] <key> <path>
#
# Returns:
#   0 on success
#   1 when the key already exists
#   2 for invalid arguments
paths::register()
{
    # ...
}
```

Private helpers do not require full documentation unless their behavior is subtle.

---

# 40. Command handler structure

Use this order:

```text
1. Local defaults
2. Option parsing
3. Argument validation
4. Dependency checks
5. Path resolution
6. Operation
7. Final log/output
```

Example:

```bash
artifact::download()
{
    local __version__="latest"
    local __output__=
    local __force__=false

    while [[ $# -gt 0 ]]; do
        # Parse options
    done

    validate::not_empty "version" "${__version__}" || return
    proc::require curl || return

    if [[ -z "${__output__}" ]]; then
        __output__=$(artifact::download::_default_output "${__version__}") ||
            return
    fi

    net::download \
        "${__url__}" \
        "${__output__}" ||
        return

    log::info "Downloaded artifact: ${__output__}"
}
```

---

# 41. Module structure

Recommended module structure:

```text
runtime/modules/example/
├── module.sh
├── vars.sh
├── run.sh
├── verify/
│   └── command.sh
├── README.md
└── tests/
    ├── run.bats
    └── verify.bats
```

Responsibilities:

```text
module.sh   Module and public command declarations only
vars.sh     Optional module variable initialization
run.sh      Small run command implementation
verify/     Files related to the verify command
tests/      Module-specific tests
```

Rune recursively sources module `*.sh` files in lexical order, excluding
`module.sh` and `tests/`. It sources `module.sh` last. Module authors must not
resolve the module directory or manually source sibling files.

Minimal declaration:

```bash
module::register \
    --name example \
    --description "Example commands"

command::register \
    --cmd run \
    --handler example::run \
    --description "Run the example"
```

Only `--name` is required by `module::register`. Only `--cmd` and `--handler`
are required by `command::register`. The command registry infers the module
currently being loaded. The longer metadata flags remain optional.

Do not put implementations in `module.sh`.

Prefer one file for a small command. A directory such as `verify/` is useful
only when a command has enough options, helpers, or distinct operations to make
the separation easier to read.

---

# 42. Module variable conventions

Module constants:

```bash
readonly EXAMPLE_DEFAULT_VERSION="latest"
readonly EXAMPLE_API_BASE_URL="https://example.com"
```

Module environment overrides:

```bash
local __cache__=${EXAMPLE_CACHE_DIR:-"${HOME}/.cache/app/example"}
```

Initialize module-owned values through the context passed by Rune:

```bash
example::_init()
{
    local -n _example_=${1}

    local __cache__=${EXAMPLE_CACHE_DIR:-"${HOME}/.cache/rune/example"}
    local __build__="${__cache__}/build"

    var::update _example_ cache "${__cache__}"
    var::update _example_ build "${__build__}"
    var::update _example_ binary "${__build__}/example"
}
```

Rune detects `<module>::_init` automatically and invokes it once before the
selected command. Read values with `module::var example cache`. The generic
path registry remains available for values that truly need global inspection,
but ordinary module variables do not need registry boilerplate.

Command help belongs to the runtime. `rune example run --help` prints the
registered usage and description without calling the handler. A command only
needs custom help logic when it deliberately implements a richer help system.

Do not create generic globals such as:

```bash
CACHE_DIR
DOWNLOAD_DIR
VERSION
```

Prefix module-level globals:

```bash
EXAMPLE_CACHE_DIR
EXAMPLE_DEFAULT_VERSION
```

---

# 43. Output conventions

Machine-readable functions should print only the requested value to stdout.

Example:

```bash
paths::get "runtime|modules"
```

Output:

```text
/project/runtime/modules
```

Diagnostics must go to stderr through logging functions.

This allows:

```bash
local __modules__
__modules__=$(paths::require "runtime|modules") || return
```

Do not mix logs and data on stdout.

User-facing table commands may intentionally print formatted output to stdout:

```bash
paths::print
module::list
command::list
```

---

# 44. Heredocs

Use quoted heredoc delimiters when variable expansion is not wanted:

```bash
cat <<'EOF'
${HOME} is printed literally.
EOF
```

Use unquoted delimiters only when interpolation is intentional:

```bash
cat <<EOF
Project root: ${PROJECT_ROOT}
EOF
```

Indent carefully. Use `<<-EOF` only when tab indentation is intentionally required.

---

# 45. Reading input

Read safely:

```bash
while IFS= read -r __line__; do
    # ...
done < "${__file__}"
```

For null-delimited filenames:

```bash
while IFS= read -r -d '' __file__; do
    # ...
done < <(
    find "${__directory__}" -type f -print0
)
```

Do not use:

```bash
for file in $(find ...)
```

It breaks on whitespace and special characters.

---

# 46. Sorting and deterministic behavior

Module discovery, command listing, and generated output should be deterministic.

Use sorting:

```bash
printf '%s\n' "${!MODULE_REGISTRY[@]}" |
    sort
```

Never rely on associative-array iteration order.

Tests should not depend on filesystem directory order.

---

# 47. Concurrency and locks

Use `lock::` for operations that cannot safely run concurrently.

```bash
local __lock_file__
__lock_file__=$(paths::require "runtime|lock") || return

lock::acquire "${__lock_file__}" || {
    log::error "Another operation is already running"
    return 1
}
```

Release through cleanup handling.

Do not invent lock files in random locations. Register them in the path registry.

---

# 48. Retry behavior

Use retries only for transient failures.

Appropriate:

```text
network timeouts
connection refusal
temporary service unavailability
device detection
```

Not appropriate:

```text
invalid arguments
missing files
syntax errors
checksum mismatch
permission errors
```

Example:

```bash
retry::run \
    --attempts 5 \
    --delay 2 \
    curl --fail "${__url__}"
```

Defaults should handle common cases without requiring every option.

---

# 49. Testing conventions

Use Bats.

Test names describe behavior:

```bash
@test "fs::copy_file preserves file content" {
}
```

Avoid:

```bash
@test "copy test" {
}
```

Use temporary test directories:

```bash
local __source__="${BATS_TEST_TMPDIR}/source file"
```

Test paths with spaces and leading hyphens.

Each public API should have tests for:

```text
success
missing required argument
invalid argument
operational failure
important edge cases
```

Commands should test exit statuses and output separately.

```bash
run example::run --bad-option

[ "${status}" -eq 2 ]
[[ "${output}" == *"unknown option"* ]]
```

---

# 50. Static analysis and formatting

Run ShellCheck:

```bash
shellcheck shell/lib/*.sh
shellcheck shell/runtime/*.sh
shellcheck runtime/modules/**/*.sh
```

Format with:

```bash
shfmt -w -i 4 -ci shell runtime/modules
```

Validate syntax:

```bash
find shell runtime/modules \
    -type f \
    -name '*.sh' \
    -print0 |
    xargs -0 -n1 bash -n
```

Do not suppress ShellCheck warnings without a specific reason.

When suppression is needed:

```bash
# shellcheck disable=SC1090
source "${__dynamic_file__}"
```

Keep suppressions as narrow as possible.

---

# 51. Preferred complete function example

```bash
# Download a file with automatic resume and optional SHA-256 validation.
#
# Usage:
#   example::download [--sha256 HASH] [--force] <url> <destination>
example::download()
{
    local __sha256__=
    local __force__=false

    while [[ $# -gt 0 ]]; do
        case "${1}" in
            --sha256)
                [[ $# -ge 2 ]] || {
                    log::error "example download: --sha256 requires a value"
                    return 2
                }

                __sha256__=${2}
                shift 2
                ;;

            --force)
                __force__=true
                shift
                ;;

            --help|-h)
                example::download::_usage
                return 0
                ;;

            --)
                shift
                break
                ;;

            -*)
                log::error "example download: unknown option: ${1}"
                return 2
                ;;

            *)
                break
                ;;
        esac
    done

    if [[ $# -ne 2 ]]; then
        log::error "Usage: example download [options] <url> <destination>"
        return 2
    fi

    local __url__=${1}
    local __destination__=${2}

    validate::not_empty "URL" "${__url__}" || return
    validate::not_empty "destination" "${__destination__}" || return

    local -a __arguments__=()

    if [[ "${__force__}" == true ]]; then
        __arguments__+=(--overwrite)
    fi

    if [[ -n "${__sha256__}" ]]; then
        __arguments__+=(--sha256 "${__sha256__}")
    fi

    net::download \
        "${__arguments__[@]}" \
        "${__url__}" \
        "${__destination__}"
}
```

This example demonstrates:

* namespace naming;
* `__variable__` locals;
* safe default values;
* option parsing;
* usage errors returning `2`;
* arrays for forwarding options;
* quoted arguments;
* delegation to a reusable library API;
* no `eval`;
* no direct `exit`.

---

# 52. Code review checklist

Before merging Bash code, verify:

```text
[ ] Uses Bash shebang.
[ ] Executable enables strict mode.
[ ] Sourced library does not change caller shell options.
[ ] Functions are namespaced.
[ ] Private helpers use `_` after namespace.
[ ] Local variables use `__name__`.
[ ] Namerefs use `_name_`.
[ ] Globals are uppercase and appropriately scoped.
[ ] Variable expansions are quoted.
[ ] External commands receive `--` before path arguments when supported.
[ ] Commands are passed as arrays, not strings.
[ ] No use of eval.
[ ] Low-level functions return instead of exit.
[ ] Invalid usage returns 2.
[ ] Destructive paths are validated.
[ ] Secrets are not logged.
[ ] stdout contains data; stderr contains diagnostics.
[ ] Module loading has no expensive side effects.
[ ] Initialization is idempotent.
[ ] Public APIs have tests.
[ ] Bash syntax validation passes.
[ ] ShellCheck passes.
[ ] shfmt formatting passes.
```

---

# 53. Final convention summary

Use these patterns consistently:

```bash
namespace::public_function()
namespace::_private_function()

local __normal_variable__
local -a __array__=()
local -A __map__=()
local -n _nameref_=${1}

readonly GLOBAL_CONSTANT="value"
declare -gA GLOBAL_REGISTRY=()
```

Use:

```bash
[[ ... ]]
(( ... ))
"${array[@]}"
command -- "${path}"
operation || return
```

Avoid:

```bash
[ ... ]
eval
unquoted variables
command strings
generic global names
exit inside library functions
raw rm -rf on unchecked paths
```

The core philosophy is:

```text
Explicit inputs
Explicit registration
Explicit errors
Safe defaults
Small namespaced functions
No hidden global behavior
Deterministic output
Composable return-based APIs
```
