#!/usr/bin/env bash

declare -gA ARG_KIND_REGISTRY=()
declare -gA ARG_SHORT_REGISTRY=()
declare -gA ARG_VALUE_REGISTRY=()
declare -gA ARG_SEEN_REGISTRY=()
declare -gA ARG_REQUIRED_REGISTRY=()
declare -gA ARG_DEFAULT_REGISTRY=()
declare -gA ARG_HAS_DEFAULT_REGISTRY=()
declare -ga ARG_POSITIONAL_VALUES=()

# Clear all argument definitions and parsed values.
args::reset()
{
    ARG_KIND_REGISTRY=()
    ARG_SHORT_REGISTRY=()
    ARG_VALUE_REGISTRY=()
    ARG_SEEN_REGISTRY=()
    ARG_REQUIRED_REGISTRY=()
    ARG_DEFAULT_REGISTRY=()
    ARG_HAS_DEFAULT_REGISTRY=()
    ARG_POSITIONAL_VALUES=()
}

# Register a boolean flag.
#
# Usage:
#   args::flag --name <name> [--short <character>]
args::flag()
{
    args::_define flag "$@"
}

# Register an option that consumes a value.
#
# Usage:
#   args::option --name <name> [--short <character>] [--default <value>] [--required]
args::option()
{
    args::_define option "$@"
}

args::_define()
{
    local __kind__=${1:-}
    shift || true

    local __name__=
    local __short__=
    local __default__=
    local __has_default__=false
    local __required__=false

    while (($# > 0)); do
        case "$1" in
            --name | --short | --default)
                args::_require_option_value args::_define "$1" "$#" || return

                case "$1" in
                    --name) __name__=$2 ;;
                    --short) __short__=$2 ;;
                    --default)
                        __default__=$2
                        __has_default__=true
                        ;;
                esac

                shift 2
                ;;
            --required)
                __required__=true
                shift
                ;;
            *)
                log::error "args::${__kind__}: unknown definition option: $1"
                return 2
                ;;
        esac
    done

    args::_validate_definition \
        "${__kind__}" \
        "${__name__}" \
        "${__short__}" \
        "${__has_default__}" \
        "${__required__}" || return

    ARG_KIND_REGISTRY["${__name__}"]=${__kind__}
    ARG_REQUIRED_REGISTRY["${__name__}"]=${__required__}
    ARG_HAS_DEFAULT_REGISTRY["${__name__}"]=${__has_default__}
    ARG_SEEN_REGISTRY["${__name__}"]=false

    if [[ -n "${__short__}" ]]; then
        ARG_SHORT_REGISTRY["${__short__}"]=${__name__}
    fi

    if [[ "${__kind__}" == flag ]]; then
        ARG_VALUE_REGISTRY["${__name__}"]=false
    elif [[ "${__has_default__}" == true ]]; then
        ARG_DEFAULT_REGISTRY["${__name__}"]=${__default__}
        ARG_VALUE_REGISTRY["${__name__}"]=${__default__}
    fi
}

args::_validate_definition()
{
    local __kind__=${1:-}
    local __name__=${2:-}
    local __short__=${3:-}
    local __has_default__=${4:-false}
    local __required__=${5:-false}

    validate::enum "${__kind__}" flag option || return

    if [[ ! "${__name__}" =~ ^[a-z][a-z0-9-]*$ ]]; then
        log::error "Invalid argument name: ${__name__}"
        return 2
    fi

    if [[ -v "ARG_KIND_REGISTRY[${__name__}]" ]]; then
        log::error "Argument is already defined: --${__name__}"
        return 1
    fi

    if [[ -n "${__short__}" ]]; then
        if [[ ! "${__short__}" =~ ^[a-zA-Z0-9]$ ]]; then
            log::error "Invalid short argument: -${__short__}"
            return 2
        fi

        if [[ -v "ARG_SHORT_REGISTRY[${__short__}]" ]]; then
            log::error "Short argument is already defined: -${__short__}"
            return 1
        fi
    fi

    if [[ "${__kind__}" == flag && ("${__has_default__}" == true || "${__required__}" == true) ]]; then
        log::error "Flags do not accept --default or --required: --${__name__}"
        return 2
    fi

    if [[ "${__has_default__}" == true && "${__required__}" == true ]]; then
        log::error "An option cannot be both required and have a default: --${__name__}"
        return 2
    fi
}

# Parse long options, single short options, and positional values.
args::parse()
{
    args::_reset_values

    local __argument__
    while (($# > 0)); do
        __argument__=$1

        case "${__argument__}" in
            --)
                shift
                ARG_POSITIONAL_VALUES+=("$@")
                break
                ;;
            --*=*)
                args::_parse_long_with_value "${__argument__}" || return
                shift
                ;;
            --*)
                args::_parse_named "${__argument__#--}" "${2-}" "$#" || return

                if [[ "${ARG_KIND_REGISTRY[${__argument__#--}]}" == option ]]; then
                    shift 2
                else
                    shift
                fi
                ;;
            -?)
                local __name__=${ARG_SHORT_REGISTRY[${__argument__#-}]:-}

                if [[ -z "${__name__}" ]]; then
                    log::error "Unknown option: ${__argument__}"
                    return 2
                fi

                args::_parse_named "${__name__}" "${2-}" "$#" || return

                if [[ "${ARG_KIND_REGISTRY[${__name__}]}" == option ]]; then
                    shift 2
                else
                    shift
                fi
                ;;
            -*)
                log::error "Combined short options are not supported: ${__argument__}"
                return 2
                ;;
            *)
                ARG_POSITIONAL_VALUES+=("${__argument__}")
                shift
                ;;
        esac
    done

    args::_validate_required
}

args::_reset_values()
{
    ARG_POSITIONAL_VALUES=()

    local __name__
    for __name__ in "${!ARG_KIND_REGISTRY[@]}"; do
        ARG_SEEN_REGISTRY["${__name__}"]=false

        if [[ "${ARG_KIND_REGISTRY[${__name__}]}" == flag ]]; then
            ARG_VALUE_REGISTRY["${__name__}"]=false
        elif [[ "${ARG_HAS_DEFAULT_REGISTRY[${__name__}]}" == true ]]; then
            ARG_VALUE_REGISTRY["${__name__}"]=${ARG_DEFAULT_REGISTRY[${__name__}]}
        else
            unset 'ARG_VALUE_REGISTRY['"${__name__}"']'
        fi
    done
}

args::_parse_long_with_value()
{
    local __argument__=${1:-}
    local __name__=${__argument__%%=*}
    local __value__=${__argument__#*=}
    __name__=${__name__#--}

    if [[ ! -v "ARG_KIND_REGISTRY[${__name__}]" ]]; then
        log::error "Unknown option: --${__name__}"
        return 2
    fi

    if [[ "${ARG_KIND_REGISTRY[${__name__}]}" != option ]]; then
        log::error "Flag does not accept a value: --${__name__}"
        return 2
    fi

    ARG_VALUE_REGISTRY["${__name__}"]=${__value__}
    ARG_SEEN_REGISTRY["${__name__}"]=true
}

args::_parse_named()
{
    local __name__=${1:-}
    local __next_value__=${2-}
    local __argument_count__=${3:-0}

    if [[ ! -v "ARG_KIND_REGISTRY[${__name__}]" ]]; then
        log::error "Unknown option: --${__name__}"
        return 2
    fi

    if [[ "${ARG_KIND_REGISTRY[${__name__}]}" == flag ]]; then
        ARG_VALUE_REGISTRY["${__name__}"]=true
        ARG_SEEN_REGISTRY["${__name__}"]=true
        return 0
    fi

    if ((__argument_count__ < 2)); then
        log::error "Option requires a value: --${__name__}"
        return 2
    fi

    ARG_VALUE_REGISTRY["${__name__}"]=${__next_value__}
    ARG_SEEN_REGISTRY["${__name__}"]=true
}

args::_validate_required()
{
    local __name__

    for __name__ in "${!ARG_REQUIRED_REGISTRY[@]}"; do
        if [[ "${ARG_REQUIRED_REGISTRY[${__name__}]}" == true && "${ARG_SEEN_REGISTRY[${__name__}]}" != true ]]; then
            log::error "Required option is missing: --${__name__}"
            return 2
        fi
    done
}

args::_require_option_value()
{
    local __caller__=${1:-function}
    local __option__=${2:-}
    local __argument_count__=${3:-0}

    if ((__argument_count__ < 2)); then
        log::error "${__caller__}: ${__option__} requires a value"
        return 2
    fi
}

args::has()
{
    local __name__=${1:-}

    args::_require_definition "${__name__}" || return
    [[ "${ARG_SEEN_REGISTRY[${__name__}]}" == true ]]
}

args::get()
{
    local __name__=${1:-}
    local __fallback__=${2-}

    args::_require_definition "${__name__}" || return

    if [[ -v "ARG_VALUE_REGISTRY[${__name__}]" ]]; then
        printf '%s\n' "${ARG_VALUE_REGISTRY[${__name__}]}"
    elif (($# >= 2)); then
        printf '%s\n' "${__fallback__}"
    else
        log::error "Argument has no value: --${__name__}"
        return 1
    fi
}

args::_require_definition()
{
    local __name__=${1:-}

    if [[ ! -v "ARG_KIND_REGISTRY[${__name__}]" ]]; then
        log::error "Argument is not defined: ${__name__}"
        return 1
    fi
}

args::positionals()
{
    ((${#ARG_POSITIONAL_VALUES[@]} > 0)) || return 0
    printf '%s\n' "${ARG_POSITIONAL_VALUES[@]}"
}

args::position()
{
    local __index__=${1:-}

    validate::non_negative_integer position "${__index__}" || return

    if ((__index__ >= ${#ARG_POSITIONAL_VALUES[@]})); then
        log::error "Positional argument does not exist: ${__index__}"
        return 1
    fi

    printf '%s\n' "${ARG_POSITIONAL_VALUES[${__index__}]}"
}

args::count()
{
    printf '%s\n' "${#ARG_POSITIONAL_VALUES[@]}"
}
