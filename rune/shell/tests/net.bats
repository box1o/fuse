#!/usr/bin/env bats
load test_helper

@test "net::is_url accepts HTTP URLs" {
    run net::is_url 'https://example.com/file'
    [ "$status" -eq 0 ]
}

@test "net::is_url rejects local paths" {
    run net::is_url '/tmp/file'
    [ "$status" -ne 0 ]
}

@test "checksum mismatch never publishes completed download" {
    curl() {
        local output=''
        while [[ $# -gt 0 ]]; do
            case "$1" in --output)
                output=$2
                shift 2
                ;;
            *) shift ;; esac
        done
        printf 'corrupt' >"${output}"
    }
    destination="${BATS_TEST_TMPDIR}/artifact"
    run net::download --sha256 deadbeef https://example.com/artifact "${destination}"
    [ "${status}" -ne 0 ]
    [ ! -e "${destination}" ]
    [ -f "${destination}.part" ]
}
