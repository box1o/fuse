#!/usr/bin/env bash

# Set a key in a caller-owned associative array.
#
# Usage:
#   var::update <array-reference> <key> <value>
var::update()
{
    local __reference__=${1:-}
    local __key__=${2:-}
    local __value__=${3-}

    var::_validate_reference "${__reference__}" || return
    validate::identifier "${__key__}" || return

    # The validated name refers to an associative array.
    # shellcheck disable=SC2178
    local -n _variables_=${__reference__}
    _variables_["${__key__}"]=${__value__}
}

var::has()
{
    local __reference__=${1:-}
    local __key__=${2:-}

    var::_validate_reference "${__reference__}" || return
    validate::identifier "${__key__}" || return

    # The validated name refers to an associative array.
    # shellcheck disable=SC2178
    local -n _variables_=${__reference__}
    [[ -n ${_variables_[${__key__}]+defined} ]]
}

var::get()
{
    local __reference__=${1:-}
    local __key__=${2:-}
    local __default__=${3-}

    var::_validate_reference "${__reference__}" || return
    validate::identifier "${__key__}" || return

    # The validated name refers to an associative array.
    # shellcheck disable=SC2178
    local -n _variables_=${__reference__}

    if [[ -n ${_variables_[${__key__}]+defined} ]]; then
        printf '%s\n' "${_variables_[${__key__}]}"
    elif (($# >= 3)); then
        printf '%s\n' "${__default__}"
    else
        log::error "Variable is not defined: ${__key__}"
        return 1
    fi
}

var::_validate_reference()
{
    local __reference__=${1:-}

    if [[ ! "${__reference__}" =~ ^[a-zA-Z_][a-zA-Z0-9_]*$ ]]; then
        log::error "Invalid variable reference: ${__reference__}"
        return 2
    fi
}
