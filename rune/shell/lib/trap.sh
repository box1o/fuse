#!/usr/bin/env bash

declare -gA TRAP_HANDLERS=()

trap::_dispatch()
{
    local __status__=$?
    local __signal__=${1:-}
    local __handler__

    while IFS= read -r __handler__; do
        [[ -n "${__handler__}" ]] || continue
        "${__handler__}" || true
    done <<<"${TRAP_HANDLERS[${__signal__}]:-}"

    return "${__status__}"
}

# Add a function handler without replacing existing handlers for the signal.
trap::add()
{
    local __signal__=${1:-EXIT}
    local __handler__=${2:-}

    validate::not_empty signal "${__signal__}" || return
    validate::not_empty handler "${__handler__}" || return

    if (($# != 2)); then
        log::error "Usage: trap::add <signal> <handler-function>"
        return 2
    fi

    if ! declare -F "${__handler__}" >/dev/null 2>&1; then
        log::error "Trap handler is not a function: ${__handler__}"
        return 2
    fi

    TRAP_HANDLERS["${__signal__}"]+="${__handler__}"$'\n'

    # The signal name is captured now; handlers are still invoked without eval.
    # shellcheck disable=SC2064
    trap "trap::_dispatch ${__signal__}" "${__signal__}"
}
