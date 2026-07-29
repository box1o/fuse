#!/usr/bin/env bats
load runtime_test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
    testmod::run() {
        printf '%s\n' "$@"
    }
    module::register --name testmod --summary test --root "${BATS_TEST_TMPDIR}"
}

@test 'register and list a command' {
    command::register --module testmod --name run --summary test --handler testmod::run
    run command::list
    [ "${status}" -eq 0 ]
    [ "${output}" = 'testmod run' ]
}

@test 'register infers the module and defaults command text' {
    MODULE_LOADING_NAME=testmod
    command::register --cmd run --handler testmod::run
    MODULE_LOADING_NAME=

    [ "${COMMAND_SUMMARY_REGISTRY["testmod run"]}" = run ]
    [ "${COMMAND_USAGE_REGISTRY["testmod run"]}" = 'rune testmod run [options]' ]
}

@test 'reject duplicate command' {
    command::register --module testmod --name run --summary test --handler testmod::run
    run command::register --module testmod --name run --summary test --handler testmod::run
    [ "${status}" -ne 0 ]
}

@test 'reject missing handler' {
    run command::register --module testmod --name run --summary test --handler testmod::missing
    [ "${status}" -ne 0 ]
}

@test 'dispatch forwards arguments' {
    command::register --module testmod --name run --summary test --handler testmod::run
    run command::dispatch testmod run 'one two' three
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = 'one two' ]
    [ "${lines[1]}" = three ]
}

@test 'dispatch preserves handler status' {
    testmod::fail() {
        return 42
    }
    command::register --module testmod --name fail --summary test --handler testmod::fail
    run command::dispatch testmod fail
    [ "${status}" -eq 42 ]
}

@test 'dispatch handles command help without calling the handler' {
    testmod::fail_if_called() {
        return 99
    }
    command::register --module testmod --cmd help-test --handler testmod::fail_if_called --description 'Generated help'

    run command::dispatch testmod help-test --help

    [ "${status}" -eq 0 ]
    [[ "${output}" == *'rune testmod help-test [options]'* ]]
    [[ "${output}" == *'Generated help'* ]]
}

@test 'unknown command returns usage status' {
    run command::dispatch testmod missing
    [ "${status}" -eq 2 ]
}
