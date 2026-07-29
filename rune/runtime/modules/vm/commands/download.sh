#!/usr/bin/env bash

firecracker::download()
{
    local __default_version__
    __default_version__=$(module::var vm default_version) || return

    args::reset
    args::option --name version --default "${__default_version__}"
    args::option --name output
    args::option --name sha256
    args::flag --name force
    args::parse "$@" || return

    if (($(args::count) > 0)); then
        log::error "firecracker download: unexpected argument: $(args::position 0)"
        return 2
    fi

    local __version__
    local __output__
    local __checksum__
    __version__=$(args::get version) || return
    __output__=$(args::get output '') || return
    __checksum__=$(args::get sha256 '') || return

    local __releases_url__
    __releases_url__=$(module::var vm releases_url) || return

    if [[ "${__version__}" == latest ]]; then
        local __latest_url__
        __latest_url__=$(proc::capture \
            curl -fsSLI -o /dev/null -w '%{url_effective}' \
            "${__releases_url__}/latest") || return
        __version__=${__latest_url__##*/}
    elif [[ "${__version__}" != v* ]]; then
        __version__="v${__version__}"
    fi

    local __architecture__
    __architecture__=$(system::architecture) || return

    if [[ -z "${__output__}" ]]; then
        local __downloads__
        __downloads__=$(module::var vm downloads) || return
        __output__="${__downloads__}/firecracker-${__version__}-${__architecture__}.tgz"
    fi

    __output__=$(path::absolute "${__output__}") || return

    local __url__
    printf -v __url__ '%s/download/%s/firecracker-%s-%s.tgz' \
        "${__releases_url__}" \
        "${__version__}" \
        "${__version__}" \
        "${__architecture__}"

    local -a __download_args__=()
    args::has force && __download_args__+=(--overwrite)
    [[ -n "${__checksum__}" ]] && __download_args__+=(--sha256 "${__checksum__}")

    log::info "Downloading Firecracker: ${__url__}"
    net::download "${__download_args__[@]}" "${__url__}" "${__output__}"
}
