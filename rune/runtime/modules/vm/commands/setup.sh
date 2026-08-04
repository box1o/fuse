#!/usr/bin/env bash

firecracker::setup()
{
    args::reset
    args::flag --name force
    args::parse "$@" || return

    if (($(args::count) > 0)); then
        log::error "firecracker setup: unexpected argument: $(args::position 0)"
        return 2
    fi

    local __runtime__
    local __downloads__
    __runtime__=$(module::var vm runtime) || return
    __downloads__=$(module::var vm downloads) || return
    fs::mkdir "${__runtime__}" || return
    fs::mkdir "$(module::var vm vms)" || return

    if args::has force || [[ ! -x "${__runtime__}/firecracker" || ! -x "${__runtime__}/jailer" ]]; then
        local __archive__
        __archive__=$(find "${__downloads__}" -maxdepth 1 -type f -name 'firecracker-*.tgz' -print 2>/dev/null | sort -V | tail -n 1)

        if [[ -z "${__archive__}" ]]; then
            firecracker::download || return
            __archive__=$(find "${__downloads__}" -maxdepth 1 -type f -name 'firecracker-*.tgz' -print | sort -V | tail -n 1)
        fi

        firecracker::_install_release "${__archive__}" "${__runtime__}" || return
    fi

    log::info "Firecracker is ready" "runtime=${__runtime__}"
}

firecracker::_install_release()
{
    local __archive__=${1:-}
    local __runtime__=${2:-}
    local __staging__="${__runtime__}/extract"

    fs::empty_dir "${__staging__}" || return
    archive::extract "${__archive__}" "${__staging__}" || return

    local __firecracker__
    local __jailer__
    __firecracker__=$(find "${__staging__}" -type f -name 'firecracker-v*-*' ! -name '*.debug' -print -quit)
    __jailer__=$(find "${__staging__}" -type f -name 'jailer-v*-*' ! -name '*.debug' -print -quit)

    fs::require_file firecracker-binary "${__firecracker__}" || return
    fs::require_file jailer-binary "${__jailer__}" || return
    fs::copy_file "${__firecracker__}" "${__runtime__}/firecracker" || return
    fs::copy_file "${__jailer__}" "${__runtime__}/jailer" || return
    fs::chmod 755 "${__runtime__}/firecracker" || return
    fs::chmod 755 "${__runtime__}/jailer" || return
    fs::remove_tree "${__staging__}"
}
