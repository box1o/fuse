#!/usr/bin/env bash

proc::format_command()
{
    local __formatted__

    printf -v __formatted__ '%q ' "$@"
    printf '%s' "${__formatted__% }"
}

proc::exists()
{
    local __command__=${1:-}

    [[ -n "${__command__}" ]] && command -v "${__command__}" >/dev/null 2>&1
}

proc::require()
{
    local __command__
    for __command__ in "$@"; do
        if ! proc::exists "${__command__}"; then
            log::error "Required command not found: ${__command__}"
            return 127
        fi
    done
}

proc::_require_command()
{
    local __caller__=${1:-function}
    local __argument_count__=${2:-0}

    if ((__argument_count__ == 0)); then
        log::error "${__caller__}: command is required"
        return 2
    fi
}

proc::run()
{
    proc::_require_command proc::run "$#" || return

    log::debug "Running: $(proc::format_command "$@")"
    "$@"
}

proc::run_quiet()
{
    proc::_require_command proc::run_quiet "$#" || return

    "$@" >/dev/null 2>&1
}

proc::run_allow_fail()
{
    proc::_require_command proc::run_allow_fail "$#" || return

    local __status__=0
    "$@" || __status__=$?

    if ((__status__ != 0)); then
        log::warn "Command failed (${__status__}): $(proc::format_command "$@")"
    fi

    return "${__status__}"
}

proc::capture()
{
    proc::_require_command proc::capture "$#" || return

    local __output__
    local __status__=0

    __output__=$("$@") || __status__=$?

    if ((__status__ != 0)); then
        log::error "Command failed (${__status__}): $(proc::format_command "$@")"
        return "${__status__}"
    fi

    printf '%s\n' "${__output__}"
}

proc::run_as()
{
    local __user__=${1:-}
    shift || true

    validate::not_empty user "${__user__}" || return
    proc::_require_command proc::run_as "$#" || return
    proc::require sudo || return

    sudo -u "${__user__}" -- "$@"
}

proc::timeout()
{
    local __seconds__=${1:-}
    shift || true

    validate::positive_integer seconds "${__seconds__}" || return
    proc::_require_command proc::timeout "$#" || return
    proc::require timeout || return

    timeout --preserve-status "${__seconds__}" "$@"
}
