#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    FILE="${BATS_TEST_TMPDIR}/rune.yaml"
    cat >"${FILE}" <<'EOF'
version: 1
kernels:
  test:
    version: "6.1"
    architecture: x86_64
    source: url
    url: https://example.invalid/vmlinux
images:
  worker:
    version: v1
    architecture: x86_64
    size: 2G
    system:
      distribution: debian
      release: bookworm
    kernel: test
EOF
}

@test "project configuration validates across modules" {
    run env RUNE_CONFIG_FILE="${FILE}" "${PROJECT_ROOT}/rune" config validate

    [ "${status}" -eq 0 ]
}

@test "project configuration reads expressions" {
    run env RUNE_CONFIG_FILE="${FILE}" "${PROJECT_ROOT}/rune" config get '.images.worker.version'

    [ "${status}" -eq 0 ]
    [ "${output}" = v1 ]
}

@test "project configuration rejects non-mapping image sections" {
    yq eval -i '.images = []' "${FILE}"

    run env RUNE_CONFIG_FILE="${FILE}" "${PROJECT_ROOT}/rune" config validate

    [ "${status}" -ne 0 ]
    [[ "${output}" == *"must be a mapping"* ]]
}
