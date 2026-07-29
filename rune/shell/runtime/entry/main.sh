#!/usr/bin/env bash

runtime::_load_modules()
{
    local __override__=${1:-}
    local __default__
    local __entry__
    local -a __roots__=()
    __default__=$(paths::require 'runtime|modules') || return
    if [[ -n "${__override__}" ]]; then
        __roots__+=("${__override__}")
    fi
    __roots__+=("${__default__}")

    if [[ -n "${RUNE_MODULE_PATH:-}" ]]; then
        local -a __extra__=()
        IFS=: read -r -a __extra__ <<<"${RUNE_MODULE_PATH}"
        for __entry__ in "${__extra__[@]}"; do
            if [[ -n "${__entry__}" ]]; then
                __roots__+=("${__entry__}")
            fi
        done
    fi

    module::load_all "${__roots__[@]}"
}

runtime::main()
{
    local __modules_dir__=
    local __want_help__=false
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --help | -h)
                __want_help__=true
                shift
                break
                ;;
            --version)
                runtime::builtin::version
                return
                ;;
            --debug)
                LOG_LEVEL=${LOG_LEVEL_DEBUG}
                shift
                ;;
            --quiet)
                LOG_LEVEL=${LOG_LEVEL_WARN}
                shift
                ;;
            --no-color)
                LOG_COLOR=never
                shift
                ;;
            --json)
                runtime::config::set 'output|format' json
                LOG_COLOR=never
                shift
                ;;
            --modules-dir)
                [[ $# -ge 2 ]] || {
                    log::error '--modules-dir requires a path'
                    return 2
                }
                __modules_dir__=$2
                shift 2
                ;;
            --)
                shift
                break
                ;;
            -*)
                log::error "Unknown global option: $1"
                return 2
                ;;
            *)
                break
                ;;
        esac
    done
    runtime::_load_modules "${__modules_dir__}" || return
    if [[ "${__want_help__}" == true ]]; then
        help::global
        return
    fi
    if (($# == 0)); then
        help::global
        return 2
    fi
    if [[ $1 == help ]]; then
        shift
        case $# in
            0) help::global ;;
            1) help::module "$1" ;;
            2) help::command "$1" "$2" ;;
            *)
                log::error 'Usage: rune help [module [command]]'
                return 2
                ;;
        esac
        return
    fi
    local __module__=$1
    shift
    if [[ $# -eq 0 || $1 == --help || $1 == -h ]]; then
        help::module "${__module__}"
        return
    fi
    local __command__=$1
    shift
    command::dispatch "${__module__}" "${__command__}" "$@"
}
