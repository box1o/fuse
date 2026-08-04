#!/usr/bin/env bash

{{MODULE_NAMESPACE}}::{{COMMAND_FUNCTION}}()
{
    if (( $# > 0 )); then
        log::error "{{MODULE_NAME}} {{COMMAND_NAME}}: unexpected argument: $1"
        return 2
    fi

    local __cache__
    __cache__=$(module::var {{MODULE_NAME}} cache) || return

    log::info "{{MODULE_SUMMARY}}: command not yet implemented" "cache=${__cache__}"
}
