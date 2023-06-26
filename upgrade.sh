#! /bin/bash



echo "Upgrade success!"
#!/bin/bash

echo "
  ____             __          _       _                
 / ___|    __ _   / _|   ___  | |     (_)  _ __     ___ 
 \___ \   / _\` | | |_   / _ \ | |     | | | '_ \   / _ \\
  ___) | | (_| | |  _| |  __/ | |___  | | | | | | |  __/
 |____/   \__,_| |_|    \___| |_____| |_| |_| |_|  \___|
"

qrcode() {
    echo
    echo "█████████████████████████████████████████"
    echo "█████████████████████████████████████████"
    echo "████ ▄▄▄▄▄ █▀█ █▄█▀▀ ██▄▀▄▀▄██ ▄▄▄▄▄ ████"
    echo "████ █   █ █▀▀▀█ ▀▄ █▀█▄ ▄  ▀█ █   █ ████"
    echo "████ █▄▄▄█ █▀ █▀▀▀▄▀ ▀ ▄▀▀▀▄▄█ █▄▄▄█ ████"
    echo "████▄▄▄▄▄▄▄█▄▀ ▀▄█ ▀ █ ▀▄▀ █▄█▄▄▄▄▄▄▄████"
    echo "████▄  ▄ ▀▄  ▄▀▄▀▄▀▄▄▄▄▀▀██ ▄▄▀▄█▄▀ ▀████"
    echo "████▄▄▄▄▀█▄▄▄█▄█▀██▄▄▄ ██▀ ▀█▀▄██▄▀▄█████"
    echo "████ ▄█▄▄ ▄▄ ▄▄█▄▀█▄▀▄▀▀▄▀▀ ▄ ▀▄▄▄█▀▀████"
    echo "█████▄████▄█▀ ▄ ▄▄█  █▄██  ▀▀█▄█▄█▀▄█████"
    echo "██████ █▀▄▄█▄▄ ▄▀▀█▄▄▄▀▀▄▀▄▀▄▄▀█▄▀█ ▀████"
    echo "████▀▄██▀ ▄▄  ▀█▀ ▄ █▄▀█ ▀▄▀▄▀▄▄ ▄▀▄█████"
    echo "████ ▄▄█▀ ▄█▀ ██▄█▄▄▄▄▀▀▄▀▀  ▄▀ ▄▀█ ▀████"
    echo "████ █ ██▄▄█▄█▄ ▄█▀ ▀███▄ ██▀▄ ▀  ▀▄█████"
    echo "████▄██▄▄█▄█  ▀▄▀▀▀▄▄▄▄▀▀▀▀▀ ▄▄▄  ▀█ ████"
    echo "████ ▄▄▄▄▄ █▄ ▄█▀ ▄ ▀█▀▀█ ▀  █▄█  ▀ ▀████"
    echo "████ █   █ █ ▀▄█▄█▄▄▀▄▀▀▄▀▀▀ ▄▄  ▀█  ████"
    echo "████ █▄▄▄█ █ █  ▄█▀ ▄█▀█▀ █▀▄▀ ▀▀ ▀██████"
    echo "████▄▄▄▄▄▄▄█▄▄█▄███▄█▄████▄▄▄█▄▄█▄█▄█████"
    echo "█████████████████████████████████████████"
    echo "█████████████████████████████████████████"
    echo
    echo "微信扫描上方二维码加入雷池项目讨论组"
}

command_exists() {
	command -v "$1" 2>&1
}

space_left() {
    dir="$1"
    while [ ! -d "$dir" ]; do
        dir=`dirname "$dir"`;
    done
    echo `df -h "$dir" --output='avail' | tail -n 1`
}

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

info() {
    echo -e "\033[37m[SafeLine] $*\033[0m"
}

warning() {
    echo -e "\033[33m[SafeLine] $*\033[0m"
}

abort() {
    qrcode
    echo -e "\033[31m[SafeLine] $*\033[0m"
    exit 1
}

if [[ "$#" -ne "0" ]]; then
    abort "当前脚本无需任何参数, 直接运行即可"
fi
info "运行参数确认正常"

if [ "$EUID" -ne "0" ]; then
    abort "请以 root 权限运行"
fi
info "运行权限确认正常"

if [ -z `command_exists docker` ]; then
    warning "缺少 Docker 环境"
    if confirm "是否需要自动安装 Docker"; then
        curl -sSLk https://get.docker.com/ | bash
        if [ $? -ne "0" ]; then
            abort "Docker 安装失败"
        fi
        info "Docker 安装完成"
    else
        abort "中止安装"
    fi
fi
info "发现 Docker 环境: '`command -v docker`'"

docker version > /dev/null 2>&1
if [ $? -ne "0" ]; then
    abort "Docker 服务工作异常"
fi
info "Docker 工作状态正常"

compose_plugin=true
compose_command="docker compose"
docker compose version > /dev/null 2>&1 || compose_plugin=false || compose_command="docker-compose"

if [[ "x${compose_plugin}" = "xfalse" ]]; then
    warning "未发现 Docker Compose Plugin"
    if [ -z `command_exists "docker-compose"` ]; then
        warning "未发现 docker-compose 组件"
        if confirm "是否需要自动安装 Docker Compose Plugin"; then
            curl -sSLk https://get.docker.com/ | bash
            if [ $? -ne "0" ]; then
                abort "Docker Compose Plugin 安装失败"
            fi
            info "Docker Compose Plugin 安装完成"
            compose_plugin=true
            compose_command="docker compose"
        else
            abort "中止安装"
        fi
    else
        info "发现 docker-compose 组件: '`command -v docker-compose`'"
    fi
else
    info "发现 Docker Compose Plugin"
fi

compose_path=`$compose_command ls | grep safeline | awk '{print $3}'`

if [[ -z "$compose_path" ]]; then
    abort "没有正在运行的雷池环境"
fi

safeline_path=`dirname $compose_path`
info "发现位于 '$safeline_path' 的雷池环境"

cd "$safeline_path"

mv compose.yaml compose.yaml.old

wget "https://waf-ce.chaitin.cn/release/latest/compose.yaml" --no-check-certificate -O compose.yaml
if [ $? -ne "0" ]; then
    abort "下载 compose.yaml 脚本失败"
fi
info "下载 compose.yaml 脚本成功"

sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" ".env"

grep "SAFELINE_DIR" ".env" > /dev/null || echo "SAFELINE_DIR=$(pwd)" >> ".env"
grep "SUBNET_PREFIX" ".env" > /dev/null || echo "SUBNET_PREFIX=169.254.0" >> ".env"
grep "REDIS_PASSWORD" ".env" > /dev/null || echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> ".env"

info "升级 .env 脚本成功"

info "即将开始下载新版本 Docker 镜像"

$compose_command pull
if [ $? -ne "0" ]; then
    abort "下载新版本 Docker 镜像失败"
fi
info "下载新版本 Docker 镜像成功"

info "即将开始替换 Docker 容器"

$compose_command down && $compose_command up -d
if [ $? -ne "0" ]; then
    abort "替换 Docker 容器失败"
fi
info "雷池升级成功"

qrcode

warning "雷池 WAF 社区版安装成功, 请访问以下地址访问控制台"
warning "https://0.0.0.0:9443/"

