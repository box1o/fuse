#!/usr/bin/env bash

resource::validate_cpu()
{
    validate::positive_integer cpu "${1:-}"
}

resource::validate_memory()
{
    validate::positive_integer memory "${1:-}"
}

resource::host_cpu()
{
    proc::require nproc || return
    nproc
}

resource::host_memory()
{
    resource::_meminfo MemTotal
}

resource::host_available_memory()
{
    resource::_meminfo MemAvailable
}

resource::_meminfo()
{
    local __field__=${1:-}
    local __value__
    __value__=$(awk -v field="${__field__}:" '$1 == field { print int($2 / 1024); exit }' /proc/meminfo) || return
    validate::not_empty "memory field ${__field__}" "${__value__}" || return
    printf '%s\n' "${__value__}"
}

resource::process_usage()
{
    local __pid__=${1:-}
    validate::positive_integer pid "${__pid__}" || return

    if ! kill -0 "${__pid__}" 2>/dev/null; then
        log::error "Process does not exist: ${__pid__}"
        return 1
    fi

    ps -o pid=,pcpu=,rss= -p "${__pid__}" |
        awk '{ printf "pid: %s\ncpu_percent: %s\nmemory_mib: %.1f\n", $1, $2, $3 / 1024 }'
}

resource::format_memory()
{
    local __memory__=${1:-}
    resource::validate_memory "${__memory__}" || return
    printf '%s MiB\n' "${__memory__}"
}

resource::host()
{
    if (($# != 0)); then
        log::error "Usage: rune resource host"
        return 2
    fi

    printf 'cpu: %s\n' "$(resource::host_cpu)" || return
    printf 'memory_mib: %s\n' "$(resource::host_memory)" || return
    printf 'available_memory_mib: %s\n' "$(resource::host_available_memory)"
}

resource::usage()
{
    if (($# != 1)); then
        log::error "Usage: rune resource usage <pid>"
        return 2
    fi

    resource::process_usage "$1"
}

resource::validate()
{
    args::reset
    args::option --name cpu --required
    args::option --name memory --required
    args::parse "$@" || return

    if (($(args::count) != 0)); then
        log::error "Usage: rune resource validate --cpu N --memory MiB"
        return 2
    fi

    resource::validate_cpu "$(args::get cpu)" || return
    resource::validate_memory "$(args::get memory)" || return
    log::info "Resource values are valid"
}
