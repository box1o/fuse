#!/usr/bin/env bash

system::os()
{
    uname -s
}

system::architecture()
{
    uname -m
}

system::kernel()
{
    uname -r
}

system::hostname()
{
    hostname
}

system::cpu_count()
{
    getconf _NPROCESSORS_ONLN 2>/dev/null || nproc
}

system::memory_mb()
{
    awk '/MemTotal:/ { print int($2 / 1024) }' /proc/meminfo
}

system::require_linux()
{
    local __operating_system__
    __operating_system__=$(system::os) || return

    if [[ "${__operating_system__}" != Linux ]]; then
        log::error "Linux is required"
        return 1
    fi
}

system::distribution()
{
    if [[ ! -r /etc/os-release ]]; then
        printf 'unknown\n'
        return 0
    fi

    local __distribution__
    __distribution__=$(
        # shellcheck disable=SC1091
        source /etc/os-release
        printf '%s\n' "${ID:-unknown}"
    ) || return

    printf '%s\n' "${__distribution__}"
}

system::is_container()
{
    [[ -f /.dockerenv ]] && return 0
    grep -qaE '(docker|containerd|lxc|kubepods)' /proc/1/cgroup 2>/dev/null
}

system::is_wsl()
{
    grep -qi microsoft /proc/version 2>/dev/null
}
