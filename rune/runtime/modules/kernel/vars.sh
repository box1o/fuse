#!/usr/bin/env bash

kernel::_init()
{
    local -n _kernel_=${1}

    local __cache__
    __cache__=$(paths::require 'runtime|cache') || return
    var::update _kernel_ cache "${__cache__}/kernels"
}
