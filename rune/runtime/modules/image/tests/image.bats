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
images:
  cpp-worker:
    version: v1
    architecture: x86_64
    size: 2G
    system:
      distribution: debian
      release: bookworm
    kernel: test-kernel
    packages:
      - build-essential
      - cmake
    capabilities:
      - ci
      - cpp
  render-worker:
    version: v2
    architecture: x86_64
    size: 4G
    system:
      distribution: debian
      release: bookworm
    kernel: test-kernel
    packages:
      - blender
EOF
}

@test "image catalog lists every configured image" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" image catalog

    [ "${status}" -eq 0 ]
    [[ "${output}" == *"cpp-worker"* ]]
    [[ "${output}" == *"render-worker"* ]]
    [[ "${output}" == *"missing"* ]]
}

@test "image build-all dry-run processes every configured image" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" image build-all --dry-run

    [ "${status}" -eq 0 ]
    [[ "${output}" == *"Image: cpp-worker:v1"* ]]
    [[ "${output}" == *"Image: render-worker:v2"* ]]
}

@test "image configuration validates" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" image validate cpp-worker
    [ "${status}" -eq 0 ]
}

@test "image dry-run includes baseline and requested packages" {
    run env RUNE_CONFIG_FILE="${CONFIG}" "${PROJECT_ROOT}/rune" image build cpp-worker --dry-run
    [ "${status}" -eq 0 ]
    [[ "${output}" == *"busybox-static"* ]]
    [[ "${output}" == *"build-essential"* ]]
    [[ "${output}" == *"Kernel: test-kernel"* ]]
}
