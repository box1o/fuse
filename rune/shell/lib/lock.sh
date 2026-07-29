#!/usr/bin/env bash

declare -gA LOCK_FDS=()

lock::acquire()
{
    local __file__=${1:-}
    local __timeout__=${2:-0}

    validate::not_empty lock-file "${__file__}" || return
    validate::non_negative_integer timeout "${__timeout__}" || return
    proc::require flock || return
    fs::mkdir_parent "${__file__}" || return

    local __file_descriptor__
    exec {__file_descriptor__}>"${__file__}" || return

    if ! lock::_wait "${__file_descriptor__}" "${__timeout__}"; then
        exec {__file_descriptor__}>&-
        return 1
    fi

    LOCK_FDS["${__file__}"]=${__file_descriptor__}
}

lock::_wait()
{
    local __file_descriptor__=${1:-}
    local __timeout__=${2:-0}

    if ((__timeout__ == 0)); then
        flock -n "${__file_descriptor__}"
    else
        flock -w "${__timeout__}" "${__file_descriptor__}"
    fi
}

lock::release()
{
    local __file__=${1:-}
    local __file_descriptor__=${LOCK_FDS[${__file__}]:-}

    [[ -n "${__file_descriptor__}" ]] || return 0

    flock -u "${__file_descriptor__}" || true
    exec {__file_descriptor__}>&-
    unset 'LOCK_FDS['"${__file__}"']'
}

lock::with()
{
    local __file__=${1:-}
    shift || true

    proc::_require_command lock::with "$#" || return
    lock::acquire "${__file__}" || return

    local __status__=0
    "$@" || __status__=$?

    lock::release "${__file__}" || {
        if ((__status__ == 0)); then
            return 1
        fi
    }

    return "${__status__}"
}
