#!/usr/bin/env bash

firecracker::start()
{
    if (($# != 1)); then
        log::error "Usage: rune vm start <name>"
        return 2
    fi

    local __name__=$1
    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_start "${__name__}"
}

firecracker::_start()
{
    local __name__=${1:-}
    firecracker::_require_vm "${__name__}" || return
    if firecracker::_is_running "${__name__}"; then
        log::info "VM is already running: ${__name__}"
        return 0
    fi

    local __runtime__
    local __kvm__
    local __vm__
    __runtime__=$(module::var vm runtime) || return
    __kvm__=$(module::var vm kvm) || return
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    fs::require_file firecracker "${__runtime__}/firecracker" || return
    proc::require nohup || return

    if [[ ! -r "${__kvm__}" || ! -w "${__kvm__}" ]]; then
        log::error "KVM is not accessible: ${__kvm__}"
        return 1
    fi

    fs::remove_file "${__vm__}/firecracker.pid" || return
    fs::remove_file "${__vm__}/firecracker.sock" || return

    nohup "${__runtime__}/firecracker" \
        --api-sock "${__vm__}/firecracker.sock" \
        --config-file "${__vm__}/config.json" \
        </dev/null >>"${__vm__}/console.log" 2>&1 &
    local __pid__=$!
    printf '%s\n' "${__pid__}" | fs::atomic_write "${__vm__}/firecracker.pid"

    local __attempt__
    for ((__attempt__ = 0; __attempt__ < 20; __attempt__++)); do
        if [[ -S "${__vm__}/firecracker.sock" ]]; then
            log::info "Started VM: ${__name__}" "pid=${__pid__}"
            return 0
        fi
        if ! kill -0 "${__pid__}" 2>/dev/null; then
            log::error "Firecracker exited while starting: ${__name__}"
            tail -n 20 "${__vm__}/console.log" >&2 || true
            return 1
        fi
        sleep 0.1
    done

    log::error "Timed out waiting for VM socket: ${__name__}"
    firecracker::_stop "${__name__}" >/dev/null 2>&1 || true
    return 1
}
