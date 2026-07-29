#!/usr/bin/env bats
load runtime_test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    modules_root="${BATS_TEST_TMPDIR}/modules"
    mkdir -p "${modules_root}"
}

@test 'scaffold creates expected files with requested command' {
    run scaffold::module sample --command greet --summary 'Sample module' --directory "${modules_root}"
    [ "${status}" -eq 0 ]
    [ -f "${modules_root}/sample/module.sh" ]
    [ -f "${modules_root}/sample/vars.sh" ]
    [ -f "${modules_root}/sample/greet/command.sh" ]
    [ -f "${modules_root}/sample/tests/greet.bats" ]
    bash -n "${modules_root}/sample/module.sh"
}

@test 'scaffold refuses existing nonempty destination' {
    mkdir -p "${modules_root}/sample"
    touch "${modules_root}/sample/existing"
    run scaffold::module sample --force --directory "${modules_root}"
    [ "${status}" -ne 0 ]
}

@test 'scaffold rejects invalid and traversal names' {
    run scaffold::module '../bad' --directory "${modules_root}"
    [ "${status}" -eq 2 ]
}

@test 'scaffold dry-run writes nothing' {
    run scaffold::module sample --dry-run --directory "${modules_root}"
    [ "${status}" -eq 0 ]
    [ ! -e "${modules_root}/sample" ]
    [[ "${output}" == *'/sample/module.sh'* ]]
}

@test 'generated module loads successfully' {
    scaffold::module sample --directory "${modules_root}"
    command::_reset
    module::_reset
    run module::load "${modules_root}/sample"
    [ "${status}" -eq 0 ]
}
