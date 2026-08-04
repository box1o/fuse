#!/usr/bin/env bash

archive::detect_format()
{
    local __archive__=${1:-}

    case "${__archive__}" in
        *.tar.gz | *.tgz) printf 'tar.gz\n' ;;
        *.tar.xz | *.txz) printf 'tar.xz\n' ;;
        *.tar.bz2 | *.tbz2) printf 'tar.bz2\n' ;;
        *.tar) printf 'tar\n' ;;
        *.zip) printf 'zip\n' ;;
        *)
            log::error "Unsupported archive format: ${__archive__}"
            return 2
            ;;
    esac
}

archive::_require_tool()
{
    local __format__=${1:-}

    case "${__format__}" in
        tar | tar.*) proc::require tar ;;
        zip) proc::require unzip ;;
        *) return 2 ;;
    esac
}

archive::list()
{
    local __archive__=${1:-}

    fs::require_file archive "${__archive__}" || return

    local __format__
    __format__=$(archive::detect_format "${__archive__}") || return
    archive::_require_tool "${__format__}" || return

    case "${__format__}" in
        tar | tar.*) tar -tf "${__archive__}" ;;
        zip) unzip -Z1 "${__archive__}" ;;
    esac
}

archive::verify()
{
    local __archive__=${1:-}
    local __entries__

    __entries__=$(archive::list "${__archive__}") || return

    local __entry__
    while IFS= read -r __entry__; do
        [[ -n "${__entry__}" ]] || continue

        if archive::_entry_is_unsafe "${__entry__}"; then
            log::error "Unsafe archive entry: ${__entry__}"
            return 1
        fi
    done <<<"${__entries__}"
}

archive::_entry_is_unsafe()
{
    local __entry__=${1:-}

    [[ "${__entry__}" == /* || "${__entry__}" == ../* || "${__entry__}" == */../* ]]
}

archive::extract()
{
    local __archive__=${1:-}
    local __destination__=${2:-.}

    fs::require_file archive "${__archive__}" || return

    local __format__
    __format__=$(archive::detect_format "${__archive__}") || return
    archive::_require_tool "${__format__}" || return
    archive::verify "${__archive__}" || return
    fs::mkdir "${__destination__}" || return

    case "${__format__}" in
        tar | tar.*) tar -xf "${__archive__}" -C "${__destination__}" ;;
        zip) unzip -q "${__archive__}" -d "${__destination__}" ;;
    esac
}

archive::create_tar_gz()
{
    archive::_create_tar z "${1:-.}" "${2:-archive.tar.gz}"
}

archive::create_tar_xz()
{
    archive::_create_tar J "${1:-.}" "${2:-archive.tar.xz}"
}

archive::create_tar()
{
    archive::_create_tar '' "${1:-.}" "${2:-archive.tar}"
}

archive::_create_tar()
{
    local __compression__=${1-}
    local __source__=${2:-.}
    local __destination__=${3:-archive.tar}

    fs::require_dir source "${__source__}" || return
    proc::require tar || return
    fs::mkdir_parent "${__destination__}" || return

    tar "-c${__compression__}f" "${__destination__}" -C "${__source__}" .
}
