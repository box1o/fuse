#!/usr/bin/env bash

module::has()
{
    [[ -v "MODULE_ROOT_REGISTRY[${1:-}]" ]]
}

module::get()
{
    local __name__=${1:-}
    local __field__=${2:-}

    if ! module::has "${__name__}"; then
        log::error "Unknown module: ${__name__}"
        return 1
    fi

    case "${__field__}" in
        root) printf '%s\n' "${MODULE_ROOT_REGISTRY[${__name__}]}" ;;
        summary) printf '%s\n' "${MODULE_SUMMARY_REGISTRY[${__name__}]}" ;;
        description) printf '%s\n' "${MODULE_DESCRIPTION_REGISTRY[${__name__}]}" ;;
        version) printf '%s\n' "${MODULE_VERSION_REGISTRY[${__name__}]}" ;;
        runtime-api) printf '%s\n' "${MODULE_API_REGISTRY[${__name__}]}" ;;
        state) printf '%s\n' "${MODULE_STATE_REGISTRY[${__name__}]}" ;;
        dependencies) printf '%s\n' "${MODULE_DEPENDENCY_REGISTRY[${__name__}]}" ;;
        *)
            log::error "Unknown module field: ${__field__}"
            return 2
            ;;
    esac
}

module::root()
{
    module::get "${1:-}" root
}

module::require()
{
    local __name__
    for __name__ in "$@"; do
        if ! module::has "${__name__}"; then
            log::error "Required module is not loaded: ${__name__}"
            return 1
        fi
    done
}

module::use()
{
    local __name__=${1:-}
    module::require "${__name__}" || return
    module::initialize "${__name__}"
}

module::var()
{
    local __name__=${1:-}
    local __key__=${2:-}
    local __default__=${3-}

    if ! module::has "${__name__}"; then
        log::error "Unknown module: ${__name__}"
        return 1
    fi

    local __context__=${MODULE_CONTEXT_REGISTRY[${__name__}]}

    if (($# >= 3)); then
        var::get "${__context__}" "${__key__}" "${__default__}"
    else
        var::get "${__context__}" "${__key__}"
    fi
}

module::list()
{
    printf '%s\n' "${!MODULE_ROOT_REGISTRY[@]}" |
        sed '/^$/d' |
        sort
}
