#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    RUNE="${PROJECT_ROOT}/rune"
}

@test "firecracker download exposes help" {
    run "${RUNE}" vm download --help

    [ "${status}" -eq 0 ]
    [[ "${output}" == *"Usage:"* ]]
    [[ "${output}" == *"automatic resume"* ]]
}

@test "firecracker download rejects a missing version value" {
    run "${RUNE}" vm download --version

    [ "${status}" -eq 2 ]
    [[ "${output}" == *"requires a value"* ]]
}

@test "firecracker download rejects unknown options" {
    run "${RUNE}" vm download --unknown

    [ "${status}" -eq 2 ]
    [[ "${output}" == *"Unknown option"* ]]
}

@test "firecracker download rejects positional arguments" {
    run "${RUNE}" vm download unexpected

    [ "${status}" -eq 2 ]
    [[ "${output}" == *"unexpected argument"* ]]
}
