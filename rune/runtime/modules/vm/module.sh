#!/usr/bin/env bash

module::register \
    --name vm \
    --requires resource,volume,image,kernel \
    --init-handler firecracker::_init \
    --description "Create and manage lightweight Firecracker virtual machines."

command::register --cmd download --handler firecracker::download \
    --description "Download a Firecracker release with automatic resume."
command::register --cmd setup --handler firecracker::setup \
    --description "Install the Firecracker and jailer binaries."
command::register --cmd create --handler firecracker::create \
    --description "Create a new virtual machine."
command::register --cmd start --handler firecracker::start \
    --description "Start a virtual machine."
command::register --cmd status --handler firecracker::status \
    --description "Show virtual machine status."
command::register --cmd list --handler firecracker::list \
    --description "List local virtual machines."
command::register --cmd stop --handler firecracker::stop \
    --description "Stop a virtual machine."
command::register --cmd delete --handler firecracker::delete \
    --description "Delete a stopped virtual machine."
command::register --cmd logs --handler firecracker::logs \
    --description "Read a virtual machine console log."
command::register --cmd job-logs --handler firecracker::job_logs \
    --description "Read logs produced by rune-guest inside a virtual machine."
command::register --cmd job-status --handler firecracker::job_status \
    --description "Show the state and exit code of a guest job."
command::register --cmd wait-job --handler firecracker::wait_job \
    --description "Wait for a guest job to complete."
command::register --cmd enqueue --handler firecracker::enqueue \
    --description "Queue a command for execution on the next virtual machine boot."
command::register --cmd chroot --handler firecracker::chroot \
    --description "Mount a stopped VM rootfs and enter it with chroot."
command::register --cmd run --handler firecracker::run \
    --description "Set up Firecracker, create a VM when missing, and start it."
command::register --cmd launch --handler firecracker::run \
    --description "Create a VM from a cached image when needed and start it."
command::register --cmd resources --handler firecracker::resources \
    --description "Show virtual machine CPU, memory, and disk resources."
command::register --cmd resize --handler firecracker::resize \
    --description "Resize a stopped virtual machine."
command::register --cmd set-memory --handler firecracker::set_memory \
    --description "Set memory for a stopped virtual machine."
command::register --cmd cleanup --handler firecracker::cleanup \
    --description "Remove stale virtual machine runtime files."
