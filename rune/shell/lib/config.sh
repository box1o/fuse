#!/usr/bin/env bash

declare -gA CONFIG_FILE_REGISTRY=()

config::_valid_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9_-]*$ ]]; then
        log::error "Invalid configuration name: ${1:-}"
        return 2
    fi
}

config::load()
{
    local __name__=${1:-}
    local __file__=${2:-}
    local __version__=${3:-}

    config::_valid_name "${__name__}" || return
    validate::not_empty configuration-file "${__file__}" || return
    __file__=$(path::absolute "${__file__}") || return
    yaml::validate "${__file__}" || return

    if [[ -n "${__version__}" ]]; then
        local __actual__
        __actual__=$(yaml::get "${__file__}" '.version') || return
        if [[ "${__actual__}" != "${__version__}" ]]; then
            log::error "Unsupported configuration version: ${__actual__}; expected ${__version__}"
            return 1
        fi
    fi

    CONFIG_FILE_REGISTRY["${__name__}"]=${__file__}
}

config::has_document()
{
    [[ -v "CONFIG_FILE_REGISTRY[${1:-}]" ]]
}

config::file()
{
    local __name__=${1:-}
    if ! config::has_document "${__name__}"; then
        log::error "Configuration is not loaded: ${__name__}"
        return 1
    fi
    printf '%s\n' "${CONFIG_FILE_REGISTRY[${__name__}]}"
}

config::get()
{
    local __name__=${1:-}
    local __expression__=${2:-.}
    local __file__
    __file__=$(config::file "${__name__}") || return

    if (($# >= 3)); then
        yaml::get "${__file__}" "${__expression__}" "${3-}"
    else
        yaml::get "${__file__}" "${__expression__}"
    fi
}

config::has()
{
    local __file__
    __file__=$(config::file "${1:-}") || return
    yaml::has "${__file__}" "${2:-}"
}

config::require()
{
    local __name__=${1:-}
    shift || true
    local __file__
    __file__=$(config::file "${__name__}") || return
    yaml::require "${__file__}" "$@"
}

config::keys()
{
    local __file__
    __file__=$(config::file "${1:-}") || return
    yaml::keys "${__file__}" "${2:-.}"
}

config::values()
{
    local __file__
    __file__=$(config::file "${1:-}") || return
    yq eval -r "${2:-.}[]?" "${__file__}"
}

config::entries()
{
    local __file__
    __file__=$(config::file "${1:-}") || return
    yq eval -o=json -I=0 "${2:-.} | to_entries[]?" "${__file__}"
}

config::reset()
{
    CONFIG_FILE_REGISTRY=()
}
