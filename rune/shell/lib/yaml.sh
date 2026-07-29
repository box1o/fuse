#!/usr/bin/env bash

# YAML operations require Mike Farah yq v4.
yaml::_require_yq()
{
    proc::require yq || return

    local __version__
    __version__=$(yq --version 2>/dev/null) || return

    if [[ "${__version__}" != *'mikefarah/yq'* || "${__version__}" != *'version v4.'* ]]; then
        log::error "yaml::* requires mikefarah/yq v4"
        return 1
    fi
}

yaml::_require_document()
{
    local __file__=${1:-}

    fs::require_file YAML-file "${__file__}" || return
    yaml::_require_yq
}

yaml::_ensure_document()
{
    local __file__=${1:-}

    validate::not_empty file "${__file__}" || return
    yaml::_require_yq || return

    if [[ ! -f "${__file__}" ]]; then
        fs::atomic_write "${__file__}" <<<'{}'
    fi
}

yaml::validate()
{
    local __file__=${1:-}

    yaml::_require_document "${__file__}" || return
    yq eval '.' "${__file__}" >/dev/null
}

yaml::get()
{
    local __file__=${1:-}
    local __expression__=${2:-.}
    local __default__=${3-}

    yaml::_require_document "${__file__}" || return

    local __value__
    __value__=$(
        YAML_DEFAULT_VALUE=${__default__} \
            yq eval -r "${__expression__} // strenv(YAML_DEFAULT_VALUE)" "${__file__}"
    ) || return

    printf '%s\n' "${__value__}"
}

yaml::has()
{
    local __file__=${1:-}
    local __expression__=${2:-}

    validate::not_empty expression "${__expression__}" || return
    yaml::_require_document "${__file__}" || return

    [[ $(yq eval "${__expression__} != null" "${__file__}") == true ]]
}

yaml::type()
{
    local __file__=${1:-}
    local __expression__=${2:-.}

    yaml::_require_document "${__file__}" || return
    yq eval "${__expression__} | type" "${__file__}"
}

yaml::set()
{
    local __file__=${1:-}
    local __expression__=${2:-}
    local __value__=${3-}

    validate::not_empty expression "${__expression__}" || return
    yaml::_ensure_document "${__file__}" || return

    YAML_VALUE=${__value__} \
        yq eval -i "${__expression__} = strenv(YAML_VALUE)" "${__file__}"
}

yaml::set_raw()
{
    local __file__=${1:-}
    local __expression__=${2:-}
    local __value_expression__=${3:-}

    validate::not_empty expression "${__expression__}" || return
    validate::not_empty value-expression "${__value_expression__}" || return
    yaml::_ensure_document "${__file__}" || return

    yq eval -i "${__expression__} = ${__value_expression__}" "${__file__}"
}

yaml::delete()
{
    local __file__=${1:-}
    local __expression__=${2:-}

    validate::not_empty expression "${__expression__}" || return
    yaml::_require_document "${__file__}" || return

    yq eval -i "del(${__expression__})" "${__file__}"
}

yaml::merge()
{
    local __destination__=${1:-}
    shift || true

    validate::not_empty destination "${__destination__}" || return

    if (($# == 0)); then
        log::error "yaml::merge: at least one source file is required"
        return 2
    fi

    yaml::_require_yq || return

    local __source__
    for __source__ in "$@"; do
        fs::require_file YAML-source "${__source__}" || return
    done

    local __temporary__
    __temporary__=$(temp::file yaml-merge) || return

    yq eval-all '. as $item ireduce ({}; . * $item)' "$@" >"${__temporary__}" || return
    fs::mkdir_parent "${__destination__}" || return
    fs::move "${__temporary__}" "${__destination__}"
}

yaml::keys()
{
    local __file__=${1:-}
    local __expression__=${2:-.}

    yaml::_require_document "${__file__}" || return
    yq eval -r "${__expression__} | keys | .[]" "${__file__}"
}

yaml::length()
{
    local __file__=${1:-}
    local __expression__=${2:-.}

    yaml::_require_document "${__file__}" || return
    yq eval "${__expression__} | length" "${__file__}"
}

yaml::to_json()
{
    local __file__=${1:-}
    local __expression__=${2:-.}

    yaml::_require_document "${__file__}" || return
    yq eval -o=json "${__expression__}" "${__file__}"
}

yaml::from_json()
{
    local __json_file__=${1:-}
    local __yaml_file__=${2:-}

    fs::require_file JSON-file "${__json_file__}" || return
    validate::not_empty YAML-destination "${__yaml_file__}" || return
    yaml::_require_yq || return
    fs::mkdir_parent "${__yaml_file__}" || return

    yq eval -P -o=yaml '.' "${__json_file__}" >"${__yaml_file__}"
}

yaml::require()
{
    local __file__=${1:-}
    shift || true

    yaml::validate "${__file__}" || return

    local __expression__
    for __expression__ in "$@"; do
        if ! yaml::has "${__file__}" "${__expression__}"; then
            log::error "Required YAML value missing: ${__expression__}"
            return 1
        fi
    done
}
