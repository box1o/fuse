#!/usr/bin/env bash

firecracker::status()
{
    if (($# != 1)); then
        log::error "Usage: rune vm status <name>"
        return 2
    fi

    local __name__=$1
    firecracker::_require_vm "${__name__}" || return

    local __state__=stopped
    local __pid__=
    if firecracker::_is_running "${__name__}"; then
        __state__=running
        __pid__=$(firecracker::_pid "${__name__}") || return
    fi
    if [[ $(runtime::config::get 'output|format' text) == json ]]; then
        printf '{"name":"%s","status":"%s","pid":%s}\n' \
            "$(firecracker::_json_string "${__name__}")" \
            "${__state__}" \
            "${__pid__:-null}"
        return
    fi
    [[ -n "${__pid__}" ]] && printf '%s %s pid=%s\n' "${__name__}" "${__state__}" "${__pid__}" ||
        printf '%s %s\n' "${__name__}" "${__state__}"
}

firecracker::list()
{
    if (($# != 0)); then
        log::error "Usage: rune vm list"
        return 2
    fi

    local __vms__
    __vms__=$(module::var vm vms) || return
    local __json__=false
    [[ $(runtime::config::get 'output|format' text) == json ]] && __json__=true
    [[ "${__json__}" == true ]] && printf '[' || printf '%-24s %-10s %s\n' NAME STATUS PID
    if [[ ! -d "${__vms__}" ]]; then
        [[ "${__json__}" == true ]] && printf ']\n'
        return 0
    fi

    local __vm__
    local __name__
    local __pid__
    local __separator__=
    while IFS= read -r -d '' __vm__; do
        __name__=${__vm__##*/}
        if firecracker::_is_running "${__name__}"; then
            __pid__=$(firecracker::_pid "${__name__}") || return
            if [[ "${__json__}" == true ]]; then
                printf '%s{"name":"%s","status":"running","pid":%s}' \
                    "${__separator__}" "$(firecracker::_json_string "${__name__}")" "${__pid__}"
            else
                printf '%-24s %-10s %s\n' "${__name__}" running "${__pid__}"
            fi
        else
            if [[ "${__json__}" == true ]]; then
                printf '%s{"name":"%s","status":"stopped","pid":null}' \
                    "${__separator__}" "$(firecracker::_json_string "${__name__}")"
            else
                printf '%-24s %-10s %s\n' "${__name__}" stopped -
            fi
        fi
        __separator__=,
    done < <(find "${__vms__}" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)
    if [[ "${__json__}" == true ]]; then
        printf ']\n'
    fi
}
