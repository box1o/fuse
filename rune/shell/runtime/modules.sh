#!/usr/bin/env bash

readonly MODULE_RUNTIME_DIR="${RUNTIME_SHELL_DIR}/module"

# shellcheck source=/dev/null
source "${MODULE_RUNTIME_DIR}/registration.sh"
# shellcheck source=/dev/null
source "${MODULE_RUNTIME_DIR}/access.sh"
# shellcheck source=/dev/null
source "${MODULE_RUNTIME_DIR}/lifecycle.sh"
# shellcheck source=/dev/null
source "${MODULE_RUNTIME_DIR}/loader.sh"
# shellcheck source=/dev/null
source "${MODULE_RUNTIME_DIR}/state.sh"
