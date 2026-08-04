#!/usr/bin/env bash

firecracker::_valid_vm_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid VM name: ${1:-}"
        return 2
    fi
}

firecracker::_vm_dir()
{
    local __name__=${1:-}
    firecracker::_valid_vm_name "${__name__}" || return

    local __vms__
    __vms__=$(module::var vm vms) || return
    printf '%s/%s\n' "${__vms__}" "${__name__}"
}

firecracker::_require_vm()
{
    local __name__=${1:-}
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return

    if [[ ! -d "${__vm__}" ]]; then
        log::error "VM does not exist: ${__name__}"
        return 1
    fi
}

firecracker::_pid()
{
    local __name__=${1:-}
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    [[ -f "${__vm__}/firecracker.pid" ]] || return 1
    read -r __pid__ <"${__vm__}/firecracker.pid"
    printf '%s\n' "${__pid__}"
}

firecracker::_is_running()
{
    local __pid__
    __pid__=$(firecracker::_pid "${1:-}") || return 1
    [[ "${__pid__}" =~ ^[0-9]+$ ]] || return 1
    kill -0 "${__pid__}" 2>/dev/null || return 1

    local __vm__
    local __command__
    __vm__=$(firecracker::_vm_dir "${1:-}") || return 1
    [[ -r "/proc/${__pid__}/cmdline" ]] || return 1
    __command__=$(tr '\0' ' ' <"/proc/${__pid__}/cmdline") || return 1
    [[ " ${__command__} " == *" --config-file ${__vm__}/config.json "* ]]
}

firecracker::_require_stopped()
{
    if firecracker::_is_running "${1:-}"; then
        log::error "VM is running: ${1}"
        return 1
    fi
}

firecracker::_json_string()
{
    local __value__=${1-}
    __value__=${__value__//\\/\\\\}
    __value__=${__value__//\"/\\\"}
    printf '%s' "${__value__}"
}

firecracker::_is_setup()
{
    local __runtime__
    __runtime__=$(module::var vm runtime) || return
    [[ -x "${__runtime__}/firecracker" && -x "${__runtime__}/jailer" ]]
}

firecracker::_metadata()
{
    printf '%s/metadata.yaml\n' "$(firecracker::_vm_dir "${1:-}")"
}

firecracker::_rootfs()
{
    local __name__=${1:-}
    local __vm__
    local __metadata__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    __metadata__="${__vm__}/metadata.yaml"

    if [[ -f "${__metadata__}" ]] && yaml::has "${__metadata__}" '.root_volume'; then
        volume::path "$(yaml::get "${__metadata__}" '.root_volume')"
    else
        printf '%s/rootfs.ext4\n' "${__vm__}"
    fi
}

firecracker::_lock()
{
    local __name__=${1:-}
    firecracker::_valid_vm_name "${__name__}" || return
    runtime::lock_file vm "${__name__}"
}
