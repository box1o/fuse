#!/usr/bin/env bash

volume::_valid_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid volume name: ${1:-}"
        return 2
    fi
}

volume::_directory()
{
    local __name__=${1:-}
    volume::_valid_name "${__name__}" || return
    printf '%s/%s\n' "$(module::var volume cache)" "${__name__}"
}

volume::path()
{
    if (($# != 1)); then
        log::error "Usage: rune volume path <name>"
        return 2
    fi

    local __path__
    __path__="$(volume::_directory "$1")/volume.ext4" || return
    fs::require_file volume "${__path__}" || return
    printf '%s\n' "${__path__}"
}

volume::_metadata()
{
    printf '%s/metadata.yaml\n' "$(volume::_directory "${1:-}")"
}

volume::_write_metadata()
{
    local __name__=${1:-}
    local __source__=${2:-}
    local __owner__=${3:-}
    local __path__=${4:-}
    local __metadata__=${5:-}
    local __size__
    local __document__
    __size__=$(stat -c '%s' -- "${__path__}") || return

    __document__=$(VOLUME_NAME=${__name__} \
        VOLUME_SOURCE=${__source__} \
        VOLUME_OWNER=${__owner__} \
        VOLUME_PATH=${__path__} \
        VOLUME_SIZE=${__size__} \
        VOLUME_CREATED=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
        yq -n -o=yaml '
            .name = strenv(VOLUME_NAME) |
            .filesystem = "ext4" |
            .source = strenv(VOLUME_SOURCE) |
            .owner = strenv(VOLUME_OWNER) |
            .path = strenv(VOLUME_PATH) |
            .size_bytes = (strenv(VOLUME_SIZE) | tonumber) |
            .created_at = strenv(VOLUME_CREATED)
        ') || return
    printf '%s\n' "${__document__}" | fs::atomic_write "${__metadata__}"
}

volume::create()
{
    args::reset
    args::option --name from --required
    args::option --name owner --default ''
    args::parse "$@" || return

    if (($(args::count) != 1)); then
        log::error "Usage: rune volume create <name> --from PATH [--owner ID]"
        return 2
    fi

    local __name__
    local __source__
    local __owner__
    __name__=$(args::position 0) || return
    __source__=$(args::get from) || return
    __owner__=$(args::get owner '') || return
    volume::_valid_name "${__name__}" || return
    __source__=$(path::absolute "${__source__}") || return
    fs::require_file source-volume "${__source__}" || return

    local __lock__
    __lock__=$(runtime::lock_file volume "${__name__}") || return
    lock::with "${__lock__}" volume::_create "${__name__}" "${__source__}" "${__owner__}"
}

volume::_create()
{
    local __name__=${1:-}
    local __source__=${2:-}
    local __owner__=${3:-}
    local __directory__
    local __staging__
    __directory__=$(volume::_directory "${__name__}") || return
    __staging__="${__directory__}.building"

    if [[ -e "${__directory__}" || -e "${__staging__}" ]]; then
        log::error "Volume already exists: ${__name__}"
        return 1
    fi

    fs::mkdir "${__staging__}" || return
    if ! cp --reflink=auto --sparse=always -- "${__source__}" "${__staging__}/volume.ext4"; then
        fs::remove_tree "${__staging__}"
        return 1
    fi
    volume::_write_metadata \
        "${__name__}" "${__source__}" "${__owner__}" \
        "${__staging__}/volume.ext4" "${__staging__}/metadata.yaml" || {
        fs::remove_tree "${__staging__}"
        return 1
    }
    fs::move "${__staging__}" "${__directory__}" || return
    yaml::set "${__directory__}/metadata.yaml" '.path' "${__directory__}/volume.ext4" || return
    printf '%s\n' "${__directory__}/volume.ext4"
}

volume::clone()
{
    args::reset
    args::option --name owner --default ''
    args::parse "$@" || return

    if (($(args::count) != 2)); then
        log::error "Usage: rune volume clone <source> <destination> [--owner ID]"
        return 2
    fi

    local __source__
    local __destination__
    __source__=$(args::position 0) || return
    __destination__=$(args::position 1) || return
    volume::create "${__destination__}" \
        --from "$(volume::path "${__source__}")" \
        --owner "$(args::get owner '')"
}

volume::show()
{
    if (($# != 1)); then
        log::error "Usage: rune volume show <name>"
        return 2
    fi
    local __metadata__
    __metadata__=$(volume::_metadata "$1") || return
    fs::require_file volume-metadata "${__metadata__}" || return
    cat -- "${__metadata__}"
}

volume::list()
{
    if (($# != 0)); then
        log::error "Usage: rune volume list"
        return 2
    fi

    local __cache__
    __cache__=$(module::var volume cache) || return
    [[ -d "${__cache__}" ]] || return 0
    find "${__cache__}" -mindepth 2 -maxdepth 2 -name metadata.yaml -print0 |
        sort -z |
        while IFS= read -r -d '' __metadata__; do
            printf '%s\t%s\t%s\n' \
                "$(yaml::get "${__metadata__}" '.name')" \
                "$(yaml::get "${__metadata__}" '.size_bytes')" \
                "$(yaml::get "${__metadata__}" '.owner' '')"
        done
}

volume::mounted()
{
    local __path__
    __path__=$(volume::path "${1:-}") || return
    local __loop__
    while IFS= read -r __loop__; do
        [[ -n "${__loop__}" ]] || continue
        findmnt -rn -S "${__loop__}" >/dev/null && return 0
    done < <(losetup -j "${__path__}" 2>/dev/null | cut -d: -f1)
    return 1
}
