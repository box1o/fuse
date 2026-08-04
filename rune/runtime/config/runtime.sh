#!/usr/bin/env bash

runtime::register_paths() {
    local __project_root__=${1:-}
    validate::not_empty project-root "${__project_root__}" || return
    validate::absolute_path "${__project_root__}" || return

    local __cli_name__
    local __user_home__=${HOME}
    __cli_name__=$(runtime::config::get 'cli|name') || return

    if (( EUID == 0 )) && [[ -n "${SUDO_USER:-}" ]]; then
        __user_home__=$(user::home "${SUDO_USER}") || return
    fi

    runtime::_set_path 'runtime|project_root' "${__project_root__}" 'Project root directory' || return
    runtime::_set_path 'runtime|shell' "${__project_root__}/shell/runtime" 'Shell runtime directory' || return
    runtime::_set_path 'runtime|modules' "${__project_root__}/runtime/modules" 'Module directory' || return
    runtime::_set_path 'runtime|templates' "${__project_root__}/shell/templates/module" 'Module template directory' || return
    runtime::_set_path 'runtime|config' "${RUNE_CONFIG_FILE:-${__project_root__}/rune.yaml}" 'Project configuration file' || return
    runtime::_set_path 'runtime|cache' "${XDG_CACHE_HOME:-${__user_home__}/.cache}/${__cli_name__}" 'Runtime cache directory' || return
    runtime::_set_path 'runtime|state' "${XDG_STATE_HOME:-${__user_home__}/.local/state}/${__cli_name__}" 'Runtime state directory'
}

runtime::_set_path() {
    local __key__=${1:-} __value__=${2:-} __description__=${3:-}
    if paths::has "${__key__}"; then
        paths::set --source runtime --description "${__description__}" "${__key__}" "${__value__}"
    else
        paths::register --source runtime --description "${__description__}" "${__key__}" "${__value__}"
    fi
}

runtime::lock_file()
{
    local __namespace__=${1:-}
    local __key__=${2:-}
    if [[ ! "${__namespace__}" =~ ^[a-z][a-z0-9-]*$ || ! "${__key__}" =~ ^[a-zA-Z0-9][a-zA-Z0-9._-]*$ ]]; then
        log::error "Invalid runtime lock identifier: ${__namespace__}/${__key__}"
        return 2
    fi

    local __state__
    __state__=$(paths::require 'runtime|state') || return
    printf '%s/locks/%s/%s.lock\n' "${__state__}" "${__namespace__}" "${__key__}"
}
