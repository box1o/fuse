#!/usr/bin/env bash

if [[ ${RUNTIME_INITIALIZED:-0} == 1 ]]; then return 0; fi
RUNTIME_INITIALIZED=1

RUNTIME_SHELL_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
readonly RUNTIME_SHELL_DIR
RUNTIME_PROJECT_ROOT="$(cd -- "${RUNTIME_SHELL_DIR}/../.." && pwd -P)"
readonly RUNTIME_PROJECT_ROOT

# shellcheck source=/dev/null
source "${RUNTIME_PROJECT_ROOT}/shell/lib/init.sh"
source "${RUNTIME_SHELL_DIR}/config.sh"
source "${RUNTIME_PROJECT_ROOT}/runtime/config/runtime.sh"
source "${RUNTIME_SHELL_DIR}/commands.sh"
source "${RUNTIME_SHELL_DIR}/modules.sh"
source "${RUNTIME_SHELL_DIR}/help.sh"
source "${RUNTIME_SHELL_DIR}/scaffold.sh"
source "${RUNTIME_SHELL_DIR}/main.sh"

runtime::config::set 'cli|name' "${RUNE_CLI_NAME:-rune}"
runtime::config::set 'modules|path_env' RUNE_MODULE_PATH
runtime::config::set 'output|format' text
runtime::register_paths "${RUNTIME_PROJECT_ROOT}"
runtime::register_builtins
