#!/usr/bin/env bash

volume::delete()
{
    args::reset
    args::flag --name force
    args::parse "$@" || return
    if (($(args::count) != 1)); then
        log::error "Usage: rune volume delete <name> [--force]"
        return 2
    fi

    local __name__
    local __force__=false
    __name__=$(args::position 0) || return
    args::has force && __force__=true
    local __lock__
    __lock__=$(runtime::lock_file volume "${__name__}") || return
    lock::with "${__lock__}" volume::_delete "${__name__}" "${__force__}"
}

volume::_delete()
{
    local __name__=${1:-}
    local __force__=${2:-false}
    local __metadata__
    local __owner__
    local __directory__
    __metadata__=$(volume::_metadata "${__name__}") || return
    fs::require_file volume-metadata "${__metadata__}" || return
    __owner__=$(yaml::get "${__metadata__}" '.owner' '') || return
    if [[ -n "${__owner__}" && "${__force__}" != true ]]; then
        log::error "Volume is attached to ${__owner__}: ${__name__}"
        return 1
    fi
    volume::mounted "${__name__}" && {
        log::error "Volume is mounted: ${__name__}"
        return 1
    }
    __directory__=$(volume::_directory "${__name__}") || return
    path::assert_within "${__directory__}" "$(module::var volume cache)" || return
    fs::remove_tree "${__directory__}" || return
    log::info "Deleted volume: ${__name__}"
}

volume::cleanup()
{
    args::reset
    args::flag --name force
    args::parse "$@" || return
    if (($(args::count) != 0)); then
        log::error "Usage: rune volume cleanup [--force]"
        return 2
    fi

    local __cache__
    local __metadata__
    local __name__
    local __owner__
    local __force__=false
    args::has force && __force__=true
    __cache__=$(module::var volume cache) || return
    [[ -d "${__cache__}" ]] || return 0
    while IFS= read -r -d '' __metadata__; do
        __owner__=$(yaml::get "${__metadata__}" '.owner' '') || return
        [[ -z "${__owner__}" ]] || continue
        __name__=$(yaml::get "${__metadata__}" '.name') || return
        if [[ "${__force__}" == true ]]; then
            volume::delete "${__name__}" --force || return
        else
            printf '%s\n' "${__name__}"
        fi
    done < <(find "${__cache__}" -mindepth 2 -maxdepth 2 -name metadata.yaml -print0 | sort -z)
}
