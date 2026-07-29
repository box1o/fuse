#!/usr/bin/env bash

kernel::ensure()
{
    local __file__=
    local __name__=
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
                log::error "kernel ensure: unknown option: $1"
                return 2
                ;;
            *)
                [[ -z "${__name__}" ]] || return 2
                __name__=$1
                shift
                ;;
        esac
    done

    kernel::_valid_name "${__name__}" || return
    if [[ "${__dry_run__}" == true ]]; then
        kernel::_ensure "${__name__}" "${__file__}" "${__force__}" true
        return
    fi
    local __lock__
    __lock__=$(runtime::lock_file kernel "${__name__}") || return
    lock::with "${__lock__}" kernel::_ensure \
        "${__name__}" "${__file__}" "${__force__}" "${__dry_run__}"
}

kernel::_ensure()
{
    local __name__=${1:-}
    local __file__=${2:-}
    local __force__=${3:-false}
    local __dry_run__=${4:-false}
    kernel::_load_config __file__ "${__file__}" || return
    kernel::_validate_definition "${__file__}" "${__name__}" || return

    local __cache__
    local __directory__
    local __destination__
    __cache__=$(module::var kernel cache) || return
    __directory__="${__cache__}/${__name__}"
    __destination__="${__directory__}/vmlinux"

    if [[ -f "${__destination__}" && "${__force__}" != true ]]; then
        printf '%s\n' "${__destination__}"
        return 0
    fi

    local __expression__
    local __version__
    local __architecture__
    local __source__
    local __url__
    __expression__=$(kernel::_expression "${__name__}") || return
    __version__=$(yaml::get "${__file__}" "${__expression__}.version") || return
    __architecture__=$(yaml::get "${__file__}" "${__expression__}.architecture") || return
    __source__=$(yaml::get "${__file__}" "${__expression__}.source") || return

    if [[ "${__source__}" == url ]]; then
        __url__=$(yaml::get "${__file__}" "${__expression__}.url") || return
    else
        __url__=$(kernel::_firecracker_ci_url "${__version__}" "${__architecture__}") || return
    fi

    if [[ "${__dry_run__}" == true ]]; then
        printf 'Kernel: %s\nSource: %s\nDestination: %s\n' "${__name__}" "${__url__}" "${__destination__}"
        return 0
    fi

    fs::mkdir "${__directory__}" || return
    local -a __download_args__=()
    [[ "${__force__}" == true ]] && __download_args__+=(--overwrite)
    log::info "Downloading kernel: ${__url__}"
    net::download "${__download_args__[@]}" "${__url__}" "${__destination__}" || return

    local __checksum__
    __checksum__=$(fs::checksum "${__destination__}") || return
    fs::atomic_write "${__directory__}/metadata.yaml" <<EOF
name: ${__name__}
version: "${__version__}"
architecture: ${__architecture__}
source: ${__source__}
url: ${__url__}
path: ${__destination__}
sha256: ${__checksum__}
EOF

    if [[ -n "${SUDO_UID:-}" && -n "${SUDO_GID:-}" ]]; then
        chown -R "${SUDO_UID}:${SUDO_GID}" -- "${__directory__}" || return
    fi

    printf '%s\n' "${__destination__}"
}

kernel::_firecracker_ci_url()
{
    local __version__=${1:-}
    local __architecture__=${2:-}
    local __bucket__=https://s3.amazonaws.com/spec.ccfc.min
    local __roots__
    __roots__=$(proc::capture curl -fsSL "${__bucket__}?list-type=2&prefix=firecracker-ci/&delimiter=/") || return
    local __ci_prefix__
    __ci_prefix__=$(grep -oE '<Prefix>firecracker-ci/[0-9]{8}-[^<]*/' <<<"${__roots__}" | sed 's/<Prefix>//' | sort | tail -n 1)
    validate::not_empty firecracker-ci-prefix "${__ci_prefix__}" || return

    local __prefix__="${__ci_prefix__}${__architecture__}/vmlinux-${__version__}."
    local __listing__
    __listing__=$(proc::capture curl -fsSL "${__bucket__}?list-type=2&prefix=${__prefix__}") || return
    local __key__
    __key__=$(grep -oE '<Key>[^<]*/vmlinux-[0-9]+\.[0-9]+\.[0-9]+</Key>' <<<"${__listing__}" | sed -e 's/<Key>//' -e 's|</Key>||' | sort -V | tail -n 1)
    validate::not_empty kernel-artifact "${__key__}" || return
    printf '%s/%s\n' "${__bucket__}" "${__key__}"
}
