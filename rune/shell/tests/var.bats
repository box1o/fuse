#!/usr/bin/env bats
load test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/.." && pwd -P)"
    source "${PROJECT_ROOT}/lib/init.sh"
}

@test "var updates an associative array through a nameref" {
    local -A values=()
    local -n _values_=values

    var::update _values_ source_dir "/project/source"

    [ "${values[source_dir]}" = "/project/source" ]
    [ "$(var::get values source_dir)" = "/project/source" ]
}

@test "var get supports a default" {
    local -A values=()

    [ "$(var::get values missing fallback)" = fallback ]
}
