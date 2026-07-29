#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    CONFIG="${BATS_TEST_TMPDIR}/rune.yaml"
    cat >"${CONFIG}" <<'EOF'
version: 1
kernels:
  test-kernel:
    version: "6.1"
    architecture: x86_64
    source: url
    url: https://example.invalid/vmlinux
images: {}
EOF
}

@test "kernel configuration validates" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" kernel validate test-kernel
    [ "${status}" -eq 0 ]
}

@test "kernel dry-run reports source and destination without downloading" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" kernel ensure test-kernel --dry-run
    [ "${status}" -eq 0 ]
    [[ "${output}" == *"https://example.invalid/vmlinux"* ]]
    [ ! -e "${BATS_TEST_TMPDIR}/vmlinux" ]
}
