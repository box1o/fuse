#!/usr/bin/env bats
load runtime_test_helper

make_module() {
    local name=$1 root=$2
    mkdir -p "${root}/${name}/run"
    cat >"${root}/${name}/run/command.sh" <<EOF
${name}::run() { :; }
EOF
    cat >"${root}/${name}/module.sh" <<EOF
module::register --name '${name}' --description '${name}'
command::register --cmd run --handler '${name}::run'
EOF
}

setup() {
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/../.." && pwd -P)"
    source "${PROJECT_ROOT}/shell/runtime/init.sh"
    command::_reset
    module::_reset
}

@test 'discover loads modules in sorted order' {
    root="${BATS_TEST_TMPDIR}/modules"
    make_module zebra "${root}"
    make_module alpha "${root}"
    module::discover "${root}"
    run module::list
    [ "${status}" -eq 0 ]
    [ "${lines[0]}" = alpha ]
    [ "${lines[1]}" = zebra ]
}

@test 'discovery rejects missing module declaration' {
    mkdir -p "${BATS_TEST_TMPDIR}/modules/broken"
    run module::discover "${BATS_TEST_TMPDIR}/modules"
    [ "${status}" -ne 0 ]
}

@test 'failed load rolls back partial registrations' {
    root="${BATS_TEST_TMPDIR}/modules"
    make_module broken "${root}"
    printf '\nreturn 9\n' >>"${root}/broken/module.sh"
    run module::load "${root}/broken"
    [ "${status}" -eq 9 ]
    ! module::has broken
    ! command::has 'broken run'
}

@test 'module initialization runs once' {
    count_file="${BATS_TEST_TMPDIR}/count"
    printf 0 >"${count_file}"
    once::init() {
        local n
        n=$(<"${count_file}")
        printf '%s' "$((n + 1))" >"${count_file}"
    }
    module::register --name once --summary once --root "${BATS_TEST_TMPDIR}" --init-handler once::init
    module::initialize once
    module::initialize once
    [ "$(<"${count_file}")" = 1 ]
}

@test 'module init receives variables populated through a nameref' {
    root="${BATS_TEST_TMPDIR}/modules"
    mkdir -p "${root}/context"
    cat >"${root}/context/vars.sh" <<'EOF'
context::_init()
{
    local -n _context_=${1}
    var::update _context_ greeting hello
}
EOF
    cat >"${root}/context/module.sh" <<'EOF'
module::register --name context
EOF

    module::load "${root}/context"
    module::initialize context

    [ "$(module::var context greeting)" = hello ]
    [ "$(module::var context root)" = "${root}/context" ]
}

@test 'reserved external module names are rejected' {
    run module::register --name runtime --summary bad --root "${BATS_TEST_TMPDIR}"
    [ "${status}" -ne 0 ]
}

@test 'module initialization initializes dependencies first' {
    order="${BATS_TEST_TMPDIR}/order"
    dependency::_init() { printf 'dependency\n' >>"${order}"; }
    consumer::_init() { printf 'consumer\n' >>"${order}"; }

    module::register --name dependency --root "${BATS_TEST_TMPDIR}"
    module::register --name consumer --root "${BATS_TEST_TMPDIR}" --requires dependency
    module::initialize consumer

    [ "$(sed -n '1p' "${order}")" = dependency ]
    [ "$(sed -n '2p' "${order}")" = consumer ]
}

@test 'module initialization rejects missing dependencies' {
    module::register --name consumer --root "${BATS_TEST_TMPDIR}" --requires missing

    run module::initialize consumer

    [ "${status}" -ne 0 ]
    [[ "${output}" == *'requires missing module'* ]]
}

@test 'module initialization detects dependency cycles' {
    module::register --name alpha --root "${BATS_TEST_TMPDIR}" --requires beta
    module::register --name beta --root "${BATS_TEST_TMPDIR}" --requires alpha

    run module::initialize alpha

    [ "${status}" -ne 0 ]
    [[ "${output}" == *'dependency cycle'* ]]
}
