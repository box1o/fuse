#!/usr/bin/env bash

declare -gA PATH_REGISTRY=()
declare -gA PATH_DESCRIPTIONS=()
declare -gA PATH_SOURCES=()

paths::has()
{
    local __key__=${1:-}

    validate::identifier "${__key__}" || return
    [[ -v "PATH_REGISTRY[${__key__}]" ]]
}

# Register a normalized path under a unique key.
#
# Usage:
#   paths::register [options] <key> <path>
paths::register()
{
    local __overwrite__=false
    local __description__=
    local __source__=runtime
    local __normalize__=true

    while (($# > 0)); do
        case "$1" in
            --overwrite)
                __overwrite__=true
                shift
                ;;
            --description)
                paths::_require_option_value "$1" "$#" || return
                __description__=$2
                shift 2
                ;;
            --source)
                paths::_require_option_value "$1" "$#" || return
                __source__=$2
                shift 2
                ;;
            --no-normalize)
                __normalize__=false
                shift
                ;;
            --)
                shift
                break
                ;;
            -*)
                log::error "paths::register: unknown option: $1"
                return 2
                ;;
            *) break ;;
        esac
    done

    if (($# != 2)); then
        log::error "Usage: paths::register [options] <key> <path>"
        return 2
    fi

    local __key__=$1
    local __value__=$2

    validate::identifier "${__key__}" || return
    validate::not_empty path "${__value__}" || return

    if paths::has "${__key__}" && [[ "${__overwrite__}" != true ]]; then
        log::error "Path is already registered: ${__key__}"
        return 1
    fi

    if [[ "${__normalize__}" == true ]]; then
        __value__=$(path::absolute "${__value__}") || return
    fi

    PATH_REGISTRY["${__key__}"]=${__value__}
    PATH_DESCRIPTIONS["${__key__}"]=${__description__}
    PATH_SOURCES["${__key__}"]=${__source__}
}

paths::_require_option_value()
{
    local __option__=${1:-}
    local __argument_count__=${2:-0}

    if ((__argument_count__ < 2)); then
        log::error "paths::register: ${__option__} requires a value"
        return 2
    fi
}

paths::set()
{
    paths::register --overwrite "$@"
}

paths::get()
{
    local __key__=${1:-}

    if ! paths::has "${__key__}"; then
        log::error "Unknown path key: ${__key__}"
        return 1
    fi

    printf '%s\n' "${PATH_REGISTRY[${__key__}]}"
}

paths::require()
{
    paths::get "$@"
}

paths::remove()
{
    local __key__=${1:-}

    paths::has "${__key__}" || return 0

    unset \
        'PATH_REGISTRY['"${__key__}"']' \
        'PATH_DESCRIPTIONS['"${__key__}"']' \
        'PATH_SOURCES['"${__key__}"']'
}

paths::keys()
{
    printf '%s\n' "${!PATH_REGISTRY[@]}" | sort
}

paths::print()
{
    printf '%-28s %-12s %-52s %s\n' KEY SOURCE PATH DESCRIPTION
    printf '%-28s %-12s %-52s %s\n' \
        '----------------------------' \
        '------------' \
        '----------------------------------------------------' \
        '------------------------------'

    local __key__
    while IFS= read -r __key__; do
        printf '%-28s %-12s %-52s %s\n' \
            "${__key__}" \
            "${PATH_SOURCES[${__key__}]:-unknown}" \
            "${PATH_REGISTRY[${__key__}]}" \
            "${PATH_DESCRIPTIONS[${__key__}]:-}"
    done < <(paths::keys)
}

paths::log()
{
    log::info "Registered paths"

    local __key__
    while IFS= read -r __key__; do
        log::info \
            "${__key__}=${PATH_REGISTRY[${__key__}]}" \
            "[source=${PATH_SOURCES[${__key__}]:-unknown}]" \
            "${PATH_DESCRIPTIONS[${__key__}]:-}"
    done < <(paths::keys)
}
