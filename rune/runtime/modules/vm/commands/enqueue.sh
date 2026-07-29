#!/usr/bin/env bash

firecracker::enqueue()
{
    args::reset
    args::option --name id --required
    args::option --name workdir --default /workspace
    args::parse "$@" || return

    if (($(args::count) < 2)); then
        log::error "Usage: rune vm enqueue <vm> --id ID [--workdir PATH] -- command [arguments...]"
        return 2
    fi

    local __name__
    local __job__
    local __workdir__
    local -a __command__=()
    __name__=$(args::position 0) || return
    __job__=$(args::get id) || return
    __workdir__=$(args::get workdir) || return
    mapfile -t __command__ < <(args::positionals | tail -n +2)
    firecracker::_valid_job_id "${__job__}" || return
    [[ "${__workdir__}" == /* ]] || {
        log::error "Job workdir must be absolute: ${__workdir__}"
        return 2
    }

    local __lock__
    __lock__=$(firecracker::_lock "${__name__}") || return
    lock::with "${__lock__}" firecracker::_enqueue \
        "${__name__}" "${__job__}" "${__workdir__}" "${__command__[@]}"
}

firecracker::_enqueue()
{
    local __name__=${1:-}
    local __job__=${2:-}
    local __workdir__=${3:-}
    shift 3 || true
    firecracker::_require_vm "${__name__}" || return
    firecracker::_require_stopped "${__name__}" || return
    (($# > 0)) || return 2
    proc::require debugfs || return

    local __script__
    __script__=$(temp::file vm-job) || return
    {
        printf '#!/usr/bin/env bash\n\n'
        printf 'exec /usr/local/bin/rune-guest run --id %q --workdir %q --' "${__job__}" "${__workdir__}"
        printf ' %q' "$@"
        printf '\n'
    } | fs::atomic_write "${__script__}" || return

    local __rootfs__
    local __destination__="/var/lib/rune/queue/${__job__}.job"
    __rootfs__=$(firecracker::_rootfs "${__name__}") || return
    fs::require_file rootfs "${__rootfs__}" || return
    debugfs -w -R "rm ${__destination__}" "${__rootfs__}" >/dev/null 2>&1 || true
    debugfs -w -R "write ${__script__} ${__destination__}" "${__rootfs__}" >/dev/null 2>&1 || {
        log::error "Unable to write job into VM rootfs: ${__name__}/${__job__}"
        return 1
    }
    debugfs -R "stat ${__destination__}" "${__rootfs__}" >/dev/null 2>&1 || {
        log::error "Queued job verification failed: ${__name__}/${__job__}"
        return 1
    }
    log::info "Queued VM job: ${__name__}/${__job__}"
}
