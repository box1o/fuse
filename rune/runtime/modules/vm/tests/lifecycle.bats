#!/usr/bin/env bats

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../../../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
    module::load "${PROJECT_ROOT}/runtime/modules/resource"
    module::load "${PROJECT_ROOT}/runtime/modules/volume"
    module::load "${PROJECT_ROOT}/runtime/modules/image"
    module::load "${PROJECT_ROOT}/runtime/modules/kernel"
    module::load "${PROJECT_ROOT}/runtime/modules/vm"
    module::initialize vm

    CACHE="${BATS_TEST_TMPDIR}/firecracker"
    CONTEXT=${MODULE_CONTEXT_REGISTRY[vm]}
    var::update "${CONTEXT}" cache "${CACHE}"
    var::update "${CONTEXT}" downloads "${CACHE}/downloads"
    var::update "${CONTEXT}" runtime "${CACHE}/runtime"
    var::update "${CONTEXT}" assets "${CACHE}/assets"
    var::update "${CONTEXT}" vms "${CACHE}/vms"
    var::update "${CONTEXT}" kvm "${CACHE}/kvm"
    VOLUME_CONTEXT=${MODULE_CONTEXT_REGISTRY[volume]}
    var::update "${VOLUME_CONTEXT}" cache "${CACHE}/volumes"
}

@test "module registers the complete VM lifecycle" {
    run command::list_module vm

    [ "${status}" -eq 0 ]
    [ "${#lines[@]}" -eq 20 ]
    [[ "${output}" == *create* ]]
    [[ "${output}" == *chroot* ]]
    [[ "${output}" == *run* ]]
}

@test "job status emits structured completion data" {
    mkdir -p "${CACHE}/vms/builder"
    cat >"${CACHE}/vms/builder/console.log" <<'EOF'
RUNE_JOB_BEGIN compile 2026-01-01T00:00:00Z
output
RUNE_JOB_END compile 7 2026-01-01T00:00:01Z
EOF
    firecracker::_is_running() { return 0; }
    runtime::config::set 'output|format' json

    run firecracker::job_status builder compile

    [ "${status}" -eq 0 ]
    [ "$(yq eval -r '.state' - <<<"${output}")" = completed ]
    [ "$(yq eval -r '.exit_code' - <<<"${output}")" -eq 7 ]
}

@test "wait job returns the guest exit code" {
    mkdir -p "${CACHE}/vms/builder"
    firecracker::_is_running() { return 0; }
    firecracker::_job_status_console() { printf '9\n'; }

    run firecracker::wait_job builder compile --timeout 1

    [ "${status}" -eq 9 ]
    [[ "${output}" == *'exit_code=9'* ]]
}

@test "enqueue writes a safely quoted job into a stopped VM" {
    mkdir -p "${CACHE}/vms/builder"
    truncate -s 16M "${CACHE}/vms/builder/rootfs.ext4"
    mkfs.ext4 -q -F "${CACHE}/vms/builder/rootfs.ext4"
    debugfs -w -R 'mkdir /var' "${CACHE}/vms/builder/rootfs.ext4" >/dev/null 2>&1
    debugfs -w -R 'mkdir /var/lib' "${CACHE}/vms/builder/rootfs.ext4" >/dev/null 2>&1
    debugfs -w -R 'mkdir /var/lib/rune' "${CACHE}/vms/builder/rootfs.ext4" >/dev/null 2>&1
    debugfs -w -R 'mkdir /var/lib/rune/queue' "${CACHE}/vms/builder/rootfs.ext4" >/dev/null 2>&1

    run firecracker::enqueue builder --id compile -- bash -lc 'printf "hello world\n"'

    [ "${status}" -eq 0 ]
    run debugfs -R 'cat /var/lib/rune/queue/compile.job' "${CACHE}/vms/builder/rootfs.ext4"
    [ "${status}" -eq 0 ]
    [[ "${output}" == *'rune-guest run --id compile'* ]]
    [[ "${output}" == *'hello\ world'* ]]
}

@test "job logs extracts guest output from the serial console" {
    mkdir -p "${CACHE}/vms/builder"
    cat >"${CACHE}/vms/builder/console.log" <<'EOF'
booting
RUNE_JOB_BEGIN compile 2026-01-01T00:00:00Z
RUNE_JOB_LOG compile first line
RUNE_JOB_LOG compile second line
RUNE_JOB_END compile 0 2026-01-01T00:00:01Z
login:
EOF

    run firecracker::job_logs builder compile --source console

    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = "first line" ]
    [ "${lines[1]}" = "second line" ]
}

@test "list shows VMs in stable order with their state" {
    mkdir -p "${CACHE}/vms/zebra" "${CACHE}/vms/alpha"

    run firecracker::list

    [ "${status}" -eq 0 ]
    [[ "${lines[0]}" == NAME* ]]
    [[ "${lines[1]}" == alpha*stopped* ]]
    [[ "${lines[2]}" == zebra*stopped* ]]
}

@test "list emits machine-readable JSON" {
    mkdir -p "${CACHE}/vms/builder"
    runtime::config::set 'output|format' json

    run firecracker::list

    [ "${status}" -eq 0 ]
    [ "$(yq eval -r '.[0].name' - <<<"${output}")" = builder ]
    [ "$(yq eval -r '.[0].status' - <<<"${output}")" = stopped ]
}

@test "setup installs release binaries without downloading existing assets" {
    mkdir -p "${CACHE}/downloads" "${CACHE}/assets" "${BATS_TEST_TMPDIR}/release/release-v1-x86_64"
    printf '#!/usr/bin/env bash\n' >"${BATS_TEST_TMPDIR}/release/release-v1-x86_64/firecracker-v1-x86_64"
    printf '#!/usr/bin/env bash\n' >"${BATS_TEST_TMPDIR}/release/release-v1-x86_64/jailer-v1-x86_64"
    tar -czf "${CACHE}/downloads/firecracker-v1-x86_64.tgz" -C "${BATS_TEST_TMPDIR}/release" .
    touch "${CACHE}/assets/vmlinux.bin" "${CACHE}/assets/rootfs.ext4"

    run firecracker::setup --force

    [ "${status}" -eq 0 ]
    [ -x "${CACHE}/runtime/firecracker" ]
    [ -x "${CACHE}/runtime/jailer" ]
}

@test "create writes an isolated rootfs and machine configuration" {
    mkdir -p "${CACHE}/assets"
    printf kernel >"${CACHE}/assets/vmlinux.bin"
    printf rootfs >"${CACHE}/assets/rootfs.ext4"
    firecracker::setup() {
        :
    }

    run firecracker::create builder --cpu 4 --memory 2048 \
        --kernel "${CACHE}/assets/vmlinux.bin" \
        --rootfs "${CACHE}/assets/rootfs.ext4"

    [ "${status}" -eq 0 ]
    [ -f "${CACHE}/volumes/vm-builder-root/volume.ext4" ]
    [ -f "${CACHE}/vms/builder/metadata.yaml" ]
    [ -f "${CACHE}/vms/builder/config.json" ]
    grep -q '"vcpu_count": 4' "${CACHE}/vms/builder/config.json"
    grep -q '"mem_size_mib": 2048' "${CACHE}/vms/builder/config.json"
}

@test "resources shows VM CPU memory and disk" {
    mkdir -p "${CACHE}/vms/builder"
    truncate -s 1M "${CACHE}/vms/builder/rootfs.ext4"
    cat >"${CACHE}/vms/builder/config.json" <<'EOF'
{"machine-config":{"vcpu_count":2,"mem_size_mib":1024}}
EOF

    run firecracker::resources builder

    [ "${status}" -eq 0 ]
    [[ "${output}" == *"cpu: 2"* ]]
    [[ "${output}" == *"memory_mib: 1024"* ]]
}

@test "resize updates CPU and memory for a stopped VM" {
    mkdir -p "${CACHE}/vms/builder"
    cat >"${CACHE}/vms/builder/config.json" <<'EOF'
{"machine-config":{"vcpu_count":2,"mem_size_mib":1024}}
EOF

    run firecracker::resize builder --cpu 4 --memory 2048

    [ "${status}" -eq 0 ]
    [ "$(yaml::get "${CACHE}/vms/builder/config.json" '.machine-config.vcpu_count')" -eq 4 ]
    [ "$(yaml::get "${CACHE}/vms/builder/config.json" '.machine-config.mem_size_mib')" -eq 2048 ]
}

@test "combined resize preserves resource options when volume delegates" {
    mkdir -p "${CACHE}/vms/builder"
    cat >"${CACHE}/vms/builder/config.json" <<'EOF'
{"machine-config":{"vcpu_count":2,"mem_size_mib":1024}}
EOF
    cat >"${CACHE}/vms/builder/metadata.yaml" <<'EOF'
schema: 1
name: builder
provider: firecracker
kernel_path: KERNEL_PATH
root_volume: vm-builder-root
cpu: 2
memory_mib: 1024
EOF
    sed -i "s|KERNEL_PATH|${CACHE}/kernel|" "${CACHE}/vms/builder/metadata.yaml"
    touch "${CACHE}/kernel" "${CACHE}/rootfs.ext4"
    volume::grow() { printf '%s\n' "$*" >"${BATS_TEST_TMPDIR}/grow-call"; }
    volume::path() { printf '%s\n' "${CACHE}/rootfs.ext4"; }

    run firecracker::resize builder --cpu 4 --memory 2048 --disk 4G

    [ "${status}" -eq 0 ]
    [ "$(yaml::get "${CACHE}/vms/builder/config.json" '.machine-config.vcpu_count')" -eq 4 ]
    [ "$(yaml::get "${CACHE}/vms/builder/config.json" '.machine-config.mem_size_mib')" -eq 2048 ]
    [ "$(<"${BATS_TEST_TMPDIR}/grow-call")" = "vm-builder-root --size 4G" ]
}

@test "cleanup removes stale runtime files but preserves logs" {
    mkdir -p "${CACHE}/vms/builder"
    touch "${CACHE}/vms/builder/firecracker.pid" "${CACHE}/vms/builder/firecracker.sock"
    printf 'keep\n' >"${CACHE}/vms/builder/console.log"

    firecracker::cleanup builder

    [ ! -e "${CACHE}/vms/builder/firecracker.pid" ]
    [ ! -e "${CACHE}/vms/builder/firecracker.sock" ]
    [ -s "${CACHE}/vms/builder/console.log" ]
}

@test "status reports a created VM as stopped" {
    mkdir -p "${CACHE}/vms/builder"

    run firecracker::status builder

    [ "${status}" -eq 0 ]
    [ "${output}" = "builder stopped" ]
}

@test "start refuses an inaccessible KVM device" {
    mkdir -p "${CACHE}/vms/builder" "${CACHE}/runtime"
    touch "${CACHE}/vms/builder/config.json"
    printf '#!/usr/bin/env bash\n' >"${CACHE}/runtime/firecracker"
    chmod +x "${CACHE}/runtime/firecracker"

    run firecracker::start builder

    [ "${status}" -ne 0 ]
    [[ "${output}" == *"KVM is not accessible"* ]]
}

@test "delete removes a stopped VM" {
    mkdir -p "${CACHE}/vms/builder"
    touch "${CACHE}/vms/builder/rootfs.ext4"

    firecracker::delete builder

    [ ! -e "${CACHE}/vms/builder" ]
}

@test "run composes setup create and start" {
    calls="${BATS_TEST_TMPDIR}/calls"
    firecracker::setup() {
        printf 'setup\n' >>"${calls}"
    }
    firecracker::create() {
        printf 'create %s\n' "$*" >>"${calls}"
        mkdir -p "${CACHE}/vms/$1"
    }
    firecracker::start() {
        printf 'start %s\n' "$*" >>"${calls}"
    }

    firecracker::run builder --image custom:v1 --cpu 2 --memory 1024

    [ "$(sed -n '1p' "${calls}")" = setup ]
    [ "$(sed -n '2p' "${calls}")" = "create builder --image custom:v1 --cpu 2 --memory 1024" ]
    [ "$(sed -n '3p' "${calls}")" = "start builder" ]
}
