#!/usr/bin/env bash

kernel::_valid_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9.-]*$ ]]; then
        log::error "Invalid kernel name: ${1:-}"
        return 2
    fi
}

kernel::_load_config()
{
    local -n _file_=$1
    local __input__=${2:-}
    [[ -n "${__input__}" ]] || __input__=$(paths::require 'runtime|config') || return
    config::load rune "${__input__}" 1 || return
    _file_=$(config::file rune) || return
}

kernel::_expression()
{
    local __name__=${1:-}
    kernel::_valid_name "${__name__}" || return
    printf '.kernels."%s"' "${__name__}"
}

kernel::validate()
{
    local __file__=
    local __name__=

    while (($# > 0)); do
        case "$1" in
            --file)
                [[ $# -ge 2 ]] || {
                    log::error "kernel validate: --file requires a value"
                    return 2
                }
                __file__=$2
                shift 2
                ;;
            -*)
                log::error "kernel validate: unknown option: $1"
                return 2
                ;;
            *)
                [[ -z "${__name__}" ]] || return 2
                __name__=$1
                shift
                ;;
        esac
    done

    kernel::_load_config __file__ "${__file__}" || return
    if [[ -n "${__name__}" ]]; then
        kernel::_validate_definition "${__file__}" "${__name__}"
        return
    fi

    local __kernel__
    while IFS= read -r __kernel__; do
        kernel::_validate_definition "${__file__}" "${__kernel__}" || return
    done < <(config::keys rune '.kernels')
}

kernel::_validate_definition()
{
    local __file__=${1:-}
    local __name__=${2:-}
    local __expression__
    __expression__=$(kernel::_expression "${__name__}") || return
    yaml::require "${__file__}" \
        "${__expression__}.version" \
        "${__expression__}.architecture" \
        "${__expression__}.source" || return

    local __source__
    local __architecture__
    __source__=$(yaml::get "${__file__}" "${__expression__}.source") || return
    __architecture__=$(yaml::get "${__file__}" "${__expression__}.architecture") || return
    validate::enum "${__source__}" url firecracker-ci || return
    validate::enum "${__architecture__}" x86_64 aarch64 || return
    if [[ "${__source__}" == url ]]; then
        yaml::require "${__file__}" "${__expression__}.url" || return
        net::is_url "$(yaml::get "${__file__}" "${__expression__}.url")" || {
            log::error "Invalid kernel URL: ${__name__}"
            return 2
        }
    fi
}
