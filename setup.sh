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
    echo "████ ▄▄▄▄▄ █▀█ █▄▄▄▄▄▀█▄█▄▀▄██ ▄▄▄▄▄ ████"
    echo "████ █   █ █▀▀▀█ ▄▄▀ ▀█▄ ▄  ▀█ █   █ ████"
    echo "████ █▄▄▄█ █▀ █▀▀▄▄ █▀ ▄▀▀▀▄▄█ █▄▄▄█ ████"
    echo "████▄▄▄▄▄▄▄█▄▀ ▀▄█▄▀▄█▄▀ ▀ █▄█▄▄▄▄▄▄▄████"
    echo "████   ▄ █▄ ▄▄▀▄▀▀█ ▄ ▄▀▄██ ▄▄▀▄█▄▀ ▀████"
    echo "█████ ▄▀█▄▄ ██▄█▀ ▄ ▄█ █▄▀ ▀█▀▄██▄▀▄█████"
    echo "████▄▄▄▀█▀▄▀▄ ▄█▄█▄▄▀ ▄▀▀▀▀ ▄ ▀▄▄▄█▀▀████"
    echo "████▄▄██ ▀▄ ▀ ▄ ▄█▀ ▄█▄██▀█▀▀█▄█▄█▀▄█████"
    echo "████▀▄ ▀▄█▄▄▄▀█▄▀▀▀▄▄▄▀ ▄▀▀▀▄▄▀█▄▀█ ▀████"
    echo "█████ ▄▀ █▄█▀▀██▀ ▄  ███▀▀█▀▄▀▄▄ ▄▀▄█████"
    echo "████ ▄█▀█ ▄███ █▄█▄▄▄▄▀▀▀▀▀  ▄▀ ▄▀█ ▀████"
    echo "████ ██▀▀▄▄▄ ▄█ ▄█▀  ▄ ▀▀▀██▀▄ ▀  ▀▄█████"
    echo "████▄███▄█▄▄  ▄▄▀▀▀▄▀▄▄▀▀▀▀▀ ▄▄▄  ▀█ ████"
    echo "████ ▄▄▄▄▄ █▄▀ █▀ ▄▄▀▄▀██    █▄█  ▀ █████"
    echo "████ █   █ █ ███▄█▄▄▀▄▄▀▄▀▀▀ ▄▄  ▀█▀ ████"
    echo "████ █▄▄▄█ █ ▄▄ ▄█▀ ▄█▄▀▄ █▀▄▀ ▀▀ ▀██████"
    echo "████▄▄▄▄▄▄▄█▄▄▄▄████▄▄▄███▄▄▄█▄▄█▄█▄█████"
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

start_docker() {
    systemctl start docker && systemctl enable docker
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

trap 'onexit' INT
onexit() {
    echo
    abort "用户手动结束安装"
}


safeline_path='/data/safeline'

if [ -z "$BASH" ]; then
    abort "请用 bash 执行本脚本, 请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ ! -t 0 ]; then
    abort "STDIN 不是标准的输入设备, 请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ "$#" -ne "0" ]; then
    abort "当前脚本无需任何参数, 请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ "$EUID" -ne "0" ]; then
    abort "请以 root 权限运行"
fi
info "脚本调用方式确认正常"

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

start_docker
docker version > /dev/null 2>&1
if [ $? -ne "0" ]; then
    abort "Docker 服务工作异常"
fi
info "Docker 工作状态正常"

compose_command="docker compose"
if $compose_command version; then
    info "发现 Docker Compose Plugin"
else
    warning "未发现 Docker Compose Plugin"
    compose_command="docker-compose"
    if [ -z `command_exists "docker-compose"` ]; then
        warning "未发现 docker-compose 组件"
        if confirm "是否需要自动安装 Docker Compose Plugin"; then
            curl -sSLk https://get.docker.com/ | bash
            if [ $? -ne "0" ]; then
                abort "Docker Compose Plugin 安装失败"
            fi
            info "Docker Compose Plugin 安装完成"
            compose_command="docker compose"
        else
            abort "中止安装"
        fi
    else
        info "发现 docker-compose 组件: '`command -v docker-compose`'"
    fi
fi

while true; do
    echo -e -n "\033[34m[SafeLine] 雷池安装目录 (留空则为 '$safeline_path'): \033[0m"
    read input_path
    [[ -z "$input_path" ]] && input_path=$safeline_path

    if [[ ! $input_path == /* ]]; then
        warning "'$input_path' 不是合法的绝对路径"
        continue
    fi

    if [ -f "$input_path" ] || [ -d "$input_path" ]; then
        warning "'$input_path' 路径已经存在, 请换一个"
        continue
    fi

    safeline_path=$input_path

    if confirm "目录 '$safeline_path' 当前剩余存储空间为 `space_left \"$safeline_path\"` , 雷池至少需要 5G, 是否确定"; then
        break
    fi
done

mkdir -p "$safeline_path"
if [ $? -ne "0" ]; then
    abort "创建安装目录 '$safeline_path' 失败"
fi
info "创建安装目录 '$safeline_path' 成功"
cd "$safeline_path"

wget "https://waf-ce.chaitin.cn/release/latest/compose.yaml" --no-check-certificate -O compose.yaml
if [ $? -ne "0" ]; then
    abort "下载 compose.yaml 脚本失败"
fi
info "下载 compose.yaml 脚本成功"

touch ".env"
if [ $? -ne "0" ]; then
    abort "创建 .env 脚本失败"
fi
info "创建 .env 脚本成功"

echo "SAFELINE_DIR=$safeline_path" >> .env
echo "IMAGE_TAG=latest" >> .env
echo "MGT_PORT=9443" >> .env
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "SUBNET_PREFIX=169.254.0" >> .env

info "即将开始下载 Docker 镜像"

$compose_command up -d

if [ $? -ne "0" ]; then
    abort "启动 Docker 容器失败"
fi

qrcode

warning "雷池 WAF 社区版安装成功, 请访问以下地址访问控制台"
warning "https://0.0.0.0:9443/"

