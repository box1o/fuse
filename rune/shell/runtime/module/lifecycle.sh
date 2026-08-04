#!/usr/bin/env bash

module::initialize()
{
    local __name__=${1:-}

    if ! module::has "${__name__}"; then
        log::error "Unknown module: ${__name__}"
        return 1
    fi

    if [[ "${MODULE_STATE_REGISTRY[${__name__}]}" == initialized ]]; then
        return 0
    fi

    if [[ "${MODULE_STATE_REGISTRY[${__name__}]}" == initializing ]]; then
        log::error "Module dependency cycle detected at: ${__name__}"
        return 1
    fi

    MODULE_STATE_REGISTRY["${__name__}"]=initializing

    local __dependency__
    local __dependencies__=${MODULE_DEPENDENCY_REGISTRY[${__name__}]:-}
    local -a __items__=()
    IFS=',' read -ra __items__ <<<"${__dependencies__}"
    for __dependency__ in "${__items__[@]}"; do
        [[ -n "${__dependency__}" ]] || continue
        if ! module::has "${__dependency__}"; then
            log::error "Module ${__name__} requires missing module: ${__dependency__}"
            MODULE_STATE_REGISTRY["${__name__}"]=failed
            return 1
        fi
        if ! module::initialize "${__dependency__}"; then
            MODULE_STATE_REGISTRY["${__name__}"]=failed
            return 1
        fi
    done

    local __handler__=${MODULE_INIT_HANDLER_REGISTRY[${__name__}]:-}
    local __context__=${MODULE_CONTEXT_REGISTRY[${__name__}]}

    if [[ -n "${__handler__}" ]] && ! "${__handler__}" "${__context__}"; then
        MODULE_STATE_REGISTRY["${__name__}"]=failed
        return 1
    fi

    MODULE_STATE_REGISTRY["${__name__}"]=initialized
}

module::shutdown()
{
    local __name__=${1:-}

    module::has "${__name__}" || return 1

    if [[ "${MODULE_STATE_REGISTRY[${__name__}]}" != initialized ]]; then
        return 0
    fi

    local __handler__=${MODULE_SHUTDOWN_HANDLER_REGISTRY[${__name__}]:-}
    local __context__=${MODULE_CONTEXT_REGISTRY[${__name__}]}

    if [[ -n "${__handler__}" ]]; then
        "${__handler__}" "${__context__}" || return
    fi

    MODULE_STATE_REGISTRY["${__name__}"]=shutdown
}

module::validate()
{
    local __name__=${1:-}

    if ! module::has "${__name__}"; then
        log::error "Unknown module: ${__name__}"
        return 1
    fi

    if [[ "${MODULE_BUILTIN_REGISTRY[${__name__}]}" == true ]]; then
        return 0
    fi

    local __root__=${MODULE_ROOT_REGISTRY[${__name__}]}
    fs::require_file module-declaration "${__root__}/module.sh" || return

    local __file__
    while IFS= read -r -d '' __file__; do
        bash -n "${__file__}" || return
    done < <(find "${__root__}" -type f -name '*.sh' -print0)

    local __command_id__
    while IFS= read -r __command_id__; do
        command::validate "${__command_id__}" || return
    done < <(command::list | awk -v prefix="${__name__} " 'index($0, prefix) == 1')
}
