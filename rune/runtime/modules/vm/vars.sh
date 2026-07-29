#!/usr/bin/env bash

FIRECRACKER_DEFAULT_VERSION=${FIRECRACKER_DEFAULT_VERSION:-latest}
FIRECRACKER_RELEASES_URL=${FIRECRACKER_RELEASES_URL:-https://github.com/firecracker-microvm/firecracker/releases}
FIRECRACKER_KERNEL_URL=${FIRECRACKER_KERNEL_URL:-https://s3.amazonaws.com/spec.ccfc.min/img/hello/kernel/hello-vmlinux.bin}
FIRECRACKER_ROOTFS_URL=${FIRECRACKER_ROOTFS_URL:-https://s3.amazonaws.com/spec.ccfc.min/img/hello/fsfiles/hello-rootfs.ext4}
FIRECRACKER_KVM_PATH=${FIRECRACKER_KVM_PATH:-/dev/kvm}

firecracker::_init()
{
    local -n _firecracker_=${1}

    local __cache__
    local __runtime_cache__
    __runtime_cache__=$(paths::require 'runtime|cache') || return
    __cache__=${FIRECRACKER_CACHE_DIR:-"${__runtime_cache__}/firecracker"}
    __cache__=$(path::absolute "${__cache__}") || return

    var::update _firecracker_ cache "${__cache__}"
    var::update _firecracker_ downloads "${__cache__}/downloads"
    var::update _firecracker_ runtime "${__cache__}/runtime"
    var::update _firecracker_ assets "${__cache__}/assets"
    var::update _firecracker_ vms "${__cache__}/vms"
    var::update _firecracker_ default_version "${FIRECRACKER_DEFAULT_VERSION}"
    var::update _firecracker_ releases_url "${FIRECRACKER_RELEASES_URL}"
    var::update _firecracker_ kernel_url "${FIRECRACKER_KERNEL_URL}"
    var::update _firecracker_ rootfs_url "${FIRECRACKER_ROOTFS_URL}"
    var::update _firecracker_ kvm "${FIRECRACKER_KVM_PATH}"
}
