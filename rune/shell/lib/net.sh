#!/usr/bin/env bash

NET_DEFAULT_CONNECT_TIMEOUT=${NET_DEFAULT_CONNECT_TIMEOUT:-15}
NET_DEFAULT_RETRIES=${NET_DEFAULT_RETRIES:-5}
NET_DEFAULT_RETRY_DELAY=${NET_DEFAULT_RETRY_DELAY:-2}
NET_DEFAULT_USER_AGENT=${NET_DEFAULT_USER_AGENT:-rune-shell/1.0}

net::_require_curl()
{
    proc::require curl
}

net::is_url()
{
    [[ ${1:-} =~ ^https?:// ]]
}

net::status()
{
    local __url__=${1:-}

    validate::not_empty URL "${__url__}" || return
    net::_require_curl || return

    curl \
        --silent \
        --show-error \
        --location \
        --output /dev/null \
        --write-out '%{http_code}\n' \
        --connect-timeout "${NET_DEFAULT_CONNECT_TIMEOUT}" \
        --user-agent "${NET_DEFAULT_USER_AGENT}" \
        -- "${__url__}"
}

net::is_reachable()
{
    local __url__=${1:-}

    validate::not_empty URL "${__url__}" || return
    net::_require_curl || return

    curl \
        --silent \
        --show-error \
        --fail \
        --location \
        --head \
        --connect-timeout "${NET_DEFAULT_CONNECT_TIMEOUT}" \
        --user-agent "${NET_DEFAULT_USER_AGENT}" \
        -- "${__url__}" \
        >/dev/null
}

# Download a URL through a resumable partial file.
#
# Usage:
#   net::download [options] <url> <destination>
net::download()
{
    local __retries__=${NET_DEFAULT_RETRIES}
    local __retry_delay__=${NET_DEFAULT_RETRY_DELAY}
    local __connect_timeout__=${NET_DEFAULT_CONNECT_TIMEOUT}
    local __overwrite__=false
    local __resume__=true
    local __quiet__=false
    local __sha256__=
    local -a __headers__=()

    while (($# > 0)); do
        case "$1" in
            --overwrite)
                __overwrite__=true
                shift
                ;;
            --no-resume)
                __resume__=false
                shift
                ;;
            --quiet)
                __quiet__=true
                shift
                ;;
            --retries | --retry-delay | --connect-timeout | --header | --sha256)
                net::_require_option_value net::download "$1" "$#" || return

                case "$1" in
                    --retries) __retries__=$2 ;;
                    --retry-delay) __retry_delay__=$2 ;;
                    --connect-timeout) __connect_timeout__=$2 ;;
                    --header) __headers__+=("$2") ;;
                    --sha256) __sha256__=$2 ;;
                esac

                shift 2
                ;;
            --)
                shift
                break
                ;;
            -*)
                log::error "net::download: unknown option: $1"
                return 2
                ;;
            *) break ;;
        esac
    done

    if (($# != 2)); then
        log::error "Usage: net::download [options] <url> <destination>"
        return 2
    fi

    local __url__=$1
    local __destination__=$2
    local __partial__="${__destination__}.part"

    net::_validate_download \
        "${__url__}" \
        "${__destination__}" \
        "${__overwrite__}" \
        "${__retries__}" \
        "${__retry_delay__}" \
        "${__connect_timeout__}" || return

    fs::mkdir_parent "${__destination__}" || return

    local -a __arguments__=(
        --fail
        --show-error
        --location
        --retry "${__retries__}"
        --retry-delay "${__retry_delay__}"
        --retry-connrefused
        --connect-timeout "${__connect_timeout__}"
        --user-agent "${NET_DEFAULT_USER_AGENT}"
        --output "${__partial__}"
    )

    if [[ "${__resume__}" == true && -f "${__partial__}" ]]; then
        __arguments__+=(--continue-at -)
    fi

    if [[ "${__quiet__}" == true ]]; then
        __arguments__+=(--silent)
    else
        __arguments__+=(--progress-bar)
    fi

    net::_append_headers __arguments__ __headers__

    local __status__=0
    curl "${__arguments__[@]}" -- "${__url__}" || __status__=$?

    if ((__status__ != 0)); then
        log::error "Download failed; partial file retained: ${__partial__}"
        return "${__status__}"
    fi

    if [[ -n "${__sha256__}" ]]; then
        net::verify_sha256 "${__partial__}" "${__sha256__}" || return
    fi

    mv -f -- "${__partial__}" "${__destination__}"
}

net::_require_option_value()
{
    local __caller__=${1:-function}
    local __option__=${2:-}
    local __argument_count__=${3:-0}

    if ((__argument_count__ < 2)); then
        log::error "${__caller__}: ${__option__} requires a value"
        return 2
    fi
}

net::_validate_download()
{
    local __url__=${1:-}
    local __destination__=${2:-}
    local __overwrite__=${3:-false}
    local __retries__=${4:-}
    local __retry_delay__=${5:-}
    local __connect_timeout__=${6:-}

    if ! net::is_url "${__url__}"; then
        log::error "Unsupported URL: ${__url__}"
        return 2
    fi

    validate::not_empty destination "${__destination__}" || return
    validate::non_negative_integer retries "${__retries__}" || return
    validate::non_negative_integer retry-delay "${__retry_delay__}" || return
    validate::non_negative_integer connect-timeout "${__connect_timeout__}" || return
    net::_require_curl || return

    if [[ -e "${__destination__}" && "${__overwrite__}" != true ]]; then
        log::error "Destination already exists: ${__destination__}"
        return 1
    fi
}

net::_append_headers()
{
    local -n _arguments_=$1
    local -n _headers_=$2
    local __header__

    for __header__ in "${_headers_[@]}"; do
        _arguments_+=(--header "${__header__}")
    done
}

net::verify_sha256()
{
    local __file__=${1:-}
    local __expected__=${2:-}

    fs::require_file file "${__file__}" || return
    validate::not_empty expected-sha256 "${__expected__}" || return

    local __actual__
    __actual__=$(fs::checksum "${__file__}") || return

    if [[ "${__actual__,,}" != "${__expected__,,}" ]]; then
        log::error "SHA-256 mismatch: expected ${__expected__}, got ${__actual__}"
        return 1
    fi
}

net::request()
{
    local __method__=GET
    local __output__=
    local __data__=
    local -a __headers__=()

    while (($# > 0)); do
        case "$1" in
            --method | --header | --data | --output)
                net::_require_option_value net::request "$1" "$#" || return

                case "$1" in
                    --method) __method__=$2 ;;
                    --header) __headers__+=("$2") ;;
                    --data) __data__=$2 ;;
                    --output) __output__=$2 ;;
                esac

                shift 2
                ;;
            --)
                shift
                break
                ;;
            -*)
                log::error "net::request: unknown option: $1"
                return 2
                ;;
            *) break ;;
        esac
    done

    if (($# != 1)); then
        log::error "Usage: net::request [options] <url>"
        return 2
    fi

    local __url__=$1
    validate::not_empty URL "${__url__}" || return
    net::is_url "${__url__}" || {
        log::error "Unsupported URL: ${__url__}"
        return 2
    }
    net::_require_curl || return

    local -a __arguments__=(
        --fail-with-body
        --show-error
        --location
        --request "${__method__}"
    )

    net::_append_headers __arguments__ __headers__

    if [[ -n "${__data__}" ]]; then
        __arguments__+=(--data "${__data__}")
    fi

    if [[ -n "${__output__}" ]]; then
        fs::mkdir_parent "${__output__}" || return
        __arguments__+=(--output "${__output__}")
    fi

    curl "${__arguments__[@]}" -- "${__url__}"
}
