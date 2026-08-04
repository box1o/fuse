#!/usr/bin/env bash

declare -gA COMMAND_HANDLER_REGISTRY=()
declare -gA COMMAND_SUMMARY_REGISTRY=()
declare -gA COMMAND_DESCRIPTION_REGISTRY=()
declare -gA COMMAND_USAGE_REGISTRY=()
declare -gA COMMAND_GROUP_REGISTRY=()
declare -gA COMMAND_HIDDEN_REGISTRY=()

command::_valid_identifier()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid command identifier: ${1:-}"
        return 2
    fi
}

command::_id()
{
    command::_valid_identifier "${1:-}" || return
    command::_valid_identifier "${2:-}" || return

    printf '%s %s\n' "$1" "$2"
}

# Register a command for the module currently being loaded.
#
# A module author normally only needs:
#   command::register --cmd run --handler example::run
command::register()
{
    local __module__=${MODULE_LOADING_NAME:-}
    local __command__=
    local __handler__=
    local __summary__=
    local __description__=
    local __usage__=
    local __group__=general
    local __hidden__=false

    while (($# > 0)); do
        case "$1" in
            --module | --cmd | --name | --handler | --summary | --description | --usage | --group)
                command::_require_option_value "$1" "$#" || return

                case "$1" in
                    --module) __module__=$2 ;;
                    --cmd | --name) __command__=$2 ;;
                    --handler) __handler__=$2 ;;
                    --summary) __summary__=$2 ;;
                    --description) __description__=$2 ;;
                    --usage) __usage__=$2 ;;
                    --group) __group__=$2 ;;
                esac

                shift 2
                ;;
            --hidden)
                __hidden__=true
                shift
                ;;
            *)
                log::error "command::register: unknown option: $1"
                return 2
                ;;
        esac
    done

    validate::not_empty module "${__module__}" || return
    validate::not_empty command "${__command__}" || return
    validate::not_empty handler "${__handler__}" || return

    local __id__
    __id__=$(command::_id "${__module__}" "${__command__}") || return

    if ! module::has "${__module__}"; then
        log::error "Command module is not registered: ${__module__}"
        return 1
    fi

    if command::has "${__id__}"; then
        log::error "Command is already registered: ${__id__}"
        return 1
    fi

    if ! declare -F "${__handler__}" >/dev/null 2>&1; then
        log::error "Command handler does not exist: ${__handler__}"
        return 1
    fi

    [[ -n "${__summary__}" ]] || __summary__=${__description__:-${__command__}}
    [[ -n "${__description__}" ]] || __description__=${__summary__}

    if [[ -z "${__usage__}" ]]; then
        local __cli__
        __cli__=$(runtime::config::get 'cli|name' rune) || return
        __usage__="${__cli__} ${__module__} ${__command__} [options]"
    fi

    COMMAND_HANDLER_REGISTRY["${__id__}"]=${__handler__}
    COMMAND_SUMMARY_REGISTRY["${__id__}"]=${__summary__}
    COMMAND_DESCRIPTION_REGISTRY["${__id__}"]=${__description__}
    COMMAND_USAGE_REGISTRY["${__id__}"]=${__usage__}
    COMMAND_GROUP_REGISTRY["${__id__}"]=${__group__}
    COMMAND_HIDDEN_REGISTRY["${__id__}"]=${__hidden__}
}

command::_require_option_value()
{
    local __option__=${1:-}
    local __argument_count__=${2:-0}

    if ((__argument_count__ < 2)); then
        log::error "command::register: ${__option__} requires a value"
        return 2
    fi
}

command::has()
{
    [[ -v "COMMAND_HANDLER_REGISTRY[${1:-}]" ]]
}

command::get_handler()
{
    local __id__=${1:-}

    if ! command::has "${__id__}"; then
        log::error "Unknown command: ${__id__}"
        return 1
    fi

    printf '%s\n' "${COMMAND_HANDLER_REGISTRY[${__id__}]}"
}

command::list()
{
    printf '%s\n' "${!COMMAND_HANDLER_REGISTRY[@]}" |
        sed '/^$/d' |
        sort
}

command::list_module()
{
    local __module__=${1:-}
    local __id__

    command::_valid_identifier "${__module__}" || return

    while IFS= read -r __id__; do
        if [[ "${__id__}" == "${__module__} "* ]]; then
            printf '%s\n' "${__id__#* }"
        fi
    done < <(command::list)
}

command::validate()
{
    local __id__=${1:-}
    local __handler__

    __handler__=$(command::get_handler "${__id__}") || return

    if ! declare -F "${__handler__}" >/dev/null 2>&1; then
        log::error "Missing command handler: ${__handler__}"
        return 1
    fi
}

command::dispatch()
{
    local __module__=${1:-}
    local __command__=${2:-}

    if (($# < 2)); then
        log::error 'Usage: command::dispatch <module> <command> [arguments]'
        return 2
    fi

    shift 2

    local __id__
    local __handler__
    __id__=$(command::_id "${__module__}" "${__command__}") || return
    __handler__=$(command::get_handler "${__id__}") || {
        help::suggest "${__module__}" "${__command__}"
        return 2
    }

    if [[ ${1:-} == --help || ${1:-} == -h ]]; then
        help::command "${__module__}" "${__command__}"
        return
    fi

    module::initialize "${__module__}" || return

    local __handler_status__=0
    local __shutdown_status__=0
    "${__handler__}" "$@" || __handler_status__=$?
    module::shutdown "${__module__}" || __shutdown_status__=$?

    if ((__handler_status__ != 0)); then
        return "${__handler_status__}"
    fi

    return "${__shutdown_status__}"
}

command::log()
{
    local __id__

    while IFS= read -r __id__; do
        log::info "${__id__} -> ${COMMAND_HANDLER_REGISTRY[${__id__}]}"
    done < <(command::list)
}

command::_remove()
{
    local __id__=${1:-}

    unset \
        'COMMAND_HANDLER_REGISTRY['"${__id__}"']' \
        'COMMAND_SUMMARY_REGISTRY['"${__id__}"']' \
        'COMMAND_DESCRIPTION_REGISTRY['"${__id__}"']' \
        'COMMAND_USAGE_REGISTRY['"${__id__}"']' \
        'COMMAND_GROUP_REGISTRY['"${__id__}"']' \
        'COMMAND_HIDDEN_REGISTRY['"${__id__}"']'
}

command::_reset()
{
    COMMAND_HANDLER_REGISTRY=()
    COMMAND_SUMMARY_REGISTRY=()
    COMMAND_DESCRIPTION_REGISTRY=()
    COMMAND_USAGE_REGISTRY=()
    COMMAND_GROUP_REGISTRY=()
    COMMAND_HIDDEN_REGISTRY=()
}
