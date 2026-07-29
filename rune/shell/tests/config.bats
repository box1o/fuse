#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/lib/init.sh"
    config::reset
    FILE="${BATS_TEST_TMPDIR}/config.yaml"
    cat >"${FILE}" <<'EOF'
version: 1
workers:
  cpp:
    packages:
      - cmake
      - ninja-build
  render:
    packages:
      - blender
EOF
}

@test "config loads and reads a versioned YAML document" {
    config::load rune "${FILE}" 1

    run config::get rune '.workers.cpp.packages[0]'

    [ "${status}" -eq 0 ]
    [ "${output}" = cmake ]
}

@test "config exposes mapping keys and array values" {
    config::load rune "${FILE}" 1

    run config::keys rune '.workers'
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = cpp ]
    [ "${lines[1]}" = render ]

    run config::values rune '.workers.cpp.packages'
    [ "${status}" -eq 0 ]
    [ "${lines[1]}" = ninja-build ]
}

@test "config rejects unsupported document versions" {
    run config::load rune "${FILE}" 2

    [ "${status}" -ne 0 ]
    [[ "${output}" == *"Unsupported configuration version"* ]]
}
