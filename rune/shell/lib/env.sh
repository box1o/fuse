#!/usr/bin/env bash

env::has()
{
    local __name__=${1:-}

    env::_validate_name "${__name__}" || return
    [[ -v "${__name__}" ]]
}

env::_validate_name()
{
    local __name__=${1:-}

    if [[ ! "${__name__}" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]]; then
        log::error "Invalid environment variable name: ${__name__}"
        return 2
    fi
}

env::get()
{
    local __name__=${1:-}
    local __default__=${2-}

    env::_validate_name "${__name__}" || return
    printf '%s\n' "${!__name__-${__default__}}"
}

env::require()
{
    local __name__
    for __name__ in "$@"; do
        env::_validate_name "${__name__}" || return

        if [[ -z "${!__name__:-}" ]]; then
            log::error "Required environment variable is missing: ${__name__}"
            return 1
        fi
    done
}

env::set_default()
{
    local __name__=${1:-}
    local __value__=${2-}

    env::_validate_name "${__name__}" || return

    if [[ ! -v "${__name__}" ]]; then
        printf -v "${__name__}" '%s' "${__value__}"
    fi

    export "${__name__?}"
}

# Load a trusted shell environment file and export the variables it defines.
env::load_file()
{
    local __file__=${1:-}

    fs::require_file environment-file "${__file__}" || return

    local __restore_allexport__=false
    if [[ $- == *a* ]]; then
        __restore_allexport__=true
    else
        set -a
    fi

    local __status__=0
    # Environment files are executable shell code and must be trusted.
    # shellcheck disable=SC1090
    source "${__file__}" || __status__=$?

    if [[ "${__restore_allexport__}" == false ]]; then
        set +a
    fi

    return "${__status__}"
}

env::log()
{
    local __name__
    for __name__ in "$@"; do
        env::_validate_name "${__name__}" || return
        log::info "${__name__}=${!__name__-}"
    done
}
