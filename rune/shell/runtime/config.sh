#!/usr/bin/env bash

declare -gA RUNTIME_CONFIG_REGISTRY=()

runtime::config::has()
{
    local __key__=${1:-}
    validate::not_empty key "${__key__}" || return
    [[ -v "RUNTIME_CONFIG_REGISTRY[${__key__}]" ]]
}

runtime::config::set()
{
    local __key__=${1:-} __value__=${2-}
    validate::not_empty key "${__key__}" || return
    [[ "${__key__}" =~ ^[a-z][a-z0-9_-]*(\|[a-z][a-z0-9_-]*)+$ ]] || {
        log::error "Invalid runtime configuration key: ${__key__}"
        return 2
    }
    RUNTIME_CONFIG_REGISTRY["${__key__}"]=${__value__}
}

runtime::config::get()
{
    local __key__=${1:-} __default__=${2-}
    runtime::config::has "${__key__}" || {
        [[ $# -ge 2 ]] && {
            printf '%s\n' "${__default__}"
            return 0
        }
        log::error "Unknown runtime configuration key: ${__key__}"
        return 1
    }
    printf '%s\n' "${RUNTIME_CONFIG_REGISTRY[${__key__}]}"
}

runtime::config::_reset()
{
    RUNTIME_CONFIG_REGISTRY=()
}
