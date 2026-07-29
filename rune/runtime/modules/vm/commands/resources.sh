#!/usr/bin/env bash

firecracker::resources()
{
    if (($# != 1)); then
        log::error "Usage: rune vm resources <name>"
        return 2
    fi

    local __name__=$1
    local __rootfs__
    local __disk__
    firecracker::_require_vm "${__name__}" || return
    __rootfs__=$(firecracker::_rootfs "${__name__}") || return
    fs::require_file rootfs "${__rootfs__}" || return
    __disk__=$(stat -c '%s' -- "${__rootfs__}") || return

    local __state__=stopped
    local __cpu__
    local __memory__
    firecracker::_is_running "${__name__}" && __state__=running
    __cpu__=$(firecracker::_config_value "${__name__}" '.machine-config.vcpu_count') || return
    __memory__=$(firecracker::_config_value "${__name__}" '.machine-config.mem_size_mib') || return

    if [[ $(runtime::config::get 'output|format' text) == json ]]; then
        printf '{"name":"%s","status":"%s","cpu":%s,"memory_mib":%s,"disk_bytes":%s,"rootfs":"%s"}\n' \
            "$(firecracker::_json_string "${__name__}")" "${__state__}" \
            "${__cpu__}" "${__memory__}" "${__disk__}" \
            "$(firecracker::_json_string "${__rootfs__}")"
        return
    fi

    printf 'name: %s\n' "${__name__}"
    printf 'status: %s\n' "${__state__}"
    printf 'cpu: %s\n' "${__cpu__}"
    printf 'memory_mib: %s\n' "${__memory__}"
    printf 'disk_bytes: %s\n' "${__disk__}"
    printf 'disk: %s\n' "$(numfmt --to=iec-i --suffix=B "${__disk__}")"
    printf 'rootfs: %s\n' "${__rootfs__}"
}

firecracker::_validate_cpu()
{
    local __cpu__=${1:-}
    resource::validate_cpu "${__cpu__}" || return
    if ((__cpu__ != 1 && (__cpu__ % 2 != 0 || __cpu__ > 32))); then
        log::error "Firecracker CPU count must be 1 or an even number up to 32: ${__cpu__}"
        return 2
    fi
}
