#!/usr/bin/env bash

set -Eeuo pipefail

__ENTRYPOINT_FILE__="$(readlink -f -- "${BASH_SOURCE[0]}")"
readonly __ENTRYPOINT_FILE__
__ENTRYPOINT_DIR__="$(cd -- "$(dirname -- "${__ENTRYPOINT_FILE__}")" && pwd -P)"
readonly __ENTRYPOINT_DIR__
__PROJECT_ROOT__="$(cd -- "${__ENTRYPOINT_DIR__}/../.." && pwd -P)"
readonly __PROJECT_ROOT__

# shellcheck source=/dev/null
source "${__PROJECT_ROOT__}/shell/runtime/init.sh"

runtime::main "$@"
