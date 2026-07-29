#!/usr/bin/env bash

module::register \
    --name volume \
    --description "Create and manage file-backed volumes."

command::register --cmd list --handler volume::list \
    --description "List managed volumes."
command::register --cmd show --handler volume::show \
    --description "Show volume metadata."
command::register --cmd create --handler volume::create \
    --description "Create a volume from an image file."
command::register --cmd clone --handler volume::clone \
    --description "Clone an existing volume."
command::register --cmd path --handler volume::path \
    --description "Print a volume path."
command::register --cmd check --handler volume::check \
    --description "Check an ext4 volume."
command::register --cmd grow --handler volume::grow \
    --description "Grow an ext4 volume."
command::register --cmd delete --handler volume::delete \
    --description "Delete an unattached volume."
command::register --cmd cleanup --handler volume::cleanup \
    --description "Find or delete unattached volumes."
