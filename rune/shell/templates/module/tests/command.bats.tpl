#!/usr/bin/env bats

@test "{{MODULE_NAME}} {{COMMAND_NAME}} exposes help" {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    run "${PROJECT_ROOT}/rune" {{MODULE_NAME}} {{COMMAND_NAME}} --help
    [ "${status}" -eq 0 ]
    [[ "${output}" == *'Usage:'* ]]
}
