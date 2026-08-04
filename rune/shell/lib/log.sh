#!/usr/bin/env bash

readonly LOG_LEVEL_DEBUG=10
readonly LOG_LEVEL_INFO=20
readonly LOG_LEVEL_WARN=30
readonly LOG_LEVEL_ERROR=40

LOG_LEVEL=${LOG_LEVEL:-${LOG_LEVEL_INFO}}
LOG_COLOR=${LOG_COLOR:-auto}

log::_supports_color()
{
    case "${LOG_COLOR}" in
        always) return 0 ;;
        never) return 1 ;;
        auto) [[ -t 2 ]] ;;
        *) return 1 ;;
    esac
}

log::_color_for()
{
    local __level__=${1:-}

    case "${__level__}" in
        DEBUG) printf '\033[36m' ;;
        INFO) printf '\033[32m' ;;
        WARN) printf '\033[33m' ;;
        ERROR) printf '\033[31m' ;;
    esac
}

log::_write()
{
    local __level_name__=${1:-}
    local __level_value__=${2:-0}
    shift 2 || return 2

    if ((__level_value__ < LOG_LEVEL)); then
        return 0
    fi

    local __timestamp__
    __timestamp__=$(date '+%Y-%m-%d %H:%M:%S') || return

    if log::_supports_color; then
        local __color__
        __color__=$(log::_color_for "${__level_name__}") || return

        printf '%s[%s] %-5s\033[0m %s\n' \
            "${__color__}" \
            "${__timestamp__}" \
            "${__level_name__}" \
            "$*" >&2
        return
    fi

    printf '[%s] %-5s %s\n' \
        "${__timestamp__}" \
        "${__level_name__}" \
        "$*" >&2
}

log::debug()
{
    log::_write DEBUG "${LOG_LEVEL_DEBUG}" "$@"
}

log::info()
{
    log::_write INFO "${LOG_LEVEL_INFO}" "$@"
}

log::warn()
{
    log::_write WARN "${LOG_LEVEL_WARN}" "$@"
}

log::error()
{
    log::_write ERROR "${LOG_LEVEL_ERROR}" "$@"
}

log::fatal()
{
    local __exit_code__=1

    if [[ ${1:-} == --code ]]; then
        __exit_code__=${2:-1}
        shift 2
    fi

    log::error "$@"
    exit "${__exit_code__}"
}
