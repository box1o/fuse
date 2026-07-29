#!/usr/bin/env bash

mount::is_mounted()
{
    local __mountpoint__=${1:-}

    validate::not_empty mountpoint "${__mountpoint__}" || return
    proc::require mountpoint || return

    mountpoint -q -- "${__mountpoint__}"
}

mount::bind()
{
    local __source__=${1:-}
    local __target__=${2:-}

    mount::_prepare "${__source__}" "${__target__}" || return
    mount --bind "${__source__}" "${__target__}"
}

mount::bind_recursive()
{
    local __source__=${1:-}
    local __target__=${2:-}

    mount::_prepare "${__source__}" "${__target__}" || return
    mount --rbind "${__source__}" "${__target__}" || return
    mount --make-rslave "${__target__}"
}

mount::_prepare()
{
    local __source__=${1:-}
    local __target__=${2:-}

    user::require_root || return
    proc::require mount || return
    fs::_require_source "${__source__}" || return
    validate::not_empty mount-target "${__target__}" || return
    fs::mkdir "${__target__}"
}

mount::tmpfs()
{
    local __target__=${1:-}

    user::require_root || return
    proc::require mount || return
    fs::mkdir "${__target__}" || return

    mount -t tmpfs tmpfs "${__target__}"
}

mount::unmount()
{
    mount::_unmount false "${1:-}"
}

mount::unmount_recursive()
{
    mount::_unmount true "${1:-}"
}

mount::_unmount()
{
    local __recursive__=${1:-false}
    local __target__=${2:-}

    user::require_root || return
    proc::require umount || return
    mount::is_mounted "${__target__}" || return 0

    if [[ "${__recursive__}" == true ]]; then
        umount -R -- "${__target__}"
    else
        umount -- "${__target__}"
    fi
}
