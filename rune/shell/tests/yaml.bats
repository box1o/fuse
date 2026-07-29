#!/usr/bin/env bats
load test_helper

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/.." && pwd -P)"
    source "${PROJECT_ROOT}/lib/init.sh"
    command -v yq >/dev/null 2>&1 || skip "yq v4 is not installed"
    file="${BATS_TEST_TMPDIR}/config.yaml"
    cat >"${file}" <<'YAML'
server:
  host: localhost
  port: 8080
features:
  enabled: true
YAML
}

@test "validate YAML" {
    run yaml::validate "${file}"
    [ "$status" -eq 0 ]
}

@test "get YAML value" {
    run yaml::get "${file}" '.server.port'
    [ "$status" -eq 0 ]
    [ "$output" = 8080 ]
}

@test "set raw YAML value" {
    yaml::set_raw "${file}" '.server.port' '9090'
    run yaml::get "${file}" '.server.port'
    [ "$output" = 9090 ]
}

@test "require YAML values" {
    run yaml::require "${file}" '.server.host' '.server.port'
    [ "$status" -eq 0 ]
}
