#!/usr/bin/env bash

module::_snapshot_registries()
{
    local -n _modules_=$1
    local -n _commands_=$2
    local __key__

    for __key__ in "${!MODULE_ROOT_REGISTRY[@]}"; do
        _modules_["${__key__}"]=1
    done

    for __key__ in "${!COMMAND_HANDLER_REGISTRY[@]}"; do
        _commands_["${__key__}"]=1
    done
}

module::_rollback()
{
    local -n _modules_before_=$1
    local -n _commands_before_=$2
    local __key__

    for __key__ in "${!COMMAND_HANDLER_REGISTRY[@]}"; do
        [[ -v "_commands_before_[${__key__}]" ]] || command::_remove "${__key__}"
    done

    for __key__ in "${!MODULE_ROOT_REGISTRY[@]}"; do
        [[ -v "_modules_before_[${__key__}]" ]] || module::_remove "${__key__}"
    done
}

module::load_all()
{
    local -A __seen__=()
    local __root__

    for __root__ in "$@"; do
        [[ -n "${__root__}" ]] || continue

        local __normalized__
        __normalized__=$(path::absolute "${__root__}") || return

        [[ -v "__seen__[${__normalized__}]" ]] && continue
        __seen__["${__normalized__}"]=1

        module::discover "${__normalized__}" || return
    done
}

module::log()
{
    local __name__
    while IFS= read -r __name__; do
        log::info \
            "${__name__}" \
            "root=${MODULE_ROOT_REGISTRY[${__name__}]}" \
            "state=${MODULE_STATE_REGISTRY[${__name__}]}"
    done < <(module::list)
}

module::_remove()
{
    local __name__=${1:-}
    local __context__=${MODULE_CONTEXT_REGISTRY[${__name__}]:-}

    if [[ -n "${__context__}" ]]; then
        unset "${__context__}"
    fi

    unset \
        'MODULE_ROOT_REGISTRY['"${__name__}"']' \
        'MODULE_SUMMARY_REGISTRY['"${__name__}"']' \
        'MODULE_DESCRIPTION_REGISTRY['"${__name__}"']' \
        'MODULE_VERSION_REGISTRY['"${__name__}"']' \
        'MODULE_API_REGISTRY['"${__name__}"']' \
        'MODULE_INIT_HANDLER_REGISTRY['"${__name__}"']' \
        'MODULE_SHUTDOWN_HANDLER_REGISTRY['"${__name__}"']' \
        'MODULE_CONTEXT_REGISTRY['"${__name__}"']' \
        'MODULE_STATE_REGISTRY['"${__name__}"']' \
        'MODULE_BUILTIN_REGISTRY['"${__name__}"']' \
        'MODULE_DEPENDENCY_REGISTRY['"${__name__}"']'
}

module::_reset()
{
    local __name__
    for __name__ in "${!MODULE_CONTEXT_REGISTRY[@]}"; do
        unset "${MODULE_CONTEXT_REGISTRY[${__name__}]}"
    done

    MODULE_ROOT_REGISTRY=()
    MODULE_SUMMARY_REGISTRY=()
    MODULE_DESCRIPTION_REGISTRY=()
    MODULE_VERSION_REGISTRY=()
    MODULE_API_REGISTRY=()
    MODULE_INIT_HANDLER_REGISTRY=()
    MODULE_SHUTDOWN_HANDLER_REGISTRY=()
    MODULE_CONTEXT_REGISTRY=()
    MODULE_STATE_REGISTRY=()
    MODULE_BUILTIN_REGISTRY=()
    MODULE_DEPENDENCY_REGISTRY=()
    MODULE_LOADING_NAME=
    MODULE_LOADING_ROOT=
}
