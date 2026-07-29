#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    GUEST="${PROJECT_ROOT}/runtime/modules/image/assets/guest/rune-guest"
    LOG_ROOT="${BATS_TEST_TMPDIR}/jobs"
}

@test "guest runner captures output and exit status" {
    run env RUNE_GUEST_LIB=/nonexistent RUNE_GUEST_LOG_DIR="${LOG_ROOT}" \
        "${GUEST}" run --id compile --workdir "${BATS_TEST_TMPDIR}/work" -- \
        bash -c 'printf "building\n"; exit 7'

    [ "${status}" -eq 7 ]
    [ "$(<"${LOG_ROOT}/compile/status")" = 7 ]
    [ "$(<"${LOG_ROOT}/compile/output.log")" = building ]
}

@test "guest runner lists and reads job logs" {
    env RUNE_GUEST_LIB=/nonexistent RUNE_GUEST_LOG_DIR="${LOG_ROOT}" \
        "${GUEST}" run --id render --workdir "${BATS_TEST_TMPDIR}/work" -- printf 'frame 1\n' >/dev/null

    run env RUNE_GUEST_LIB=/nonexistent RUNE_GUEST_LOG_DIR="${LOG_ROOT}" "${GUEST}" list
    [ "${status}" -eq 0 ]
    [ "${output}" = render ]

    run env RUNE_GUEST_LIB=/nonexistent RUNE_GUEST_LOG_DIR="${LOG_ROOT}" "${GUEST}" logs render
    [ "${status}" -eq 0 ]
    [ "${output}" = "frame 1" ]
}

@test "image installs the Rune library and guest runner" {
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
    module::load "${PROJECT_ROOT}/runtime/modules/image"
    module::initialize image
    root="${BATS_TEST_TMPDIR}/root"
    mkdir -p "${root}"

    image::_install_guest_runtime "${root}"

    [ -r "${root}/usr/local/lib/rune/shell/lib/init.sh" ]
    [ -x "${root}/usr/local/bin/rune-guest" ]
    [ -x "${root}/usr/local/bin/rune-guest-queue" ]
    [ -L "${root}/etc/systemd/system/multi-user.target.wants/rune-guest-queue.service" ]
    [ -d "${root}/var/log/rune/jobs" ]
}

@test "guest queue executes queued jobs and records completion" {
    queue="${BATS_TEST_TMPDIR}/queue"
    running="${BATS_TEST_TMPDIR}/running"
    completed="${BATS_TEST_TMPDIR}/completed"
    mkdir -p "${queue}"
    printf '#!/usr/bin/env bash\nprintf queued >%q\n' "${BATS_TEST_TMPDIR}/result" >"${queue}/test.job"

    RUNE_GUEST_QUEUE_DIR="${queue}" \
        RUNE_GUEST_RUNNING_DIR="${running}" \
        RUNE_GUEST_COMPLETED_DIR="${completed}" \
        "${PROJECT_ROOT}/runtime/modules/image/assets/guest/rune-guest-queue"

    [ "$(<"${BATS_TEST_TMPDIR}/result")" = queued ]
    [ "$(<"${completed}/test.job.status")" = 0 ]
    [ -f "${completed}/test.job" ]
}
