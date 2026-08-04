#!/usr/bin/env bash

module::register \
    --name {{MODULE_NAME}} \
    --description "{{MODULE_DESCRIPTION}}"

command::register \
    --cmd {{COMMAND_NAME}} \
    --handler {{MODULE_NAMESPACE}}::{{COMMAND_FUNCTION}} \
    --description "{{COMMAND_SUMMARY}}"
