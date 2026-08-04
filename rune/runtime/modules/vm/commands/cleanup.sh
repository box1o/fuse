#!/usr/bin/env bash

firecracker::cleanup()
{
    args::reset
    args::flag --name all
    args::flag --name logs
    args::parse "$@" || return

    if args::has all; then
        if (($(args::count) != 0)); then
            log::error "Usage: rune vm cleanup <name> [--logs] | --all [--logs]"
            return 2
        fi
        local __vms__
        local __directory__
        __vms__=$(module::var vm vms) || return
        [[ -d "${__vms__}" ]] || return 0
        while IFS= read -r -d '' __directory__; do
            firecracker::_cleanup_one "${__directory__##*/}" "$(args::has logs && printf true || printf false)" || return
        done < <(find "${__vms__}" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)
        return 0
    fi

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm cleanup <name> [--logs] | --all [--logs]"
        return 2
    fi
    firecracker::_cleanup_one \
        "$(args::position 0)" \
        "$(args::has logs && printf true || printf false)"
}

firecracker::_cleanup_one()
{
    local __name__=${1:-}
    local __logs__=${2:-false}
    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_cleanup_locked "${__name__}" "${__logs__}"
}

firecracker::_cleanup_locked()
{
    local __name__=${1:-}
    local __logs__=${2:-false}
    local __vm__
    firecracker::_require_vm "${__name__}" || return
    firecracker::_require_stopped "${__name__}" || return
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    fs::remove_file "${__vm__}/firecracker.pid" || return
    fs::remove_file "${__vm__}/firecracker.sock" || return
    if [[ "${__logs__}" == true && -f "${__vm__}/console.log" ]]; then
        truncate -s 0 -- "${__vm__}/console.log" || return
    fi
    log::info "Cleaned VM runtime files: ${__name__}"
}
