#!/usr/bin/env bash

scaffold::_render()
{
    local __template__=${1:-} __destination__=${2:-} __module__=${3:-} __command__=${4:-}
    local __summary__=${5:-} __description__=${6:-} __cli__=${7:-}
    fs::require_file template "${__template__}" || return
    local __content__ __upper__ __namespace__
    __content__=$(<"${__template__}") || return
    __upper__=${__module__^^}
    __upper__=${__upper__//-/_}
    __namespace__=${__module__//-/_}
    __content__=${__content__//\{\{MODULE_NAME\}\}/${__module__}}
    __content__=${__content__//\{\{MODULE_UPPER\}\}/${__upper__}}
    __content__=${__content__//\{\{MODULE_NAMESPACE\}\}/${__namespace__}}
    __content__=${__content__//\{\{MODULE_SUMMARY\}\}/${__summary__}}
    __content__=${__content__//\{\{MODULE_DESCRIPTION\}\}/${__description__}}
    __content__=${__content__//\{\{COMMAND_NAME\}\}/${__command__}}
    __content__=${__content__//\{\{COMMAND_FUNCTION\}\}/${__command__//-/_}}
    __content__=${__content__//\{\{COMMAND_FILE\}\}/${__command__}.sh}
    __content__=${__content__//\{\{COMMAND_SUMMARY\}\}/${__summary__}}
    __content__=${__content__//\{\{CLI_NAME\}\}/${__cli__}}
    printf '%s\n' "${__content__}" | fs::atomic_write "${__destination__}"
}

scaffold::module()
{
    local __command__=run __summary__='' __description__='' __parent__='' __force__=false __dry_run__=false __tests__=true
    local __module__=''
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --command)
                [[ $# -ge 2 ]] || return 2
                __command__=$2
                shift 2
                ;;
            --summary)
                [[ $# -ge 2 ]] || return 2
                __summary__=$2
                shift 2
                ;;
            --description)
                [[ $# -ge 2 ]] || return 2
                __description__=$2
                shift 2
                ;;
            --directory)
                [[ $# -ge 2 ]] || return 2
                __parent__=$2
                shift 2
                ;;
            --force)
                __force__=true
                shift
                ;;
            --dry-run)
                __dry_run__=true
                shift
                ;;
            --no-tests)
                __tests__=false
                shift
                ;;
            -*)
                log::error "module scaffold: unknown option: $1"
                return 2
                ;;
            *)
                [[ -z "${__module__}" ]] || {
                    log::error "module scaffold: unexpected argument: $1"
                    return 2
                }
                __module__=$1
                shift
                ;;
        esac
    done
    module::_valid_name "${__module__}" || return
    command::_valid_identifier "${__command__}" || return
    module::_is_reserved "${__module__}" && {
        log::error "Reserved module name: ${__module__}"
        return 1
    }
    [[ -n "${__summary__}" ]] || __summary__="${__module__} module"
    [[ -n "${__description__}" ]] || __description__=${__summary__}
    [[ -n "${__parent__}" ]] || __parent__=$(paths::require 'runtime|modules') || return
    __parent__=$(path::absolute "${__parent__}") || return
    local __destination__
    __destination__=$(path::absolute "${__parent__}/${__module__}") || return
    path::assert_within "${__destination__}" "${__parent__}" || return
    if [[ -e "${__destination__}" ]]; then
        [[ "${__force__}" == true && -d "${__destination__}" ]] || {
            log::error "Destination already exists: ${__destination__}"
            return 1
        }
        [[ -z "$(find "${__destination__}" -mindepth 1 -maxdepth 1 -print -quit)" ]] || {
            log::error "--force only accepts an empty destination: ${__destination__}"
            return 1
        }
    fi
    local __cli__ __templates__
    __cli__=$(runtime::config::get 'cli|name') || return
    __templates__=$(paths::require 'runtime|templates') || return
    local -a __files__=(module.sh vars.sh "${__command__}/command.sh" README.md)
    [[ "${__tests__}" == true ]] && __files__+=("tests/${__command__}.bats")
    if [[ "${__dry_run__}" == true ]]; then
        local __file__
        for __file__ in "${__files__[@]}"; do printf '%s\n' "${__destination__}/${__file__}"; done
        return 0
    fi
    fs::mkdir "${__destination__}" || return
    local __file__ __template__
    for __file__ in "${__files__[@]}"; do
        case "${__file__}" in
            module.sh) __template__="${__templates__}/module.sh.tpl" ;;
            vars.sh) __template__="${__templates__}/vars.sh.tpl" ;;
            "${__command__}/command.sh") __template__="${__templates__}/command.sh.tpl" ;;
            README.md) __template__="${__templates__}/README.md.tpl" ;;
            tests/*) __template__="${__templates__}/tests/command.bats.tpl" ;;
        esac
        scaffold::_render "${__template__}" "${__destination__}/${__file__}" "${__module__}" "${__command__}" "${__summary__}" "${__description__}" "${__cli__}" || return
    done
    while IFS= read -r -d '' __file__; do bash -n "${__file__}" || return; done < <(find "${__destination__}" -type f -name '*.sh' -print0)
    printf 'Created module: %s\nPath: %s\n\nNext:\n  1. Edit %s/vars.sh\n  2. Implement %s::%s\n  3. Run: %s module validate %s\n' \
        "${__module__}" "${__destination__}" "${__destination__}" "${__module__//-/_}" "${__command__//-/_}" "${__cli__}" "${__module__}"
}
