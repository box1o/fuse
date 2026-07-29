#!/usr/bin/env bash

firecracker::_valid_job_id()
{
    if [[ ! ${1:-} =~ ^[a-zA-Z0-9][a-zA-Z0-9._-]{0,127}$ ]]; then
        log::error "Invalid job ID: ${1:-}"
        return 2
    fi
}

firecracker::job_logs()
{
    args::reset
    args::option --name lines --default 100
    args::option --name source --default auto
    args::parse "$@" || return

    if (($(args::count) != 2)); then
        log::error "Usage: rune vm job-logs <vm> <job> [--lines N] [--source auto|console|disk]"
        return 2
    fi

    local __name__
    local __job__
    local __lines__
    local __source__
    __name__=$(args::position 0) || return
    __job__=$(args::position 1) || return
    __lines__=$(args::get lines) || return
    __source__=$(args::get source) || return
    firecracker::_require_vm "${__name__}" || return
    firecracker::_valid_job_id "${__job__}" || return
    validate::positive_integer lines "${__lines__}" || return
    validate::enum "${__source__}" auto console disk || return

    if [[ "${__source__}" == auto ]]; then
        if firecracker::_is_running "${__name__}"; then
            __source__=console
        else
            __source__=disk
        fi
    fi

    if [[ "${__source__}" == console ]]; then
        firecracker::_job_logs_console "${__name__}" "${__job__}" "${__lines__}"
    else
        firecracker::_job_logs_disk "${__name__}" "${__job__}" "${__lines__}"
    fi
}

firecracker::_job_logs_console()
{
    local __name__=${1:-}
    local __job__=${2:-}
    local __lines__=${3:-100}
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    fs::require_file console-log "${__vm__}/console.log" || return

    awk -v prefix="RUNE_JOB_LOG ${__job__} " '
        index($0, prefix) == 1 {
            print substr($0, length(prefix) + 1)
        }
    ' "${__vm__}/console.log" | tail -n "${__lines__}"
}

firecracker::_job_logs_disk()
{
    local __name__=${1:-}
    local __job__=${2:-}
    local __lines__=${3:-100}
    firecracker::_require_stopped "${__name__}" || return
    proc::require debugfs || return

    local __rootfs__
    __rootfs__=$(firecracker::_rootfs "${__name__}") || return
    fs::require_file rootfs "${__rootfs__}" || return
    debugfs -R "cat /var/log/rune/jobs/${__job__}/output.log" "${__rootfs__}" 2>/dev/null |
        tail -n "${__lines__}"
}

firecracker::job_status()
{
    if (($# != 2)); then
        log::error "Usage: rune vm job-status <vm> <job>"
        return 2
    fi

    local __name__=$1
    local __job__=$2
    firecracker::_require_vm "${__name__}" || return
    firecracker::_valid_job_id "${__job__}" || return

    local __value__
    if firecracker::_is_running "${__name__}"; then
        __value__=$(firecracker::_job_status_console "${__name__}" "${__job__}") || return
    else
        __value__=$(firecracker::_job_status_disk "${__name__}" "${__job__}") || return
    fi
    firecracker::_print_job_status "${__name__}" "${__job__}" "${__value__}"
}

firecracker::_job_status_console()
{
    local __name__=${1:-}
    local __job__=${2:-}
    local __vm__
    __vm__=$(firecracker::_vm_dir "${__name__}") || return
    fs::require_file console-log "${__vm__}/console.log" || return
    awk -v job="${__job__}" '
        $1 == "RUNE_JOB_BEGIN" && $2 == job { status = "running" }
        $1 == "RUNE_JOB_END" && $2 == job { status = $3 }
        END {
            if (status == "") exit 1
            print status
        }
    ' "${__vm__}/console.log" || {
        log::error "Job has not been observed: ${__name__}/${__job__}"
        return 1
    }
}

firecracker::_job_status_disk()
{
    local __name__=${1:-}
    local __job__=${2:-}
    proc::require debugfs || return
    local __rootfs__
    __rootfs__=$(firecracker::_rootfs "${__name__}") || return
    debugfs -R "cat /var/log/rune/jobs/${__job__}/status" "${__rootfs__}" 2>/dev/null || {
        log::error "Job status is unavailable: ${__name__}/${__job__}"
        return 1
    }
}

firecracker::_print_job_status()
{
    local __name__=${1:-}
    local __job__=${2:-}
    local __value__=${3:-}
    local __state__=completed
    local __exit_code__=${__value__}
    if [[ "${__value__}" == running ]]; then
        __state__=running
        __exit_code__=null
    elif [[ ! "${__value__}" =~ ^[0-9]+$ ]]; then
        log::error "Invalid guest job status: ${__value__}"
        return 1
    fi

    if [[ $(runtime::config::get 'output|format' text) == json ]]; then
        printf '{"vm":"%s","job":"%s","state":"%s","exit_code":%s}\n' \
            "$(firecracker::_json_string "${__name__}")" \
            "$(firecracker::_json_string "${__job__}")" \
            "${__state__}" "${__exit_code__}"
    elif [[ "${__state__}" == running ]]; then
        printf '%s/%s running\n' "${__name__}" "${__job__}"
    else
        printf '%s/%s completed exit_code=%s\n' "${__name__}" "${__job__}" "${__exit_code__}"
    fi
}

firecracker::wait_job()
{
    args::reset
    args::option --name timeout --default 3600
    args::parse "$@" || return
    if (($(args::count) != 2)); then
        log::error "Usage: rune vm wait-job <vm> <job> [--timeout SECONDS]"
        return 2
    fi

    local __name__
    local __job__
    local __timeout__
    __name__=$(args::position 0) || return
    __job__=$(args::position 1) || return
    __timeout__=$(args::get timeout) || return
    firecracker::_require_vm "${__name__}" || return
    firecracker::_valid_job_id "${__job__}" || return
    validate::positive_integer timeout "${__timeout__}" || return

    local __deadline__=$((SECONDS + __timeout__))
    local __value__
    while ((SECONDS < __deadline__)); do
        if firecracker::_is_running "${__name__}"; then
            __value__=$(firecracker::_job_status_console "${__name__}" "${__job__}" 2>/dev/null) || __value__=
        else
            __value__=$(firecracker::_job_status_disk "${__name__}" "${__job__}" 2>/dev/null) || __value__=
        fi
        if [[ "${__value__}" =~ ^[0-9]+$ ]]; then
            firecracker::_print_job_status "${__name__}" "${__job__}" "${__value__}" || return
            return "${__value__}"
        fi
        sleep 0.5
    done
    log::error "Timed out waiting for job: ${__name__}/${__job__}"
    return 124
}
