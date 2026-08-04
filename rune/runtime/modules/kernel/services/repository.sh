#!/usr/bin/env bash

kernel::resolve()
{
    local __name__=${1:-}
    shift || true
    kernel::ensure "${__name__}" "$@"
}

kernel::list()
{
    (($# == 0)) || return 2
    local __cache__
    __cache__=$(module::var kernel cache) || return
    [[ -d "${__cache__}" ]] || return 0
    find "${__cache__}" -mindepth 2 -maxdepth 2 -type f -name metadata.yaml -printf '%h\n' |
        sed 's|.*/||' |
        sort
}

kernel::show()
{
    (($# == 1)) || return 2
    local __cache__
    __cache__=$(module::var kernel cache) || return
    fs::require_file kernel-metadata "${__cache__}/$1/metadata.yaml" || return
    cat "${__cache__}/$1/metadata.yaml"
}

kernel::delete()
{
    (($# == 1)) || return 2
    kernel::_valid_name "$1" || return
    local __cache__
    local __target__
    __cache__=$(module::var kernel cache) || return
    __target__="${__cache__}/$1"
    path::assert_within "${__target__}" "${__cache__}" || return
    fs::remove_tree "${__target__}"
}
