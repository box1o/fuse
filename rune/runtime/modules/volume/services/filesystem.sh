#!/usr/bin/env bash

volume::check()
{
    if (($# != 1)); then
        log::error "Usage: rune volume check <name>"
        return 2
    fi
    volume::mounted "$1" && {
        log::error "Volume is mounted: $1"
        return 1
    }
    local __path__
    __path__=$(volume::path "$1") || return
    volume::_e2fsck "${__path__}"
}

volume::_e2fsck()
{
    e2fsck -pf -- "${1:-}"
    local __status__=$?
    ((__status__ == 0 || __status__ == 1))
}

volume::grow()
{
    args::reset
    args::option --name size --required
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune volume grow <name> --size SIZE"
        return 2
    fi

    local __name__
    local __requested__
    __name__=$(args::position 0) || return
    __requested__=$(args::get size) || return
    volume::_valid_name "${__name__}" || return

    local __lock__
    __lock__=$(runtime::lock_file volume "${__name__}") || return
    lock::with "${__lock__}" volume::_grow "${__name__}" "${__requested__}"
}

volume::_grow()
{
    local __name__=${1:-}
    local __requested__=${2:-}
    local __bytes__
    local __path__
    local __current__
    proc::require numfmt truncate e2fsck resize2fs losetup findmnt || return
    __bytes__=$(numfmt --from=iec "${__requested__^^}") || {
        log::error "Invalid volume size: ${__requested__}"
        return 2
    }
    __path__=$(volume::path "${__name__}") || return
    __current__=$(stat -c '%s' -- "${__path__}") || return

    if ((__bytes__ <= __current__)); then
        log::error "New size must be larger than current size: ${__current__} bytes"
        return 2
    fi
    if volume::mounted "${__name__}"; then
        log::error "Volume is mounted: ${__name__}"
        return 1
    fi

    volume::_e2fsck "${__path__}" || return
    truncate -s "${__bytes__}" -- "${__path__}" || return
    if ! resize2fs "${__path__}"; then
        log::error "Filesystem resize failed; backing file remains enlarged: ${__path__}"
        return 1
    fi
    volume::_e2fsck "${__path__}" || return
    yaml::set_raw "$(volume::_metadata "${__name__}")" '.size_bytes' "${__bytes__}" || return
    log::info "Grew volume: ${__name__}" "size=${__requested__}"
}
