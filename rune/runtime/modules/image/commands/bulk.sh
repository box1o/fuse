#!/usr/bin/env bash

image::build_all()
{
    image::_bulk build "$@"
}

image::ensure_all()
{
    image::_bulk ensure "$@"
}

image::_bulk()
{
    local __operation__=${1:-}
    shift || true
    local __file__=
    local -a __options__=()

    while (($# > 0)); do
        case "$1" in
            --file)
                [[ $# -ge 2 ]] || {
                    log::error "image ${__operation__}-all: --file requires a value"
                    return 2
                }
                __file__=$2
                shift 2
                ;;
            --force | --dry-run)
                __options__+=("$1")
                shift
                ;;
            *)
                log::error "image ${__operation__}-all: unknown argument: $1"
                return 2
                ;;
        esac
    done

    if [[ "${__operation__}" == ensure && " ${__options__[*]} " == *' --force '* ]]; then
        log::error "image ensure-all does not accept --force"
        return 2
    fi

    image::_load_config __file__ "${__file__}" || return
    image::validate --file "${__file__}" || return

    local __name__
    while IFS= read -r __name__; do
        log::info "${__operation__^} image: ${__name__}"
        "image::${__operation__}" "${__name__}" --file "${__file__}" "${__options__[@]}" || return
    done < <(config::keys rune '.images')
}

image::catalog()
{
    local __file__=
    if (($# > 0)); then
        [[ $# -eq 2 && $1 == --file ]] || {
            log::error "Usage: rune image catalog [--file PATH]"
            return 2
        }
        __file__=$2
    fi
    image::_load_config __file__ "${__file__}" || return
    image::validate --file "${__file__}" || return

    printf '%-24s %-12s %-10s %s\n' IMAGE VERSION STATE KERNEL
    local __name__
    local __version__
    local __kernel__
    local __cache__
    local __state__
    __cache__=$(module::var image cache) || return
    while IFS= read -r __name__; do
        __version__=$(config::get rune ".images.\"${__name__}\".version") || return
        __kernel__=$(config::get rune ".images.\"${__name__}\".kernel") || return
        __state__=missing
        [[ -f "${__cache__}/${__name__}/${__version__}/rootfs.ext4" ]] && __state__=cached
        printf '%-24s %-12s %-10s %s\n' "${__name__}" "${__version__}" "${__state__}" "${__kernel__}"
    done < <(config::keys rune '.images')
}
