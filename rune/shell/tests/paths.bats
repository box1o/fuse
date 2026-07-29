#!/usr/bin/env bats
load test_helper

@test "register and get path" {
    paths::register --source test 'test|output' "${BATS_TEST_TMPDIR}/out"
    run paths::get 'test|output'
    [ "$status" -eq 0 ]
    [ "$output" = "${BATS_TEST_TMPDIR}/out" ]
}

@test "duplicate path registration fails" {
    paths::register 'test|duplicate' "${BATS_TEST_TMPDIR}/one"
    run paths::register 'test|duplicate' "${BATS_TEST_TMPDIR}/two"
    [ "$status" -ne 0 ]
}

@test "default paths load when HOME is unset" {
    run env -u HOME bash -uo pipefail -c \
        'source "$1" && paths::require "user|home"' \
        bash "${PROJECT_ROOT}/lib/init.sh"

    [ "$status" -eq 0 ]
    [[ "$output" == /* ]]
}
