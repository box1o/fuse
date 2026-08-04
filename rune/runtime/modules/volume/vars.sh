#!/usr/bin/env bash

volume::_init()
{
    local -n _volume_=${1}
    local __cache__
    __cache__=$(paths::require 'runtime|cache') || return
    var::update _volume_ cache "${VOLUME_CACHE_DIR:-${__cache__}/volumes}"
}
