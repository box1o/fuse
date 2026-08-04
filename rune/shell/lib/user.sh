#!/usr/bin/env bash

user::current()
{
    id -un
}

user::uid()
{
    local __user__=${1:-}

    if [[ -n "${__user__}" ]]; then
        id -u "${__user__}"
    else
        id -u
    fi
}

user::gid()
{
    local __user__=${1:-}

    if [[ -n "${__user__}" ]]; then
        id -g "${__user__}"
    else
        id -g
    fi
}

user::exists()
{
    local __user__=${1:-}

    validate::not_empty user "${__user__}" || return
    id "${__user__}" >/dev/null 2>&1
}

user::group_exists()
{
    local __group__=${1:-}

    validate::not_empty group "${__group__}" || return
    getent group "${__group__}" >/dev/null 2>&1
}

user::is_root()
{
    ((EUID == 0))
}

user::require_root()
{
    if ! user::is_root; then
        log::error "This operation requires root privileges"
        return 1
    fi
}

user::require_non_root()
{
    if user::is_root; then
        log::error "This operation must not run as root"
        return 1
    fi
}

user::home()
{
    local __user__=${1:-}

    if [[ -z "${__user__}" ]]; then
        __user__=$(user::current) || return
    fi

    local __entry__
    __entry__=$(getent passwd "${__user__}") || return

    local __remainder__=${__entry__#*:*:*:*:*:}
    printf '%s\n' "${__remainder__%%:*}"
}

user::run_as()
{
    proc::run_as "$@"
}
