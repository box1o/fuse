#!/usr/bin/env bash

validate::_error()
{
    log::error "$1"
    return 2
}

validate::not_empty()
{
    local __name__=${1:-value}
    local __value__=${2:-}

    [[ -n "${__value__}" ]] || validate::_error "${__name__} must not be empty"
}

validate::argument_count()
{
    local __expected__=${1:-0}
    local __actual__=${2:-0}
    local __function__=${3:-function}

    if ((__actual__ != __expected__)); then
        validate::_error "${__function__}: expected ${__expected__} arguments, got ${__actual__}"
    fi
}

validate::integer()
{
    local __name__=${1:-value}
    local __value__=${2:-}

    [[ "${__value__}" =~ ^-?[0-9]+$ ]] ||
        validate::_error "${__name__} must be an integer: ${__value__}"
}

validate::non_negative_integer()
{
    local __name__=${1:-value}
    local __value__=${2:-}

    [[ "${__value__}" =~ ^[0-9]+$ ]] ||
        validate::_error "${__name__} must be a non-negative integer: ${__value__}"
}

validate::positive_integer()
{
    local __name__=${1:-value}
    local __value__=${2:-}

    [[ "${__value__}" =~ ^[1-9][0-9]*$ ]] ||
        validate::_error "${__name__} must be a positive integer: ${__value__}"
}

validate::boolean()
{
    local __name__=${1:-value}
    local __value__=${2:-}

    case "${__value__}" in
        true | false) return 0 ;;
        *) validate::_error "${__name__} must be true or false: ${__value__}" ;;
    esac
}

validate::enum()
{
    local __value__=${1:-}
    shift || true

    local __candidate__
    for __candidate__ in "$@"; do
        if [[ "${__value__}" == "${__candidate__}" ]]; then
            return 0
        fi
    done

    validate::_error "Invalid value '${__value__}'. Expected one of: $*"
}

validate::absolute_path()
{
    local __value__=${1:-}

    [[ "${__value__}" == /* ]] || validate::_error "Path must be absolute: ${__value__}"
}

validate::identifier()
{
    local __value__=${1:-}

    [[ "${__value__}" =~ ^[a-zA-Z][a-zA-Z0-9._\|-]*$ ]] ||
        validate::_error "Invalid identifier: ${__value__}"
}
