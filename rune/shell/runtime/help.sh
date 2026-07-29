#!/usr/bin/env bash

help::_cli_name()
{
    runtime::config::get 'cli|name' rune
}

help::global()
{
    local __cli__
    __cli__=$(help::_cli_name) || return

    printf 'Usage:\n    %s [global-options] <module> <command> [options]\n\n' "${__cli__}"
    printf 'Global options:\n'
    printf '    %-16s %s\n' '--debug' 'Enable diagnostic logging.'
    printf '    %-16s %s\n' '--quiet' 'Show warnings and errors only.'
    printf '    %-16s %s\n' '--no-color' 'Disable colored output.'
    printf '    %-16s %s\n' '--json' 'Use machine-readable output where supported.'
    printf '    %-16s %s\n' '--modules-dir PATH' 'Load an additional module directory.'
    printf '\nModules:\n'

    local __module__
    while IFS= read -r __module__; do
        [[ -n "${__module__}" ]] || continue
        printf '    %-16s %s\n' "${__module__}" "${MODULE_SUMMARY_REGISTRY[${__module__}]}"
    done < <(module::list)

    printf "\nRun '%s help <module>' for module commands.\n" "${__cli__}"
}

help::module()
{
    local __module__=${1:-}

    if ! module::has "${__module__}"; then
        log::error "Unknown module: ${__module__}"
        return 2
    fi

    local __cli__
    __cli__=$(help::_cli_name) || return

    printf 'Usage:\n    %s %s <command> [options]\n\n%s\n\nCommands:\n' \
        "${__cli__}" \
        "${__module__}" \
        "${MODULE_DESCRIPTION_REGISTRY[${__module__}]}"

    local __command__
    local __id__
    while IFS= read -r __command__; do
        [[ -n "${__command__}" ]] || continue
        __id__="${__module__} ${__command__}"
        [[ "${COMMAND_HIDDEN_REGISTRY[${__id__}]:-false}" == true ]] && continue
        printf '    %-16s %s\n' "${__command__}" "${COMMAND_SUMMARY_REGISTRY[${__id__}]}"
    done < <(command::list_module "${__module__}")
}

help::command()
{
    local __module__=${1:-}
    local __command__=${2:-}
    local __id__

    __id__=$(command::_id "${__module__}" "${__command__}") || return

    if ! command::has "${__id__}"; then
        log::error "Unknown command: ${__id__}"
        return 2
    fi

    printf 'Usage:\n    %s\n\n%s\n' \
        "${COMMAND_USAGE_REGISTRY[${__id__}]}" \
        "${COMMAND_DESCRIPTION_REGISTRY[${__id__}]}"
}

help::suggest()
{
    local __module__=${1:-}
    local __command__=${2:-}

    module::has "${__module__}" || return 0

    local __candidate__
    while IFS= read -r __candidate__; do
        if [[ "${__candidate__}" == "${__command__}"* || "${__command__}" == "${__candidate__}"* ]]; then
            log::info "Did you mean: $(help::_cli_name) ${__module__} ${__candidate__}"
            return 0
        fi
    done < <(command::list_module "${__module__}")
}
