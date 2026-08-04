#!/usr/bin/env bash

image::build()
{
    local __file__=
    local __reference__=
    local __force__=false
    local __dry_run__=false
    while (($# > 0)); do
        case "$1" in
            --file)
                [[ $# -ge 2 ]] || return 2
                __file__=$2
                shift 2
                ;;
            --force)
                __force__=true
                shift
                ;;
            --dry-run)
                __dry_run__=true
                shift
                ;;
            -*)
                log::error "image build: unknown option: $1"
                return 2
                ;;
            *)
                [[ -z "${__reference__}" ]] || return 2
                __reference__=$1
                shift
                ;;
        esac
    done

    local __name__
    local __requested_version__
    image::_split_reference "${__reference__}" __name__ __requested_version__ || return
    if [[ "${__dry_run__}" == true ]]; then
        image::_build "${__name__}" "${__requested_version__}" "${__file__}" "${__force__}" true
        return
    fi
    local __lock__
    __lock__=$(runtime::lock_file image "${__name__}") || return
    lock::with "${__lock__}" image::_build \
        "${__name__}" "${__requested_version__}" "${__file__}" "${__force__}" "${__dry_run__}"
}

image::_build()
{
    local __name__=${1:-}
    local __requested_version__=${2:-}
    local __file__=${3:-}
    local __force__=${4:-false}
    local __dry_run__=${5:-false}
    image::_load_config __file__ "${__file__}" || return
    image::_validate_definition "${__file__}" "${__name__}" || return

    local __expression__
    local __version__
    local __architecture__
    local __size__
    local __release__
    local __kernel__
    __expression__=$(image::_expression "${__name__}") || return
    __version__=$(yaml::get "${__file__}" "${__expression__}.version") || return
    __architecture__=$(yaml::get "${__file__}" "${__expression__}.architecture") || return
    __size__=$(yaml::get "${__file__}" "${__expression__}.size") || return
    __release__=$(yaml::get "${__file__}" "${__expression__}.system.release") || return
    __kernel__=$(yaml::get "${__file__}" "${__expression__}.kernel") || return
    [[ -z "${__requested_version__}" || "${__requested_version__}" == "${__version__}" ]] || {
        log::error "Configured image version is ${__version__}, not ${__requested_version__}"
        return 1
    }

    local __cache__
    local __directory__
    local __rootfs__
    __cache__=$(module::var image cache) || return
    __directory__="${__cache__}/${__name__}/${__version__}"
    __rootfs__="${__directory__}/rootfs.ext4"
    if [[ -f "${__rootfs__}" && "${__force__}" != true && "${__dry_run__}" != true ]]; then
        log::error "Image already exists: ${__name__}:${__version__}"
        return 1
    fi

    local -a __packages__=("${IMAGE_BASE_PACKAGES[@]}")
    local __package__
    while IFS= read -r __package__; do
        [[ -n "${__package__}" ]] && __packages__+=("${__package__}")
    done < <(yq eval -r "${__expression__}.packages[]?" "${__file__}")
    mapfile -t __packages__ < <(printf '%s\n' "${__packages__[@]}" | sort -u)

    if [[ "${__dry_run__}" == true ]]; then
        printf 'Image: %s:%s\nDistribution: debian %s\nArchitecture: %s\nSize: %s\nKernel: %s\nPackages:\n' \
            "${__name__}" "${__version__}" "${__release__}" "${__architecture__}" "${__size__}" "${__kernel__}"
        printf '  - %s\n' "${__packages__[@]}"
        printf 'Destination: %s\n' "${__rootfs__}"
        return 0
    fi

    user::require_root || return
    proc::require debootstrap truncate mkfs.ext4 e2fsck || return
    module::use kernel || return
    local __kernel_path__
    __kernel_path__=$(kernel::resolve "${__kernel__}" --file "${__file__}") || return

    local __debian_architecture__
    case "${__architecture__}" in
        x86_64) __debian_architecture__=amd64 ;;
        aarch64) __debian_architecture__=arm64 ;;
        *)
            log::error "Unsupported image architecture: ${__architecture__}"
            return 2
            ;;
    esac

    local __staging__="${__directory__}.building"
    local __root__="${__staging__}/root"
    fs::remove_tree "${__staging__}" || return
    fs::mkdir "${__root__}" || return

    local __package_list__
    __package_list__=$(
        IFS=,
        printf '%s' "${__packages__[*]}"
    )
    log::info "Bootstrapping Debian ${__release__}: ${__name__}:${__version__}"
    debootstrap \
        --arch="${__debian_architecture__}" \
        --variant=minbase \
        --include="${__package_list__}" \
        "${__release__}" \
        "${__root__}" \
        https://deb.debian.org/debian || return

    image::_configure_root "${__file__}" "${__expression__}" "${__name__}" "${__root__}" || return
    truncate -s "${__size__}" "${__staging__}/rootfs.ext4" || return
    mkfs.ext4 -q -F -d "${__root__}" "${__staging__}/rootfs.ext4" || return
    e2fsck -pf "${__staging__}/rootfs.ext4" || {
        local __fsck_status__=$?
        ((__fsck_status__ == 1)) || return "${__fsck_status__}"
    }

    local __checksum__
    __checksum__=$(fs::checksum "${__staging__}/rootfs.ext4") || return
    fs::atomic_write "${__staging__}/metadata.yaml" <<EOF
name: ${__name__}
version: ${__version__}
architecture: ${__architecture__}
distribution: debian
release: ${__release__}
kernel: ${__kernel__}
kernel_path: ${__kernel_path__}
rootfs_path: ${__rootfs__}
sha256: ${__checksum__}
EOF
    fs::append "${__staging__}/metadata.yaml" <<EOF
capabilities:
EOF
    local __capability__
    while IFS= read -r __capability__; do
        printf '  - %s\n' "${__capability__}" | fs::append "${__staging__}/metadata.yaml"
    done < <(yq eval -r "${__expression__}.capabilities[]?" "${__file__}")
    fs::remove_tree "${__root__}" || return
    fs::mkdir_parent "${__directory__}" || return
    [[ ! -e "${__directory__}" ]] || fs::remove_tree "${__directory__}" || return
    fs::move "${__staging__}" "${__directory__}" || return
    if [[ -n "${SUDO_UID:-}" && -n "${SUDO_GID:-}" ]]; then
        chown -R "${SUDO_UID}:${SUDO_GID}" -- "${__directory__}" || return
    fi
    log::info "Built image: ${__name__}:${__version__}" "path=${__rootfs__}"
}
