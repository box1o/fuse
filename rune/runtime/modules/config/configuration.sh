#!/usr/bin/env bash

project_config::_load()
{
    local -n _file_=$1
    local __input__=${2:-}
    [[ -n "${__input__}" ]] || __input__=$(paths::require 'runtime|config') || return
    config::load rune "${__input__}" 1 || return
    _file_=$(config::file rune) || return
}

project_config::validate()
{
    args::reset
    args::option --name file
    args::parse "$@" || return
    if (($(args::count) != 0)); then
        log::error "Usage: rune config validate [--file PATH]"
        return 2
    fi

    local __file__
    project_config::_load __file__ "$(args::get file '')" || return
    project_config::_require_mapping "${__file__}" '.kernels' || return
    project_config::_require_mapping "${__file__}" '.images' || return
    kernel::validate --file "${__file__}" || return
    image::validate --file "${__file__}" || return
    log::info "Configuration is valid: ${__file__}"
}

project_config::_require_mapping()
{
    local __file__=${1:-}
    local __expression__=${2:-}
    local __type__
    __type__=$(yaml::type "${__file__}" "${__expression__}") || return
    if [[ "${__type__}" != '!!map' ]]; then
        log::error "Configuration value must be a mapping: ${__expression__}"
        return 2
    fi
}

project_config::get()
{
    args::reset
    args::option --name default
    args::option --name file
    args::parse "$@" || return
    if (($(args::count) != 1)); then
        log::error "Usage: rune config get <expression> [--default VALUE] [--file PATH]"
        return 2
    fi

    local __file__
    local __expression__
    project_config::_load __file__ "$(args::get file '')" || return
    __expression__=$(args::position 0) || return
    if args::has default; then
        yaml::get "${__file__}" "${__expression__}" "$(args::get default)"
    else
        yaml::get "${__file__}" "${__expression__}"
    fi
}

project_config::show()
{
    args::reset
    args::option --name file
    args::parse "$@" || return
    if (($(args::count) != 0)); then
        log::error "Usage: rune config show [--file PATH]"
        return 2
    fi
    local __file__
    project_config::_load __file__ "$(args::get file '')" || return
    cat -- "${__file__}"
}

project_config::path()
{
    if (($# != 0)); then
        log::error "Usage: rune config path"
        return 2
    fi
    paths::require 'runtime|config'
}
