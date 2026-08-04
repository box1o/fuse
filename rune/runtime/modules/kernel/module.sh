#!/usr/bin/env bash

module::register \
    --name kernel \
    --description "Build and manage guest kernels."

command::register --cmd validate --handler kernel::validate \
    --description "Validate kernel definitions in rune.yaml."
command::register --cmd ensure --handler kernel::ensure \
    --description "Download a configured kernel when it is missing."
command::register --cmd list --handler kernel::list \
    --description "List installed kernels."
command::register --cmd show --handler kernel::show \
    --description "Show installed kernel metadata."
command::register --cmd delete --handler kernel::delete \
    --description "Delete an installed kernel."
