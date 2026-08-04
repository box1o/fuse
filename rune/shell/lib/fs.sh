#!/usr/bin/env bash

fs::exists()
{
    local __path__=${1:-}

    [[ -e "${__path__}" || -L "${__path__}" ]]
}

fs::is_file()
{
    [[ -f "${1:-}" ]]
}

fs::is_dir()
{
    [[ -d "${1:-}" ]]
}

fs::is_symlink()
{
    [[ -L "${1:-}" ]]
}

fs::is_readable()
{
    [[ -r "${1:-}" ]]
}

fs::is_writable()
{
    [[ -w "${1:-}" ]]
}

fs::is_executable()
{
    [[ -x "${1:-}" ]]
}

fs::mkdir()
{
    local __directory__=${1:-}

    validate::not_empty directory "${__directory__}" || return
    mkdir -p -- "${__directory__}"
}

fs::mkdir_parent()
{
    local __path__=${1:-}
    local __parent__

    validate::not_empty path "${__path__}" || return
    __parent__=$(path::dirname "${__path__}") || return

    fs::mkdir "${__parent__}"
}

fs::touch()
{
    local __file__=${1:-}

    validate::not_empty file "${__file__}" || return
    fs::mkdir_parent "${__file__}" || return

    touch -- "${__file__}"
}

fs::require_file()
{
    local __description__=${1:-file}
    local __file__=${2:-}

    if [[ ! -f "${__file__}" ]]; then
        log::error "${__description__} does not exist: ${__file__}"
        return 1
    fi
}

fs::require_dir()
{
    local __description__=${1:-directory}
    local __directory__=${2:-}

    if [[ ! -d "${__directory__}" ]]; then
        log::error "${__description__} does not exist: ${__directory__}"
        return 1
    fi
}

fs::_require_source()
{
    local __source__=${1:-}

    if ! fs::exists "${__source__}"; then
        log::error "Source does not exist: ${__source__}"
        return 1
    fi
}

fs::copy_file()
{
    local __source__=${1:-}
    local __destination__=${2:-}

    fs::require_file source "${__source__}" || return
    validate::not_empty destination "${__destination__}" || return
    fs::mkdir_parent "${__destination__}" || return

    cp --preserve=mode,timestamps -- \
        "${__source__}" \
        "${__destination__}"
}

fs::copy_tree()
{
    local __source__=${1:-}
    local __destination__=${2:-}

    fs::require_dir source "${__source__}" || return
    validate::not_empty destination "${__destination__}" || return
    fs::mkdir_parent "${__destination__}" || return

    cp -a -- "${__source__}" "${__destination__}"
}

fs::copy_into()
{
    local __source__=${1:-}
    local __destination_directory__=${2:-}

    fs::_require_source "${__source__}" || return
    fs::mkdir "${__destination_directory__}" || return

    cp -a -- "${__source__}" "${__destination_directory__}/"
}

fs::move()
{
    local __source__=${1:-}
    local __destination__=${2:-}

    fs::_require_source "${__source__}" || return
    validate::not_empty destination "${__destination__}" || return
    fs::mkdir_parent "${__destination__}" || return

    mv -- "${__source__}" "${__destination__}"
}

fs::rename()
{
    fs::move "$@"
}

fs::_assert_safe_removal_path()
{
    local __path__=${1:-}

    case "${__path__}" in
        '' | / | . | ..)
            log::error "Refusing unsafe removal path: '${__path__}'"
            return 1
            ;;
    esac

    validate::absolute_path "${__path__}"
}

fs::remove_file()
{
    local __path__=${1:-}

    fs::_assert_safe_removal_path "${__path__}" || return
    fs::exists "${__path__}" || return 0

    if [[ -d "${__path__}" && ! -L "${__path__}" ]]; then
        log::error "Path is a directory: ${__path__}"
        return 1
    fi

    rm -f -- "${__path__}"
}

fs::remove_dir()
{
    local __path__=${1:-}

    fs::_assert_safe_removal_path "${__path__}" || return
    [[ -d "${__path__}" ]] || return 0

    rmdir -- "${__path__}"
}

fs::remove_tree()
{
    local __path__=${1:-}

    fs::_assert_safe_removal_path "${__path__}" || return
    fs::exists "${__path__}" || return 0

    rm -rf -- "${__path__}"
}

fs::empty_dir()
{
    local __path__=${1:-}

    fs::_assert_safe_removal_path "${__path__}" || return
    fs::remove_tree "${__path__}" || return
    fs::mkdir "${__path__}"
}

fs::symlink()
{
    local __target__=${1:-}
    local __link__=${2:-}

    validate::not_empty target "${__target__}" || return
    validate::not_empty link "${__link__}" || return
    fs::mkdir_parent "${__link__}" || return

    ln -sfn -- "${__target__}" "${__link__}"
}

fs::readlink()
{
    readlink -f -- "${1:-}"
}

fs::chmod()
{
    local __mode__=${1:-}
    local __path__=${2:-}

    validate::not_empty mode "${__mode__}" || return
    validate::not_empty path "${__path__}" || return

    chmod -- "${__mode__}" "${__path__}"
}

fs::chown()
{
    local __owner__=${1:-}
    local __path__=${2:-}

    validate::not_empty owner "${__owner__}" || return
    validate::not_empty path "${__path__}" || return

    chown -- "${__owner__}" "${__path__}"
}

fs::file_size()
{
    stat -c '%s' -- "${1:-}"
}

fs::checksum()
{
    local __file__=${1:-}
    local __checksum__

    __checksum__=$(sha256sum -- "${__file__}") || return
    printf '%s\n' "${__checksum__%% *}"
}

fs::same_file()
{
    cmp -s -- "${1:-}" "${2:-}"
}

fs::write()
{
    local __destination__=${1:-}

    validate::not_empty destination "${__destination__}" || return
    fs::mkdir_parent "${__destination__}" || return

    cat >"${__destination__}"
}

fs::append()
{
    local __destination__=${1:-}

    validate::not_empty destination "${__destination__}" || return
    fs::mkdir_parent "${__destination__}" || return

    cat >>"${__destination__}"
}

# Replace a file only after all input has been written successfully.
fs::atomic_write()
{
    local __destination__=${1:-}

    validate::not_empty destination "${__destination__}" || return

    local __directory__
    __directory__=$(path::dirname "${__destination__}") || return
    fs::mkdir "${__directory__}" || return

    local __temporary__
    __temporary__=$(mktemp "${__directory__}/.tmp.XXXXXX") || return

    if ! cat >"${__temporary__}"; then
        rm -f -- "${__temporary__}"
        return 1
    fi

    mv -f -- "${__temporary__}" "${__destination__}"
}
