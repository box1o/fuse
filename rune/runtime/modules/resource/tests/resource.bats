#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
    module::load "${PROJECT_ROOT}/runtime/modules/resource"
    module::initialize resource
}

@test "resource validates CPU and memory" {
    run resource::validate --cpu 2 --memory 1024

    [ "${status}" -eq 0 ]
}

@test "resource rejects invalid memory" {
    run resource::validate --cpu 2 --memory 0

    [ "${status}" -ne 0 ]
}

@test "resource reports host capacity" {
    run resource::host

    [ "${status}" -eq 0 ]
    [[ "${output}" == *"cpu:"* ]]
    [[ "${output}" == *"memory_mib:"* ]]
}
