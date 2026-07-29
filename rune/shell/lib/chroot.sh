#!/usr/bin/env bash

# Validate a chroot filesystem. This does not make chroot a security sandbox.
chroot::validate()
{
    local __root__=${1:-}

    validate::absolute_path "${__root__}" || return
    fs::require_dir root-directory "${__root__}" || return

    chroot::_require_shell "${__root__}"
}

chroot::_require_shell()
{
    local __root__=${1:-}
    local __shell__="${__root__}/bin/sh"

    [[ -x "${__shell__}" ]] && return 0

    if [[ -L "${__shell__}" ]]; then
        local __target__
        __target__=$(readlink -- "${__shell__}") || return

        if [[ "${__target__}" == /* ]]; then
            __target__="${__root__}${__target__}"
        else
            __target__="${__root__}/bin/${__target__}"
        fi

        [[ -x "${__target__}" ]] && return 0
    fi

    log::error "No executable /bin/sh inside: ${__root__}"
    return 1
}

chroot::run()
{
    local __root__=${1:-}
    shift || true

    chroot::validate "${__root__}" || return
    user::require_root || return
    proc::require chroot env || return

    if (($# == 0)); then
        set -- /bin/sh
    fi

    command env -i \
        HOME=/root \
        PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
        TERM="${TERM:-xterm}" \
        chroot "${__root__}" "$@"
}

chroot::shell()
{
    chroot::run "${1:-}" /bin/sh
}
