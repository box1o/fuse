#!/usr/bin/env bash

firecracker::logs()
{
    args::reset
    args::option --name lines --default 100
    args::flag --name follow --short f
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm logs <name> [--lines N] [--follow]"
        return 2
    fi

    local __name__
    local __lines__
    local __vm__
    __name__=$(args::position 0) || return
    __lines__=$(args::get lines) || return
    validate::positive_integer lines "${__lines__}" || return
    firecracker::_require_vm "${__name__}" || return
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    fs::touch "${__vm__}/console.log" || return

    if args::has follow; then
        tail -n "${__lines__}" -f -- "${__vm__}/console.log"
    else
        tail -n "${__lines__}" -- "${__vm__}/console.log"
    fi
}

firecracker::run()
{
    args::reset
    args::option --name cpu --default 1
    args::option --name memory --default 512
    args::option --name image --required
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm run <name> [--cpu N] [--memory MiB]"
        return 2
    fi

    local __name__
    local __cpu__
    local __memory__
    local __image__
    __name__=$(args::position 0) || return
    __cpu__=$(args::get cpu) || return
    __memory__=$(args::get memory) || return
    __image__=$(args::get image) || return

    firecracker::setup || return
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    if [[ ! -d "${__vm__}" ]]; then
        firecracker::create "${__name__}" --image "${__image__}" --cpu "${__cpu__}" --memory "${__memory__}" || return
    fi
    firecracker::start "${__name__}"
}
