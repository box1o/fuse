#!/usr/bin/env bash

image::_init()
{
    local -n _image_=${1}

    local __cache__
    __cache__=$(paths::require 'runtime|cache') || return
    var::update _image_ cache "${__cache__}/images"
}
