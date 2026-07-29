#!/usr/bin/env bash

image::_configure_root()
{
    local __file__=${1:-}
    local __expression__=${2:-}
    local __name__=${3:-}
    local __root__=${4:-}

    printf '%s\n' "${__name__}" | fs::atomic_write "${__root__}/etc/hostname"
    fs::atomic_write "${__root__}/etc/fstab" <<'EOF'
/dev/vda / ext4 defaults 0 1
proc /proc proc defaults 0 0
sysfs /sys sysfs defaults 0 0
devtmpfs /dev devtmpfs mode=0755,nosuid 0 0
EOF
    fs::mkdir "${__root__}/etc/systemd/system/getty.target.wants" || return
    fs::symlink /lib/systemd/system/serial-getty@.service \
        "${__root__}/etc/systemd/system/getty.target.wants/serial-getty@ttyS0.service" || return
    image::_install_guest_runtime "${__root__}" || return

    local __environment__
    __environment__=$(yq eval -r "${__expression__}.environment // {} | to_entries | .[] | .key + \"=\" + (.value | tostring)" "${__file__}") || return
    if [[ -n "${__environment__}" ]]; then
        while IFS= read -r __line__; do
            [[ "${__line__}" =~ ^[A-Z_][A-Z0-9_]*=.*$ ]] || {
                log::error "Invalid image environment entry: ${__line__}"
                return 2
            }
        done <<<"${__environment__}"
        printf '%s\n' "${__environment__}" | fs::atomic_write "${__root__}/etc/environment"
    fi

    local __config_dir__
    __config_dir__=$(path::dirname "${__file__}") || return
    local __count__
    __count__=$(yaml::length "${__file__}" "${__expression__}.files // []") || return
    local __index__
    for ((__index__ = 0; __index__ < __count__; __index__++)); do
        local __source__
        local __destination__
        local __mode__
        __source__=$(yaml::get "${__file__}" "${__expression__}.files[${__index__}].source") || return
        __destination__=$(yaml::get "${__file__}" "${__expression__}.files[${__index__}].destination") || return
        __mode__=$(yaml::get "${__file__}" "${__expression__}.files[${__index__}].mode" '') || return
        [[ "${__destination__}" == /* ]] || {
            log::error "Image file destination must be absolute: ${__destination__}"
            return 2
        }
        __source__=$(path::absolute "${__config_dir__}/${__source__}") || return
        path::assert_within "${__source__}" "${__config_dir__}" || return
        fs::copy_file "${__source__}" "${__root__}${__destination__}" || return
        [[ -z "${__mode__}" ]] || fs::chmod "${__mode__}" "${__root__}${__destination__}" || return
    done
}

image::_install_guest_runtime()
{
    local __root__=${1:-}
    local __module_root__
    local __project_root__
    __module_root__=$(module::root image) || return
    __project_root__=$(paths::require 'runtime|project_root') || return

    fs::mkdir "${__root__}/usr/local/lib/rune/shell" || return
    fs::copy_tree "${__project_root__}/shell/lib" "${__root__}/usr/local/lib/rune/shell/lib" || return
    fs::copy_file \
        "${__module_root__}/assets/guest/rune-guest" \
        "${__root__}/usr/local/bin/rune-guest" || return
    fs::copy_file \
        "${__module_root__}/assets/guest/rune-guest-queue" \
        "${__root__}/usr/local/bin/rune-guest-queue" || return
    fs::copy_file \
        "${__module_root__}/assets/guest/rune-guest-queue.service" \
        "${__root__}/etc/systemd/system/rune-guest-queue.service" || return
    fs::chmod 755 "${__root__}/usr/local/bin/rune-guest" || return
    fs::chmod 755 "${__root__}/usr/local/bin/rune-guest-queue" || return
    fs::mkdir "${__root__}/etc/systemd/system/multi-user.target.wants" || return
    fs::symlink /etc/systemd/system/rune-guest-queue.service \
        "${__root__}/etc/systemd/system/multi-user.target.wants/rune-guest-queue.service" || return
    fs::mkdir "${__root__}/var/log/rune/jobs" || return
    fs::mkdir "${__root__}/var/lib/rune/queue" || return
    fs::mkdir "${__root__}/var/lib/rune/running" || return
    fs::mkdir "${__root__}/var/lib/rune/completed" || return
    fs::mkdir "${__root__}/workspace" || return
}
