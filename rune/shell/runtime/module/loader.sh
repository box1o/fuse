#!/usr/bin/env bash

module::discover()
{
    local __root__=${1:-}

    validate::not_empty module-root "${__root__}" || return
    [[ -e "${__root__}" ]] || return 0
    fs::require_dir module-root "${__root__}" || return

    if [[ ! -r "${__root__}" ]]; then
        log::error "Module root is not readable: ${__root__}"
        return 1
    fi

    local __root_absolute__
    __root_absolute__=$(path::absolute "${__root__}") || return

    local __directory__
    while IFS= read -r -d '' __directory__; do
        module::_validate_directory "${__root_absolute__}" "${__directory__}" || return
        module::load "${__directory__}" || return
    done < <(
        find "${__root_absolute__}" \
            -mindepth 1 \
            -maxdepth 1 \
            -type d \
            -print0 |
            sort -z
    )
}

module::_validate_directory()
{
    local __root__=${1:-}
    local __directory__=${2:-}

    if [[ -L "${__directory__}" ]]; then
        log::error "Symlinked module directories are not allowed: ${__directory__}"
        return 1
    fi

    local __resolved__
    __resolved__=$(path::absolute "${__directory__}") || return
    path::assert_within "${__resolved__}" "${__root__}" || return

    local __name__
    __name__=$(path::basename "${__directory__}") || return
    module::_valid_name "${__name__}" || return
    fs::require_file module-declaration "${__directory__}/module.sh"
}

module::load()
{
    local __directory__=${1:-}
    local __normalized__
    local __name__

    __normalized__=$(path::absolute "${__directory__}") || return
    __name__=$(path::basename "${__normalized__}") || return

    module::_valid_name "${__name__}" || return

    if module::has "${__name__}"; then
        log::error "Duplicate module: ${__name__}"
        return 1
    fi

    fs::require_file module-declaration "${__normalized__}/module.sh" || return

    local -A __modules_before__=()
    local -A __commands_before__=()
    module::_snapshot_registries __modules_before__ __commands_before__

    MODULE_LOADING_NAME=${__name__}
    MODULE_LOADING_ROOT=${__normalized__}

    local __status__=0
    module::_source_implementations "${__normalized__}" || __status__=$?

    if ((__status__ == 0)); then
        # shellcheck disable=SC1090
        source "${__normalized__}/module.sh" || __status__=$?
    fi

    MODULE_LOADING_NAME=
    MODULE_LOADING_ROOT=

    if ((__status__ == 0)) && ! module::has "${__name__}"; then
        log::error "Module did not register itself: ${__name__}"
        __status__=1
    fi

    if ((__status__ != 0)); then
        module::_rollback __modules_before__ __commands_before__
        return "${__status__}"
    fi

    module::validate "${__name__}"
}

module::_source_implementations()
{
    local __root__=${1:-}
    local __file__

    while IFS= read -r -d '' __file__; do
        # shellcheck disable=SC1090
        source "${__file__}" || return
    done < <(
        find "${__root__}" \
            -type f \
            -name '*.sh' \
            ! -path "${__root__}/module.sh" \
            ! -path '*/tests/*' \
            -print0 |
            sort -z
    )
}
