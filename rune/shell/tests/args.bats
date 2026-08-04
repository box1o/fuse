#!/usr/bin/env bats
load test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/.." && pwd -P)"
    source "${PROJECT_ROOT}/lib/init.sh"
    args::reset
}

@test "args parses flags, options, defaults, and positionals" {
    args::flag --name verbose --short v
    args::option --name output --short o
    args::option --name format --default text

    args::parse -v --output "file name" first -- second

    args::has verbose
    [ "$(args::get verbose)" = true ]
    [ "$(args::get output)" = "file name" ]
    [ "$(args::get format)" = text ]
    [ "$(args::count)" -eq 2 ]
    [ "$(args::position 0)" = first ]
    [ "$(args::position 1)" = second ]
}

@test "args supports long option equals syntax" {
    args::option --name output

    args::parse --output=value

    [ "$(args::get output)" = value ]
}

@test "args rejects unknown options" {
    run args::parse --unknown

    [ "${status}" -eq 2 ]
    [[ "${output}" == *"Unknown option"* ]]
}

@test "args enforces required options" {
    args::option --name output --required

    run args::parse

    [ "${status}" -eq 2 ]
    [[ "${output}" == *"Required option is missing"* ]]
}

@test "args rejects duplicate definitions" {
    args::flag --name verbose --short v

    run args::option --name verbose

    [ "${status}" -ne 0 ]
}
