#!/usr/bin/env bash

module::register \
    --name image \
    --description "Build and manage bootable root filesystem images."

command::register --cmd validate --handler image::validate \
    --description "Validate image definitions in rune.yaml."
command::register --cmd build --handler image::build \
    --description "Build a configured Debian image."
command::register --cmd build-all --handler image::build_all \
    --description "Build every image declared in the configuration."
command::register --cmd ensure --handler image::ensure \
    --description "Build a configured image when it is missing."
command::register --cmd ensure-all --handler image::ensure_all \
    --description "Build every configured image that is not cached."
command::register --cmd catalog --handler image::catalog \
    --description "List configured images and their cache state."
command::register --cmd list --handler image::list \
    --description "List built images."
command::register --cmd show --handler image::show \
    --description "Show built image metadata."
command::register --cmd path --handler image::path \
    --description "Print a built image rootfs path."
command::register --cmd delete --handler image::delete \
    --description "Delete a built image."
