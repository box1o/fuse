#!/usr/bin/env bash

if [[ ${BASHLIB_INITIALIZED:-0} == 1 ]]; then
    return 0
fi

BASHLIB_INITIALIZED=1

BASHLIB_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
readonly BASHLIB_DIR

bashlib::_source()
{
    local __relative_path__=${1:-}

    # shellcheck disable=SC1090
    source "${BASHLIB_DIR}/${__relative_path__}"
}

bashlib::_source log.sh
bashlib::_source validate.sh
bashlib::_source var.sh
bashlib::_source args.sh
bashlib::_source proc.sh
bashlib::_source path.sh
bashlib::_source paths/registry.sh
bashlib::_source paths/defaults.sh
bashlib::_source fs.sh
bashlib::_source temp.sh
bashlib::_source trap.sh
bashlib::_source lock.sh
bashlib::_source retry.sh
bashlib::_source env.sh
bashlib::_source user.sh
bashlib::_source system.sh
bashlib::_source archive.sh
bashlib::_source net.sh
bashlib::_source yaml.sh
bashlib::_source config.sh
bashlib::_source mount.sh
bashlib::_source chroot.sh

paths::register_defaults
