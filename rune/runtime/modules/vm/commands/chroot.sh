#!/usr/bin/env bash

firecracker::chroot()
{
    local __read_only__=false
    if [[ ${1:-} == --read-only ]]; then
        __read_only__=true
        shift
    fi

    local __name__=${1:-}
    [[ -n "${__name__}" ]] || {
        log::error "Usage: rune vm chroot [--read-only] <name> [-- command]"
        return 2
    }
    shift
    [[ ${1:-} == -- ]] && shift

    firecracker::_require_vm "${__name__}" || return
    firecracker::_require_stopped "${__name__}" || return
    user::require_root || return
    proc::require mount umount chroot || return

    local __vm__
    local __rootfs__
    local __mountpoint__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    __rootfs__=$(firecracker::_rootfs "${__name__}") || return
    __mountpoint__="${__vm__}/rootfs"
    fs::mkdir "${__mountpoint__}" || return

    local -a __mount_args__=(-o loop)
    [[ "${__read_only__}" == true ]] && __mount_args__=(-o 'loop,ro')
    mount "${__mount_args__[@]}" -- "${__rootfs__}" "${__mountpoint__}" || return

    local __status__=0
    if (($# == 0)); then
        chroot::shell "${__mountpoint__}" || __status__=$?
    else
        chroot::run "${__mountpoint__}" "$@" || __status__=$?
    fi

    umount -- "${__mountpoint__}" || return
    fs::remove_dir "${__mountpoint__}" || true
    return "${__status__}"
}
