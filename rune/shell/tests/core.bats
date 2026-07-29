#!/usr/bin/env bats
load test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/.." && pwd -P)"
    source "${PROJECT_ROOT}/lib/init.sh"
}

@test "temp resources are created under the requested base and cleaned" {
    base="${BATS_TEST_TMPDIR}/temporary base"
    directory="$(temp::dir work "${base}")"
    file="$(temp::file work "${base}")"
    temp::register "${directory}"
    temp::register "${file}"

    [ -d "${directory}" ]
    [ -f "${file}" ]

    temp::cleanup

    [ ! -e "${directory}" ]
    [ ! -e "${file}" ]
}

@test "lock with preserves command failure and releases the descriptor" {
    fail_operation() {
        return 17
    }

    lock_file="${BATS_TEST_TMPDIR}/operation.lock"
    status=0
    lock::with "${lock_file}" fail_operation || status=$?

    [ "${status}" -eq 17 ]
    [[ ! -v "LOCK_FDS[${lock_file}]" ]]
}

@test "trap handlers are functions and run without eval" {
    marker="${BATS_TEST_TMPDIR}/handled"
    test_cleanup() {
        touch "${marker}"
    }

    trap::add USR1 test_cleanup
    trap::_dispatch USR1
    trap - USR1

    [ -f "${marker}" ]
}

@test "environment loading restores the allexport option" {
    environment_file="${BATS_TEST_TMPDIR}/environment"
    printf 'LOADED_VALUE=ready\n' >"${environment_file}"
    set +a

    env::load_file "${environment_file}"

    [[ $- != *a* ]]
    [ "${LOADED_VALUE}" = ready ]
    export -p | grep -q 'LOADED_VALUE'
}

@test "archive verification rejects parent traversal" {
    source_directory="${BATS_TEST_TMPDIR}/source"
    archive_file="${BATS_TEST_TMPDIR}/unsafe.tar"
    mkdir -p "${source_directory}"
    printf data >"${source_directory}/file"
    tar -cf "${archive_file}" --transform='s|file|../file|' -C "${source_directory}" file

    run archive::verify "${archive_file}"

    [ "${status}" -ne 0 ]
    [[ "${output}" == *"Unsafe archive entry"* ]]
}

@test "chroot validation resolves absolute shell symlinks inside the root" {
    root="${BATS_TEST_TMPDIR}/root"
    mkdir -p "${root}/bin"
    printf '#!/usr/bin/env bash\n' >"${root}/bin/busybox"
    chmod +x "${root}/bin/busybox"
    ln -s /bin/busybox "${root}/bin/sh"

    run chroot::validate "${root}"

    [ "${status}" -eq 0 ]
}
