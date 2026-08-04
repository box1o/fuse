#!/usr/bin/env bats
load test_helper

@test "copy file with spaces" {
    source_file="${BATS_TEST_TMPDIR}/source file"
    destination_file="${BATS_TEST_TMPDIR}/target file"
    printf 'hello\n' >"${source_file}"
    run fs::copy_file "${source_file}" "${destination_file}"
    [ "$status" -eq 0 ]
    [ "$(cat "${destination_file}")" = hello ]
}

@test "remove_tree refuses root" {
    run fs::remove_tree /
    [ "$status" -ne 0 ]
}

@test "atomic write replaces destination" {
    destination="${BATS_TEST_TMPDIR}/config.txt"
    fs::atomic_write "${destination}" <<<'value'
    [ "$(cat "${destination}")" = value ]
}
