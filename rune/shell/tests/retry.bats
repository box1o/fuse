#!/usr/bin/env bats
load test_helper

@test "retry succeeds after transient failures" {
    counter="${BATS_TEST_TMPDIR}/counter"
    printf '0' >"${counter}"
    transient() {
        local n
        n=$(cat "${counter}")
        n=$((n + 1))
        printf '%s' "$n" >"${counter}"
        ((n >= 3))
    }
    run retry::run --attempts 3 --delay 1 -- transient
    [ "$status" -eq 0 ]
}
