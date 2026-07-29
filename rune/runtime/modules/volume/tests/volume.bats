#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
    module::load "${PROJECT_ROOT}/runtime/modules/volume"
    module::initialize volume
    CACHE="${BATS_TEST_TMPDIR}/volumes"
    CONTEXT=${MODULE_CONTEXT_REGISTRY[volume]}
    var::update "${CONTEXT}" cache "${CACHE}"
}

@test "volume creates an isolated copy with metadata" {
    truncate -s 1M "${BATS_TEST_TMPDIR}/source.ext4"

    run volume::create worker-root --from "${BATS_TEST_TMPDIR}/source.ext4" --owner firecracker:worker

    [ "${status}" -eq 0 ]
    [ -f "${CACHE}/worker-root/volume.ext4" ]
    [ "$(yaml::get "${CACHE}/worker-root/metadata.yaml" '.owner')" = firecracker:worker ]
}

@test "volume refuses to delete an attached volume" {
    truncate -s 1M "${BATS_TEST_TMPDIR}/source.ext4"
    volume::create worker-root --from "${BATS_TEST_TMPDIR}/source.ext4" --owner firecracker:worker >/dev/null

    run volume::delete worker-root

    [ "${status}" -ne 0 ]
    [ -d "${CACHE}/worker-root" ]
}

@test "volume cleanup lists unattached volumes without deleting them" {
    truncate -s 1M "${BATS_TEST_TMPDIR}/source.ext4"
    volume::create orphan --from "${BATS_TEST_TMPDIR}/source.ext4" >/dev/null

    run volume::cleanup

    [ "${status}" -eq 0 ]
    [ "${output}" = orphan ]
    [ -d "${CACHE}/orphan" ]
}

@test "volume grows an ext4 filesystem" {
    truncate -s 16M "${BATS_TEST_TMPDIR}/source.ext4"
    mkfs.ext4 -q -F "${BATS_TEST_TMPDIR}/source.ext4"
    volume::create growable --from "${BATS_TEST_TMPDIR}/source.ext4" >/dev/null

    run volume::grow growable --size 32M

    [ "${status}" -eq 0 ]
    [ "$(stat -c '%s' "${CACHE}/growable/volume.ext4")" -eq 33554432 ]
    e2fsck -fn "${CACHE}/growable/volume.ext4"
}
