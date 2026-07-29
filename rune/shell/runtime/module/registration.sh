#!/usr/bin/env bash

declare -gA MODULE_ROOT_REGISTRY=()
declare -gA MODULE_SUMMARY_REGISTRY=()
declare -gA MODULE_DESCRIPTION_REGISTRY=()
declare -gA MODULE_VERSION_REGISTRY=()
declare -gA MODULE_API_REGISTRY=()
declare -gA MODULE_INIT_HANDLER_REGISTRY=()
declare -gA MODULE_SHUTDOWN_HANDLER_REGISTRY=()
declare -gA MODULE_CONTEXT_REGISTRY=()
declare -gA MODULE_STATE_REGISTRY=()
declare -gA MODULE_BUILTIN_REGISTRY=()
declare -gA MODULE_DEPENDENCY_REGISTRY=()

MODULE_LOADING_NAME=
MODULE_LOADING_ROOT=

module::_valid_name()
{
    if [[ ! ${1:-} =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid module name: ${1:-}"
        return 2
    fi
}

module::_is_reserved()
{
    [[ ${1:-} == module || ${1:-} == paths || ${1:-} == runtime ]]
}

# Register the module currently being loaded.
#
# Required for external modules:
#   module::register --name <name> [--description <text>]
#
# Root, lifecycle handlers, version, API, and summary have useful defaults.
module::register()
{
    local __name__=
    local __summary__=
    local __description__=
    local __version__=0.1.0
    local __root__=
    local __api__=${RUNTIME_API_VERSION:-1}
    local __init__=
    local __shutdown__=
    local __builtin__=false
    local __requires__=

    while (($# > 0)); do
        case "$1" in
            --name | --summary | --description | --version | --root | --runtime-api | --init-handler | --shutdown-handler | --requires)
                module::_require_option_value module::register "$1" "$#" || return

                case "$1" in
                    --name) __name__=$2 ;;
                    --summary) __summary__=$2 ;;
                    --description) __description__=$2 ;;
                    --version) __version__=$2 ;;
                    --root) __root__=$2 ;;
                    --runtime-api) __api__=$2 ;;
                    --init-handler) __init__=$2 ;;
                    --shutdown-handler) __shutdown__=$2 ;;
                    --requires) __requires__=$2 ;;
                esac

                shift 2
                ;;
            --builtin)
                __builtin__=true
                shift
                ;;
            *)
                log::error "module::register: unknown option: $1"
                return 2
                ;;
        esac
    done

    module::_valid_name "${__name__}" || return

    if [[ -z "${__root__}" ]]; then
        __root__=${MODULE_LOADING_ROOT}
    fi

    if [[ -z "${__summary__}" ]]; then
        __summary__=${__description__:-${__name__}}
    fi

    if [[ -z "${__description__}" ]]; then
        __description__=${__summary__}
    fi

    module::_validate_registration \
        "${__name__}" \
        "${__root__}" \
        "${__api__}" \
        "${__builtin__}" || return

    local __namespace__=${__name__//-/_}

    if [[ -z "${__init__}" ]] && declare -F "${__namespace__}::_init" >/dev/null 2>&1; then
        __init__="${__namespace__}::_init"
    fi

    if [[ -z "${__shutdown__}" ]] && declare -F "${__namespace__}::_shutdown" >/dev/null 2>&1; then
        __shutdown__="${__namespace__}::_shutdown"
    fi

    module::_validate_handler init "${__init__}" || return
    module::_validate_handler shutdown "${__shutdown__}" || return
    module::_validate_dependencies "${__name__}" "${__requires__}" || return

    __root__=$(path::absolute "${__root__}") || return

    local __context__="RUNE_MODULE_${__name__^^}_VARS"
    __context__=${__context__//-/_}
    module::_create_context "${__context__}" "${__name__}" "${__root__}" || return

    MODULE_ROOT_REGISTRY["${__name__}"]=${__root__}
    MODULE_SUMMARY_REGISTRY["${__name__}"]=${__summary__}
    MODULE_DESCRIPTION_REGISTRY["${__name__}"]=${__description__}
    MODULE_VERSION_REGISTRY["${__name__}"]=${__version__}
    MODULE_API_REGISTRY["${__name__}"]=${__api__}
    MODULE_INIT_HANDLER_REGISTRY["${__name__}"]=${__init__}
    MODULE_SHUTDOWN_HANDLER_REGISTRY["${__name__}"]=${__shutdown__}
    MODULE_CONTEXT_REGISTRY["${__name__}"]=${__context__}
    MODULE_STATE_REGISTRY["${__name__}"]=loaded
    MODULE_BUILTIN_REGISTRY["${__name__}"]=${__builtin__}
    MODULE_DEPENDENCY_REGISTRY["${__name__}"]=${__requires__}
}

module::_validate_dependencies()
{
    local __module__=${1:-}
    local __dependencies__=${2:-}
    local __dependency__
    local -a __items__=()

    IFS=',' read -ra __items__ <<<"${__dependencies__}"
    for __dependency__ in "${__items__[@]}"; do
        [[ -n "${__dependency__}" ]] || continue
        module::_valid_name "${__dependency__}" || return
        if [[ "${__dependency__}" == "${__module__}" ]]; then
            log::error "Module cannot depend on itself: ${__module__}"
            return 1
        fi
    done
}

module::_require_option_value()
{
    local __caller__=${1:-function}
    local __option__=${2:-}
    local __argument_count__=${3:-0}

    if ((__argument_count__ < 2)); then
        log::error "${__caller__}: ${__option__} requires a value"
        return 2
    fi
}

module::_validate_registration()
{
    local __name__=${1:-}
    local __root__=${2:-}
    local __api__=${3:-}
    local __builtin__=${4:-false}

    validate::not_empty module-root "${__root__}" || return

    if [[ "${__builtin__}" != true ]] && module::_is_reserved "${__name__}"; then
        log::error "Reserved module name: ${__name__}"
        return 1
    fi

    if module::has "${__name__}"; then
        log::error "Module is already registered: ${__name__}"
        return 1
    fi

    if [[ "${__builtin__}" != true && -n "${MODULE_LOADING_NAME}" && "${__name__}" != "${MODULE_LOADING_NAME}" ]]; then
        log::error "Module directory '${MODULE_LOADING_NAME}' registered as '${__name__}'"
        return 1
    fi

    if [[ "${__api__}" != "${RUNTIME_API_VERSION:-1}" ]]; then
        log::error "Module ${__name__} requires runtime API ${__api__}"
        return 1
    fi
}

module::_validate_handler()
{
    local __kind__=${1:-}
    local __handler__=${2:-}

    [[ -n "${__handler__}" ]] || return 0

    if ! declare -F "${__handler__}" >/dev/null 2>&1; then
        log::error "Module ${__kind__} handler does not exist: ${__handler__}"
        return 1
    fi
}

module::_create_context()
{
    local __context__=${1:-}
    local __name__=${2:-}
    local __root__=${3:-}

    declare -gA "${__context__}"
    local -n _context_=${__context__}
    _context_=()

    var::update _context_ module "${__name__}"
    var::update _context_ root "${__root__}"

}
