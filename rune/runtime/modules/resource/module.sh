#!/usr/bin/env bash

module::register \
    --name resource \
    --description "Validate and inspect compute resources."

command::register --cmd host --handler resource::host \
    --description "Show host CPU and memory capacity."
command::register --cmd usage --handler resource::usage \
    --description "Show CPU and memory usage for a process."
command::register --cmd validate --handler resource::validate \
    --description "Validate CPU and memory values."
