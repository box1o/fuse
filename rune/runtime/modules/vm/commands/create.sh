#!/usr/bin/env bash

firecracker::create()
{
    args::reset
    args::option --name cpu --default 1
    args::option --name memory --default 512
    args::option --name kernel
    args::option --name rootfs
    args::option --name image
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune vm create <name> [--image NAME[:VERSION]] [--cpu N] [--memory MiB] [--kernel PATH] [--rootfs PATH]"
        return 2
    fi

    local __name__
    local __cpu__
    local __memory__
    local __kernel__
    local __rootfs__
    local __image__
    __name__=$(args::position 0) || return
    __cpu__=$(args::get cpu) || return
    __memory__=$(args::get memory) || return
    __kernel__=$(args::get kernel '') || return
    __rootfs__=$(args::get rootfs '') || return
    __image__=$(args::get image '') || return

    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_create \
        "${__name__}" "${__cpu__}" "${__memory__}" "${__kernel__}" "${__rootfs__}" "${__image__}"
}

firecracker::_create()
{
    local __name__=${1:-}
    local __cpu__=${2:-}
    local __memory__=${3:-}
    local __kernel__=${4:-}
    local __rootfs__=${5:-}
    local __image__=${6:-}

    firecracker::_valid_vm_name "${__name__}" || return
    firecracker::_validate_cpu "${__cpu__}" || return
    resource::validate_memory "${__memory__}" || return

    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    if [[ -e "${__vm__}" ]]; then
        log::error "VM already exists: ${__name__}"
        return 1
    fi

    firecracker::_is_setup || firecracker::setup || return

    if [[ -n "${__image__}" ]]; then
        module::use image || return
        local __config__
        local __ensure__
        __config__=$(paths::require 'runtime|config') || return
        if yaml::has "${__config__}" '.runtime.vm.ensure_images'; then
            __ensure__=$(yaml::get "${__config__}" '.runtime.vm.ensure_images') || return
        else
            __ensure__=$(yaml::get "${__config__}" '.runtime.firecracker.ensure_images' true) || return
        fi
        if [[ "${__ensure__}" == true ]]; then
            __image__=$(image::ensure "${__image__}" --file "${__config__}") || return
        fi
        __rootfs__=$(image::resolve "${__image__}" rootfs_path) || return
        __kernel__=$(image::resolve "${__image__}" kernel_path) || return
    elif [[ -z "${__kernel__}" || -z "${__rootfs__}" ]]; then
        log::error "vm create requires --image or both --kernel and --rootfs"
        return 2
    fi
    __kernel__=$(path::absolute "${__kernel__}") || return
    __rootfs__=$(path::absolute "${__rootfs__}") || return
    fs::require_file kernel "${__kernel__}" || return
    fs::require_file rootfs "${__rootfs__}" || return

    local __volume__="vm-${__name__}-root"
    volume::create "${__volume__}" --from "${__rootfs__}" --owner "vm:${__name__}" >/dev/null || return

    fs::mkdir "${__vm__}" || {
        volume::delete "${__volume__}" --force >/dev/null 2>&1 || true
        return 1
    }

    firecracker::_write_metadata \
        "${__name__}" "${__volume__}" "${__kernel__}" "${__cpu__}" "${__memory__}" || {
        fs::remove_tree "${__vm__}"
        volume::delete "${__volume__}" --force >/dev/null 2>&1 || true
        return 1
    }
    firecracker::_render_config "${__name__}" || {
        fs::remove_tree "${__vm__}"
        volume::delete "${__volume__}" --force >/dev/null 2>&1 || true
        return 1
    }

    printf '%s\n' "${__kernel__}" | fs::atomic_write "${__vm__}/kernel.path"
    log::info "Created VM: ${__name__}" "cpu=${__cpu__}" "memory=${__memory__}MiB"
}
