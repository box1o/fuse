#!/usr/bin/env bash

path::join()
{
    if (($# == 0)); then
        log::error "path::join: at least one path is required"
        return 2
    fi

    local __result__=$1
    shift

    local __part__
    for __part__ in "$@"; do
        __result__="${__result__%/}/${__part__#/}"
    done

    printf '%s\n' "${__result__}"
}

path::absolute()
{
    local __path__=${1:-}

    validate::not_empty path "${__path__}" || return
    proc::require realpath || return

    realpath -m -- "${__path__}"
}

path::normalize()
{
    path::absolute "$@"
}

path::basename()
{
    basename -- "${1:-}"
}

path::dirname()
{
    dirname -- "${1:-}"
}

path::extension()
{
    local __name__
    __name__=$(path::basename "${1:-}") || return

    if [[ "${__name__}" == *.* && "${__name__}" != .* ]]; then
        printf '%s\n' "${__name__##*.}"
    fi
}

path::without_extension()
{
    local __path__=${1:-}
    local __directory__
    local __name__

    __directory__=$(path::dirname "${__path__}") || return
    __name__=$(path::basename "${__path__}") || return

    if [[ "${__name__}" == *.* && "${__name__}" != .* ]]; then
        __name__=${__name__%.*}
    fi

    if [[ "${__directory__}" == . ]]; then
        printf '%s\n' "${__name__}"
    else
        printf '%s/%s\n' "${__directory__}" "${__name__}"
    fi
}

path::relative()
{
    local __target__=${1:-}
    local __base__=${2:-${PWD}}

    validate::not_empty target "${__target__}" || return
    proc::require realpath || return

    realpath --relative-to="${__base__}" -m -- "${__target__}"
}

path::is_absolute()
{
    [[ ${1:-} == /* ]]
}

path::is_within()
{
    local __target__
    local __root__

    __target__=$(path::absolute "${1:-}") || return
    __root__=$(path::absolute "${2:-}") || return

    [[ "${__target__}" == "${__root__}" || "${__target__}" == "${__root__}/"* ]]
}

path::assert_within()
{
    local __target__=${1:-}
    local __root__=${2:-}

    if path::is_within "${__target__}" "${__root__}"; then
        return 0
    fi

    log::error "Path '${__target__}' is outside allowed root '${__root__}'"
    return 1
}
