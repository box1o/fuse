#!/usr/bin/env bash

{{MODULE_NAMESPACE}}::_init()
{
    local -n _{{MODULE_NAMESPACE}}_=${1}

    local __cli__
    local __cache__
    __cli__=$(runtime::config::get 'cli|name') || return
    __cache__="${XDG_CACHE_HOME:-${HOME}/.cache}/${__cli__}/{{MODULE_NAME}}"

    var::update _{{MODULE_NAMESPACE}}_ cache "${__cache__}"
}
