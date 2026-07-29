#!/usr/bin/env bash

firecracker::resize()
{
    args::reset
    args::option --name cpu
    args::option --name memory
    args::option --name disk
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm resize <name> [--cpu N] [--memory MiB] [--disk SIZE]"
        return 2
    fi
    if ! args::has cpu && ! args::has memory && ! args::has disk; then
        log::error "firecracker resize requires --cpu, --memory, or --disk"
        return 2
    fi

    local __name__
    local __cpu__=
    local __memory__=
    local __disk__=
    __name__=$(args::position 0) || return
    args::has cpu && __cpu__=$(args::get cpu)
    args::has memory && __memory__=$(args::get memory)
    args::has disk && __disk__=$(args::get disk)

    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_resize "${__name__}" "${__cpu__}" "${__memory__}" "${__disk__}"
}

firecracker::_resize()
{
    local __name__=${1:-}
    local __cpu__=${2:-}
    local __memory__=${3:-}
    local __disk__=${4:-}
    local __vm__
    local __metadata__
    firecracker::_require_vm "${__name__}" || return
    firecracker::_require_stopped "${__name__}" || return
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    __metadata__="${__vm__}/metadata.yaml"

    [[ -z "${__cpu__}" ]] || firecracker::_validate_cpu "${__cpu__}" || return
    [[ -z "${__memory__}" ]] || resource::validate_memory "${__memory__}" || return

    if [[ -n "${__disk__}" ]]; then
        if [[ ! -f "${__metadata__}" ]] || ! yaml::has "${__metadata__}" '.root_volume'; then
            log::error "Disk resize requires a managed volume; recreate this legacy VM from an image"
            return 1
        fi
        volume::grow \
            "$(yaml::get "${__metadata__}" '.root_volume')" \
            --size "${__disk__}" || return
    fi

    if [[ -f "${__metadata__}" ]]; then
        if [[ -n "${__cpu__}" || -n "${__memory__}" ]]; then
            firecracker::_update_metadata_resources "${__name__}" "${__cpu__}" "${__memory__}" || return
            firecracker::_render_config "${__name__}" || return
        fi
    else
        [[ -z "${__cpu__}" ]] ||
            firecracker::_set_config_resource "${__vm__}/config.json" '.machine-config.vcpu_count' "${__cpu__}" || return
        [[ -z "${__memory__}" ]] ||
            firecracker::_set_config_resource "${__vm__}/config.json" '.machine-config.mem_size_mib' "${__memory__}" || return
    fi

    log::info "Resized VM: ${__name__}"
}

firecracker::set_memory()
{
    if (($# != 2)); then
        log::error "Usage: rune vm set-memory <name> <MiB>"
        return 2
    fi
    firecracker::resize "$1" --memory "$2"
}
