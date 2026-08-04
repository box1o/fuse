#!/usr/bin/env bash

readonly RETRY_DEFAULT_ATTEMPTS=3
readonly RETRY_DEFAULT_DELAY=1

retry::run()
{
    local __attempts__=${RETRY_DEFAULT_ATTEMPTS}
    local __delay__=${RETRY_DEFAULT_DELAY}
    local __backoff__=fixed

    while (($# > 0)); do
        case "$1" in
            --attempts)
                retry::_require_option_value "$1" "$#" || return
                __attempts__=$2
                shift 2
                ;;
            --delay)
                retry::_require_option_value "$1" "$#" || return
                __delay__=$2
                shift 2
                ;;
            --backoff)
                retry::_require_option_value "$1" "$#" || return
                __backoff__=$2
                shift 2
                ;;
            --)
                shift
                break
                ;;
            -*)
                log::error "retry::run: unknown option: $1"
                return 2
                ;;
            *) break ;;
        esac
    done

    validate::positive_integer attempts "${__attempts__}" || return
    validate::positive_integer delay "${__delay__}" || return
    validate::enum "${__backoff__}" fixed exponential || return
    proc::_require_command retry::run "$#" || return

    retry::_attempt "${__attempts__}" "${__delay__}" "${__backoff__}" "$@"
}

retry::_require_option_value()
{
    local __option__=${1:-}
    local __argument_count__=${2:-0}

    if ((__argument_count__ < 2)); then
        log::error "retry::run: ${__option__} requires a value"
        return 2
    fi
}

retry::_attempt()
{
    local __attempts__=$1
    local __sleep__=$2
    local __backoff__=$3
    shift 3

    local __attempt__
    local __status__=0

    for ((__attempt__ = 1; __attempt__ <= __attempts__; __attempt__++)); do
        "$@" && return 0
        __status__=$?

        if ((__attempt__ == __attempts__)); then
            break
        fi

        log::warn "Attempt ${__attempt__}/${__attempts__} failed; retrying in ${__sleep__}s"
        sleep "${__sleep__}"

        if [[ "${__backoff__}" == exponential ]]; then
            __sleep__=$((__sleep__ * 2))
        fi
    done

    return "${__status__}"
}
