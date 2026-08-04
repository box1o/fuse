#!/usr/bin/env bash

paths::register_defaults()
{
    local __home__=${HOME:-}

    if [[ -z ${__home__} ]]; then
        local __account__ __password__ __uid__ __gid__ __gecos__ __shell__

        while IFS=: read -r __account__ __password__ __uid__ __gid__ __gecos__ __home__ __shell__; do
            [[ ${__uid__} == "${EUID}" ]] && break
            __home__=
        done </etc/passwd
    fi

    if [[ -z ${__home__} ]]; then
        log::error "Unable to resolve the current user's home directory"
        return 1
    fi

    paths::_register_default system\|tmp "${TMPDIR:-/tmp}" "System temporary directory" || return
    paths::_register_default user\|home "${__home__}" "Current user home directory" || return
    paths::_register_default user\|cache "${XDG_CACHE_HOME:-${__home__}/.cache}" "User cache directory" || return
    paths::_register_default user\|config "${XDG_CONFIG_HOME:-${__home__}/.config}" "User configuration directory" || return
    paths::_register_default user\|data "${XDG_DATA_HOME:-${__home__}/.local/share}" "User data directory" || return
    paths::_register_default user\|state "${XDG_STATE_HOME:-${__home__}/.local/state}" "User state directory"
}

paths::_register_default()
{
    local __key__=${1:-}
    local __value__=${2:-}
    local __description__=${3:-}

    paths::register \
        --source default \
        --description "${__description__}" \
        "${__key__}" \
        "${__value__}"
}
