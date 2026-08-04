#!/usr/bin/env bash

firecracker::_write_metadata()
{
    local __name__=${1:-}
    local __volume__=${2:-}
    local __kernel__=${3:-}
    local __cpu__=${4:-}
    local __memory__=${5:-}
    local __metadata__
    local __document__
    __metadata__=$(firecracker::_metadata "${__name__}") || return

    __document__=$(VM_NAME=${__name__} \
        VM_VOLUME=${__volume__} \
        VM_KERNEL=${__kernel__} \
        VM_CPU=${__cpu__} \
        VM_MEMORY=${__memory__} \
        yq -n -o=yaml '
            .schema = 1 |
            .name = strenv(VM_NAME) |
            .provider = "firecracker" |
            .kernel_path = strenv(VM_KERNEL) |
            .root_volume = strenv(VM_VOLUME) |
            .cpu = (strenv(VM_CPU) | tonumber) |
            .memory_mib = (strenv(VM_MEMORY) | tonumber)
        ') || return
    printf '%s\n' "${__document__}" | fs::atomic_write "${__metadata__}"
}

firecracker::_render_config()
{
    local __name__=${1:-}
    local __vm__
    local __metadata__
    local __kernel__
    local __volume__
    local __rootfs__
    local __cpu__
    local __memory__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    __metadata__="${__vm__}/metadata.yaml"
    fs::require_file VM-metadata "${__metadata__}" || return
    __kernel__=$(yaml::get "${__metadata__}" '.kernel_path') || return
    __volume__=$(yaml::get "${__metadata__}" '.root_volume') || return
    __cpu__=$(yaml::get "${__metadata__}" '.cpu') || return
    __memory__=$(yaml::get "${__metadata__}" '.memory_mib') || return
    __rootfs__=$(volume::path "${__volume__}") || return
    fs::require_file kernel "${__kernel__}" || return
    firecracker::_validate_cpu "${__cpu__}" || return
    resource::validate_memory "${__memory__}" || return

    local __kernel_json__
    local __rootfs_json__
    __kernel_json__=$(firecracker::_json_string "${__kernel__}")
    __rootfs_json__=$(firecracker::_json_string "${__rootfs__}")
    fs::atomic_write "${__vm__}/config.json" <<EOF
{
    "boot-source": {
        "kernel_image_path": "${__kernel_json__}",
        "boot_args": "console=ttyS0 reboot=k panic=1 pci=off"
    },
    "drives": [{
        "drive_id": "rootfs",
        "path_on_host": "${__rootfs_json__}",
        "is_root_device": true,
        "is_read_only": false
    }],
    "machine-config": {
        "vcpu_count": ${__cpu__},
        "mem_size_mib": ${__memory__},
        "smt": false
    }
}
EOF
}

firecracker::_update_metadata_resources()
{
    local __name__=${1:-}
    local __cpu__=${2:-}
    local __memory__=${3:-}
    local __metadata__
    local __temporary__
    __metadata__=$(firecracker::_metadata "${__name__}") || return
    fs::require_file VM-metadata "${__metadata__}" || return
    __temporary__=$(temp::file vm-metadata) || return
    fs::copy_file "${__metadata__}" "${__temporary__}" || return
    [[ -z "${__cpu__}" ]] || yaml::set_raw "${__temporary__}" '.cpu' "${__cpu__}" || return
    [[ -z "${__memory__}" ]] || yaml::set_raw "${__temporary__}" '.memory_mib' "${__memory__}" || return
    fs::move "${__temporary__}" "${__metadata__}" || return
}

firecracker::_config_value()
{
    local __name__=${1:-}
    local __expression__=${2:-}
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    yaml::get "${__vm__}/config.json" "${__expression__}"
}

firecracker::_set_config_resource()
{
    local __config__=${1:-}
    local __expression__=${2:-}
    local __value__=${3:-}
    local __temporary__
    __temporary__=$(temp::file firecracker-config) || return
    CONFIG_VALUE=${__value__} yq eval -o=json -I=4 \
        "${__expression__} = (strenv(CONFIG_VALUE) | tonumber)" \
        "${__config__}" >"${__temporary__}" || return
    fs::move "${__temporary__}" "${__config__}"
}
