#!/bin/bash

set -e

SCRIPT_DIR=$(dirname "$0")

confirm() {
    echo -e -n "\033[34m[SafeLine] $* \033[1;36m(Y/n)\033[0m"
    read -n 1 -s opt

    [[ "$opt" == $'\n' ]] || echo

    case "$opt" in
        'y' | 'Y' ) return 0;;
        'n' | 'N' ) return 1;;
        *) confirm "$1";;
    esac
}

if ! confirm "是否重新生成 tengine 的所有配置"; then
    exit 0
fi

if [ ! -d "${SCRIPT_DIR}/resources/nginx" ]; then
    echo "website dir not found"
    exit 1
fi

mv "${SCRIPT_DIR}"/resources/nginx "${SCRIPT_DIR}"/resources/nginx."$(date +%s)"

docker restart safeline-tengine > /dev/null

docker exec safeline-mgt gentenginewebsite