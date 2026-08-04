#!/usr/bin/env bash

firecracker::stop()
{
    if (($# != 1)); then
        log::error "Usage: rune vm stop <name>"
        return 2
    fi

    local __name__=$1
    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_stop "${__name__}"
}

firecracker::_stop()
{
    local __name__=${1:-}
    firecracker::_require_vm "${__name__}" || return

    if ! firecracker::_is_running "${__name__}"; then
        log::info "VM is already stopped: ${__name__}"
        return 0
    fi

    local __pid__
    local __vm__
    __pid__=$(firecracker::_pid "${__name__}") || return
    __vm__=$(firecracker::_vm_dir "${__name__}") || return

    firecracker::_request_shutdown "${__vm__}"

    local __attempt__
    for ((__attempt__ = 0; __attempt__ < 100; __attempt__++)); do
        kill -0 "${__pid__}" 2>/dev/null || break
        sleep 0.1
    done

    if kill -0 "${__pid__}" 2>/dev/null; then
        kill -TERM "${__pid__}" || return
        for ((__attempt__ = 0; __attempt__ < 50; __attempt__++)); do
            kill -0 "${__pid__}" 2>/dev/null || break
            sleep 0.1
        done
    fi
    if kill -0 "${__pid__}" 2>/dev/null; then
        log::warn "VM did not stop gracefully; forcing process exit: ${__name__}"
        kill -KILL "${__pid__}" || return
    fi

    fs::remove_file "${__vm__}/firecracker.pid" || return
    fs::remove_file "${__vm__}/firecracker.sock" || return
    log::info "Stopped VM: ${__name__}"
}

firecracker::_request_shutdown()
{
    local __vm__=${1:-}
    local __socket__="${__vm__}/firecracker.sock"
    [[ -S "${__socket__}" ]] || return 0
    command -v curl >/dev/null 2>&1 || return 0
    curl --silent --show-error --fail \
        --unix-socket "${__socket__}" \
        -X PUT \
        -H 'Content-Type: application/json' \
        --data '{"action_type":"SendCtrlAltDel"}' \
        http://localhost/actions >/dev/null 2>&1 || true
}

firecracker::delete()
{
    args::reset
    args::flag --name force
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm delete <name> [--force]"
        return 2
    fi

    local __name__
    local __force__=false
    __name__=$(args::position 0) || return
    args::has force && __force__=true
    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_delete "${__name__}" "${__force__}"
}

firecracker::_delete()
{
    local __name__=${1:-}
    local __force__=${2:-false}
    local __vm__
    firecracker::_require_vm "${__name__}" || return

    if firecracker::_is_running "${__name__}"; then
        if [[ "${__force__}" != true ]]; then
            log::error "VM is running; stop it first or use --force: ${__name__}"
            return 1
        fi
        firecracker::_stop "${__name__}" || return
    fi

    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    local __volume__=
    if [[ -f "${__vm__}/metadata.yaml" ]] && yaml::has "${__vm__}/metadata.yaml" '.root_volume'; then
        __volume__=$(yaml::get "${__vm__}/metadata.yaml" '.root_volume') || return
    fi
    local __vms__
    __vms__=$(module::var vm vms) || return
    path::assert_within "${__vm__}" "${__vms__}" || return
    fs::remove_tree "${__vm__}" || return
    if [[ -n "${__volume__}" ]]; then
        volume::delete "${__volume__}" --force || return
    fi
    log::info "Deleted VM: ${__name__}"
}
