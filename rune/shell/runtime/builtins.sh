#!/usr/bin/env bash

readonly RUNTIME_VERSION='0.1.0'
readonly RUNTIME_API_VERSION='1'

runtime::register_builtins()
{
    local __root__
    __root__=$(paths::require 'runtime|shell') || return
    runtime::builtin::_register_module module 'Manage runtime modules' "${__root__}" || return
    runtime::builtin::_register_module paths 'Inspect registered paths' "${__root__}" || return
    runtime::builtin::_register_module runtime 'Inspect the shell runtime' "${__root__}" || return

    runtime::builtin::_register_command module list 'List runtime modules' \
        runtime::builtin::module_list 'rune module list' || return
    runtime::builtin::_register_command module show 'Show module metadata' \
        runtime::builtin::module_show 'rune module show <name>' || return
    runtime::builtin::_register_command module validate 'Validate runtime modules' \
        runtime::builtin::module_validate 'rune module validate [name]' || return
    runtime::builtin::_register_command module scaffold 'Create a new module' \
        scaffold::module 'rune module scaffold <name> [options]' || return
    runtime::builtin::_register_command paths list 'List registered paths' \
        paths::print 'rune paths list' || return
    runtime::builtin::_register_command paths get 'Get a registered path' \
        runtime::builtin::paths_get 'rune paths get <key>' || return
    runtime::builtin::_register_command paths log 'Log registered paths' \
        paths::log 'rune paths log' || return
    runtime::builtin::_register_command runtime version 'Print runtime version' \
        runtime::builtin::version 'rune runtime version' || return
    runtime::builtin::_register_command runtime doctor 'Inspect runtime dependencies' \
        runtime::builtin::doctor 'rune runtime doctor' || return
    runtime::builtin::_register_command runtime commands 'List registered commands' \
        command::list 'rune runtime commands'
}

runtime::builtin::_register_module()
{
    module::register \
        --builtin \
        --name "${1:-}" \
        --summary "${2:-}" \
        --root "${3:-}"
}

runtime::builtin::_register_command()
{
    command::register \
        --module "${1:-}" \
        --name "${2:-}" \
        --summary "${3:-}" \
        --handler "${4:-}" \
        --usage "${5:-}"
}

runtime::builtin::module_list()
{
    module::list
}

runtime::builtin::module_show()
{
    [[ $# -eq 1 ]] || {
        log::error 'Usage: rune module show <name>'
        return 2
    }
    local __name__=$1
    module::has "${__name__}" || {
        log::error "Unknown module: ${__name__}"
        return 1
    }
    printf 'Name: %s\nSummary: %s\nDescription: %s\nVersion: %s\nRuntime API: %s\nRoot: %s\nState: %s\n' \
        "${__name__}" "${MODULE_SUMMARY_REGISTRY[${__name__}]}" "${MODULE_DESCRIPTION_REGISTRY[${__name__}]}" \
        "${MODULE_VERSION_REGISTRY[${__name__}]}" "${MODULE_API_REGISTRY[${__name__}]}" \
        "${MODULE_ROOT_REGISTRY[${__name__}]}" "${MODULE_STATE_REGISTRY[${__name__}]}"
}

runtime::builtin::module_validate()
{
    if [[ $# -gt 1 ]]; then
        log::error 'Usage: rune module validate [name]'
        return 2
    fi
    local __name__
    if [[ $# -eq 1 ]]; then
        module::validate "$1"
        return
    fi
    while IFS= read -r __name__; do
        module::validate "${__name__}" || return
    done < <(module::list)
}

runtime::builtin::paths_get()
{
    [[ $# -eq 1 ]] || {
        log::error 'Usage: rune paths get <key>'
        return 2
    }
    paths::require "$1"
}

runtime::builtin::version()
{
    printf 'rune %s (API %s)\n' "${RUNTIME_VERSION}" "${RUNTIME_API_VERSION}"
}

runtime::builtin::doctor()
{
    local __project__
    local __roots__
    local __command__
    local __status__=0
    local -a __required__=(bash realpath stat sha256sum)
    local -a __optional__=(
        curl
        yq
        tar
        unzip
        flock
        debootstrap
        mkfs.ext4
        e2fsck
        resize2fs
        debugfs
        bats
        shellcheck
        shfmt
    )
    __project__=$(paths::require 'runtime|project_root') || return
    __roots__=$(paths::require 'runtime|modules') || return
    printf 'Runtime:\n  Version: %s\n  API: %s\n  Bash: %s\n  Distribution: %s\n  Architecture: %s\n  Project root: %s\n  Module roots: %s\n  Loaded modules: %s\n  Registered commands: %s\n' \
        "${RUNTIME_VERSION}" "${RUNTIME_API_VERSION}" "${BASH_VERSION}" "$(system::distribution)" "$(system::architecture)" \
        "${__project__}" "${__roots__}" "$(module::list | wc -l)" "$(command::list | wc -l)"
    printf '\nRequired:\n'
    for __command__ in "${__required__[@]}"; do
        if proc::exists "${__command__}"; then
            printf '  [ok] %s\n' "${__command__}"
        else
            printf '  [missing] %s\n' "${__command__}"
            __status__=127
        fi
    done
    printf '\nOptional:\n'
    for __command__ in "${__optional__[@]}"; do
        if proc::exists "${__command__}"; then
            printf '  [ok] %s\n' "${__command__}"
        else
            printf '  [missing] %s\n' "${__command__}"
        fi
    done
    return "${__status__}"
}
