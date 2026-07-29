#!/usr/bin/env bash

readonly -a IMAGE_BASE_PACKAGES=(
    busybox-static
    systemd-sysv
    udev
    ca-certificates
    iproute2
    iputils-ping
    curl
)

image::_valid_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid image name: ${1:-}"
        return 2
    fi
}

image::_load_config()
{
    local -n _file_=$1
    local __input__=${2:-}
    [[ -n "${__input__}" ]] || __input__=$(paths::require 'runtime|config') || return
    config::load rune "${__input__}" 1 || return
    _file_=$(config::file rune) || return
}

image::_expression()
{
    image::_valid_name "${1:-}" || return
    printf '.images."%s"' "$1"
}

image::_split_reference()
{
    local __reference__=${1:-}
    local -n _name_=$2
    local -n _version_=$3

    _name_=${__reference__%%:*}
    if [[ "${__reference__}" == *:* ]]; then
        _version_=${__reference__#*:}
    else
        _version_=
    fi
    image::_valid_name "${_name_}"
}

image::validate()
{
    local __file__=
    local __name__=
    while (($# > 0)); do
        case "$1" in
            --file)
                [[ $# -ge 2 ]] || return 2
                __file__=$2
                shift 2
                ;;
            -*)
                log::error "image validate: unknown option: $1"
                return 2
                ;;
            *)
                [[ -z "${__name__}" ]] || return 2
                __name__=$1
                shift
                ;;
        esac
    done

    image::_load_config __file__ "${__file__}" || return
    if [[ -n "${__name__}" ]]; then
        image::_validate_definition "${__file__}" "${__name__}"
        return
    fi

    local __image__
    while IFS= read -r __image__; do
        image::_validate_definition "${__file__}" "${__image__}" || return
    done < <(config::keys rune '.images')
}

image::_validate_definition()
{
    local __file__=${1:-}
    local __name__=${2:-}
    local __expression__
    __expression__=$(image::_expression "${__name__}") || return
    yaml::require "${__file__}" \
        "${__expression__}.version" \
        "${__expression__}.architecture" \
        "${__expression__}.size" \
        "${__expression__}.system.distribution" \
        "${__expression__}.system.release" \
        "${__expression__}.kernel" || return

    [[ $(yaml::get "${__file__}" "${__expression__}.system.distribution") == debian ]] || {
        log::error "Only Debian images are currently supported: ${__name__}"
        return 1
    }
    [[ $(yaml::get "${__file__}" "${__expression__}.size") =~ ^[1-9][0-9]*[MG]$ ]] || {
        log::error "Image size must use M or G units: ${__name__}"
        return 2
    }

    local __architecture__
    local __version__
    local __kernel__
    __architecture__=$(yaml::get "${__file__}" "${__expression__}.architecture") || return
    __version__=$(yaml::get "${__file__}" "${__expression__}.version") || return
    __kernel__=$(yaml::get "${__file__}" "${__expression__}.kernel") || return
    validate::enum "${__architecture__}" x86_64 aarch64 || return
    if [[ ! "${__version__}" =~ ^[a-zA-Z0-9][a-zA-Z0-9._-]*$ ]]; then
        log::error "Invalid image version: ${__name__}:${__version__}"
        return 2
    fi
    if ! yaml::has "${__file__}" ".kernels.\"${__kernel__}\""; then
        log::error "Image ${__name__} references unknown kernel: ${__kernel__}"
        return 1
    fi

    local __package__
    while IFS= read -r __package__; do
        [[ "${__package__}" =~ ^[a-zA-Z0-9][a-zA-Z0-9+.-]*$ ]] || {
            log::error "Invalid Debian package name: ${__package__}"
            return 2
        }
    done < <(yq eval -r "${__expression__}.packages[]?" "${__file__}")

    local __environment_key__
    while IFS= read -r __environment_key__; do
        [[ "${__environment_key__}" =~ ^[A-Z_][A-Z0-9_]*$ ]] || {
            log::error "Invalid image environment variable: ${__environment_key__}"
            return 2
        }
    done < <(yq eval -r "${__expression__}.environment // {} | keys | .[]" "${__file__}")

    local __file_count__
    local __index__
    local __destination__
    local __mode__
    __file_count__=$(yaml::length "${__file__}" "${__expression__}.files // []") || return
    for ((__index__ = 0; __index__ < __file_count__; __index__++)); do
        yaml::require "${__file__}" \
            "${__expression__}.files[${__index__}].source" \
            "${__expression__}.files[${__index__}].destination" || return
        __destination__=$(yaml::get "${__file__}" "${__expression__}.files[${__index__}].destination") || return
        __mode__=$(yaml::get "${__file__}" "${__expression__}.files[${__index__}].mode" '') || return
        [[ "${__destination__}" == /* ]] || {
            log::error "Unsafe image file destination: ${__destination__}"
            return 2
        }
        case "${__destination__}" in
            / | */../* | */.. | */./*)
                log::error "Unsafe image file destination: ${__destination__}"
                return 2
                ;;
        esac
        [[ -z "${__mode__}" || "${__mode__}" =~ ^0?[0-7]{3,4}$ ]] || {
            log::error "Invalid image file mode: ${__mode__}"
            return 2
        }
    done
}
