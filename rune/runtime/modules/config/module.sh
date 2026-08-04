#!/usr/bin/env bash

module::register \
    --name config \
    --requires image,kernel \
    --description "Inspect and validate the Rune project configuration."

command::register --cmd validate --handler project_config::validate \
    --description "Validate the complete project configuration."
command::register --cmd get --handler project_config::get \
    --description "Read a value from the project configuration."
command::register --cmd show --handler project_config::show \
    --description "Print the project configuration."
command::register --cmd path --handler project_config::path \
    --description "Print the active project configuration path."
