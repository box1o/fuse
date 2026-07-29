#!/usr/bin/env bash

image::ensure()
{
    local __reference__=${1:-}
    shift || true
    local __name__
    local __version__
    image::_split_reference "${__reference__}" __name__ __version__ || return

    local __file__=
    local __argument__
    local -a __forward__=()
    while (($# > 0)); do
        __argument__=$1
        __forward__+=("$1")
        if [[ "$1" == --file ]]; then
            [[ $# -ge 2 ]] || return 2
            __file__=$2
            __forward__+=("$2")
            shift 2
        else
            shift
        fi
    done
    image::_load_config __file__ "${__file__}" || return
    local __expression__
    __expression__=$(image::_expression "${__name__}") || return
    [[ -n "${__version__}" ]] || __version__=$(yaml::get "${__file__}" "${__expression__}.version") || return

    local __cache__
    __cache__=$(module::var image cache) || return
    if [[ ! -f "${__cache__}/${__name__}/${__version__}/rootfs.ext4" ]]; then
        image::build "${__name__}:${__version__}" "${__forward__[@]}" || return
    fi
    printf '%s:%s\n' "${__name__}" "${__version__}"
}

image::_metadata()
{
    local __reference__=${1:-}
    local __field__=${2:-}
    local __name__
    local __version__
    image::_split_reference "${__reference__}" __name__ __version__ || return
    local __cache__
    __cache__=$(module::var image cache) || return
    [[ -n "${__version__}" ]] || {
        log::error "Image version is required: ${__reference__}"
        return 2
    }
    local __metadata__="${__cache__}/${__name__}/${__version__}/metadata.yaml"
    fs::require_file image-metadata "${__metadata__}" || return
    yaml::get "${__metadata__}" ".${__field__}"
}

image::resolve()
{
    local __reference__=${1:-}
    local __field__=${2:-rootfs_path}
    image::_metadata "${__reference__}" "${__field__}"
}

image::path()
{
    (($# == 1)) || return 2
    image::resolve "$1" rootfs_path
}

image::list()
{
    (($# == 0)) || return 2
    local __cache__
    __cache__=$(module::var image cache) || return
    [[ -d "${__cache__}" ]] || return 0
    local __metadata__
    while IFS= read -r -d '' __metadata__; do
        printf '%s:%s\n' \
            "$(yaml::get "${__metadata__}" '.name')" \
            "$(yaml::get "${__metadata__}" '.version')"
    done < <(find "${__cache__}" -mindepth 3 -maxdepth 3 -type f -name metadata.yaml -print0 | sort -z)
}

image::show()
{
    (($# == 1)) || return 2
    local __name__
    local __version__
    image::_split_reference "$1" __name__ __version__ || return
    local __cache__
    __cache__=$(module::var image cache) || return
    fs::require_file image-metadata "${__cache__}/${__name__}/${__version__}/metadata.yaml" || return
    cat "${__cache__}/${__name__}/${__version__}/metadata.yaml"
}

image::delete()
{
    (($# == 1)) || return 2
    local __name__
    local __version__
    image::_split_reference "$1" __name__ __version__ || return
    [[ -n "${__version__}" ]] || return 2
    local __cache__
    local __target__
    __cache__=$(module::var image cache) || return
    __target__="${__cache__}/${__name__}/${__version__}"
    path::assert_within "${__target__}" "${__cache__}" || return
    fs::remove_tree "${__target__}"
}
