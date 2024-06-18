#! /bin/bash

echo "
  ____             __          _       _                
 / ___|    __ _   / _|   ___  | |     (_)  _ __     ___ 
 \___ \   / _\` | | |_   / _ \ | |     | | | '_ \   / _ \\
  ___) | | (_| | |  _| |  __/ | |___  | | | | | | |  __/
 |____/   \__,_| |_|    \___| |_____| |_| |_| |_|  \___|
"

export STREAM=${STREAM:-0}
export CDN=${CDN:-1}

echo $1

qrcode() {
    echo

    echo "█████████████████████████████████████████"
    echo "█████████████████████████████████████████"
    echo "████ ▄▄▄▄▄ █▀ █▀▀██▀▄▀▀▄▀▄▀▄██ ▄▄▄▄▄ ████"
    echo "████ █   █ █▀ ▄ █▀▄▄▀▀ ▄█▄  ▀█ █   █ ████"
    echo "████ █▄▄▄█ █▀█ █▄█▄▀▀▄▀▄ ▀▀▄▄█ █▄▄▄█ ████"
    echo "████▄▄▄▄▄▄▄█▄█▄█ █▄▀ █ ▀▄▀ █▄█▄▄▄▄▄▄▄████"
    echo "████▄ ▄▄ █▄▄  ▄█▄▄▄▄▀▄▀▀▄██ ▄▄▀▄█▄▀ ▀████"
    echo "████▄ ▄▀▄ ▄▀▄ ▀ ▄█▀ ▀▄ █▀▀ ▀█▀▄██▄▀▄█████"
    echo "█████ ▀▄█ ▄ ▄▄▀▄▀▀█▄▀▄▄▀▄▀▄ ▄ ▀▄▄▄█▀▀████"
    echo "████ █▀▄▀ ▄▀▄▄▀█▀ ▄▄ █▄█▀▀▄▀▀█▄█▄█▀▄█████"
    echo "████ █ ▀  ▄▀▀ ██▄█▄▄▄▄▄▀▄▀▀▀▄▄▀█▄▀█ ▀████"
    echo "████ █ ▀▄ ▄██▀▀ ▄█▀ ▀███▄  ▀▄▀▄▄ ▄▀▄█████"
    echo "████▀▄▄█  ▄▀▄▀ ▄▀▀▀▄▀▄▀ ▄▀▄  ▄▀ ▄▀█ ▀████"
    echo "████ █ █ █▄▀ █▄█▀ ▄▄███▀▀▀▄█▀▄ ▀  ▀▄█████"
    echo "████▄███▄█▄▄▀▄ █▄█▄▄▄▄▀▀▄█▀▀ ▄▄▄  ▀█ ████"
    echo "████ ▄▄▄▄▄ █▄▀█ ▄█▀▄ █▀█▄ ▀  █▄█  ▀▄▀████"
    echo "████ █   █ █  █▄▀▀▀▄▄▄▀▀▀▀▀▀ ▄▄  ▀█  ████"
    echo "████ █▄▄▄█ █  ▀█▀ ▄▄▄▄ ▀█ ▀▀▄▀ ▀▀ ▀██████"
    echo "████▄▄▄▄▄▄▄█▄▄██▄█▄▄█▄██▄██▄▄█▄▄█▄█▄█████"
    echo "█████████████████████████████████████████"
    echo "█████████████████████████████████████████"

    echo
    echo "微信扫描上方二维码加入雷池项目讨论组"
}

command_exists() {
    command -v "$1" 2>&1
}

check_container_health() {
    local container_name=$1
    local max_retry=30
    local retry=0
    local health_status="unhealthy"
    echo "Waiting for $container_name to be healthy"
    while [[ "$health_status" == "unhealthy" && $retry -lt $max_retry ]]; do
        health_status=$(docker inspect --format='{{.State.Health.Status}}' $container_name 2>/dev/null || echo 'unhealthy')
        sleep 5
        retry=$((retry+1))
    done
    if [[ "$health_status" == "unhealthy" ]]; then
        abort "Container $container_name is unhealthy"
    fi
    echo "Container $container_name is healthy"
}

space_left() {
    dir="$1"
    while [ ! -d "$dir" ]; do
        dir=$(dirname "$dir")
    done
    echo $(df -h "$dir" --output='avail' | tail -n 1)
}

confirm() {
    echo -e -n "\033[34m[SafeLine] $* \033[1;36m(Y/n)\033[0m"
    read -n 1 -s opt

    [[ "$opt" == $'\n' ]] || echo

    case "$opt" in
    'y' | 'Y') return 0 ;;
    'n' | 'N') return 1 ;;
    *) confirm "$1" ;;
    esac
}

info() {
    echo -e "\033[37m[SafeLine] $*\033[0m"
}

warning() {
    echo -e "\033[33m[SafeLine] $*\033[0m"
}

abort() {
    mv $compose_name.old $compose_name 2>/dev/null || true
    qrcode
    echo -e "\033[31m[SafeLine] $*\033[0m"
    exit 1
}

trap 'onexit' INT
onexit() {
    echo
    abort "用户手动结束升级"
}

get_average_delay() {
    local source=$1
    local total_delay=0
    local iterations=3

    for ((i = 0; i < iterations; i++)); do
        # check timeout
        if ! curl -o /dev/null -m 1 -s -w "%{http_code}\n" "$source" > /dev/null; then
            delay=999
        else
            delay=$(curl -o /dev/null -s -w "%{time_total}\n" "$source")
        fi
        total_delay=$(awk "BEGIN {print $total_delay + $delay}")
    done

    average_delay=$(awk "BEGIN {print $total_delay / $iterations}")
    echo "$average_delay"
}

install_docker() {
    curl -fsSL "https://waf-ce.chaitin.cn/release/latest/get-docker.sh" -o get-docker.sh
    sources=(
        "https://mirrors.aliyun.com/docker-ce"
        "https://mirrors.tencent.com/docker-ce"
        "https://download.docker.com"
    )
    min_delay=${#sources[@]}
    selected_source=""
    for source in "${sources[@]}"; do
        average_delay=$(get_average_delay "$source")
        echo "source: $source, delay: $average_delay"
        if (( $(awk 'BEGIN { print '"$average_delay"' < '"$min_delay"' }') )); then
            min_delay=$average_delay
            selected_source=$source
        fi
    done

    echo "selected source: $selected_source"
    export DOWNLOAD_URL="$selected_source"
    bash get-docker.sh

    docker version > /dev/null 2>&1
    if [ $? -ne "0" ]; then
        echo "Docker 安装失败, 请检查网络连接或手动安装 Docker"
        echo "参考文档: https://docs.docker.com/engine/install/"
        abort "Docker 安装失败"
    fi
    info "Docker 安装成功"
}

start_docker() {
    echo "start docker"
    systemctl enable docker
    systemctl daemon-reload
    systemctl start docker
}

check_depend() {
    # CPU ssse3 指令集检查
    support_ssse3=1
    lscpu | grep ssse3 > /dev/null 2>&1
    if [ $? -ne "0" ]; then
        echo "not found info in lscpu"
        support_ssse3=0
    fi
    cat /proc/cpuinfo | grep ssse3 > /dev/null 2>&1
    if [ $support_ssse3 -eq "0" -a $? -ne "0" ]; then
      abort "雷池需要运行在支持 ssse3 指令集的 CPU 上，虚拟机请自行配置开启 CPU ssse3 指令集支持"
    fi
    if [ -z "$BASH" ]; then
        abort "请用 bash 执行本脚本，请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
    fi

    if [ ! -t 0 ]; then
        abort "STDIN 不是标准的输入设备，请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
    fi

    if [ "$EUID" -ne "0" ]; then
        abort "请以 root 权限运行"
    fi

    if [ -z `command_exists docker` ]; then
        warning "缺少 Docker 环境"
        if confirm "是否需要自动安装 Docker"; then
            install_docker
        else
            abort "中止安装"
        fi
    fi

    info "发现 Docker 环境: '`command -v docker`'"

    docker version > /dev/null 2>&1
    if [ $? -ne "0" ]; then
        abort "Docker 服务工作异常"
    fi

    compose_command="docker compose"
    if $compose_command version; then
        info "发现 Docker Compose Plugin"
    else
        compose_command="docker-compose"
        if $compose_command version; then
            info "发现 Docker Compose"
        else
            warning "未发现 Docker Compose"
            if confirm "是否需要自动安装 Docker Compose"; then
                install_docker
                if [ $? -ne "0" ]; then
                    abort "Docker Compose 安装失败"
                fi
                info "Docker Compose 安装完成"
            else
                abort "中止安装"
            fi
        fi
    fi

    # check docker compose support -d
    if ! $compose_command up -d --help > /dev/null 2>&1; then
        warning "Docker Compose Plugin 不支持 '-d' 参数"
        if confirm "是否需要自动升级 Docker Compose Plugin"; then
            install_docker
            if [ $? -ne "0" ]; then
                abort "Docker Compose Plugin 升级失败"
            fi
            info "Docker Compose Plugin 升级完成"
        else
            abort "中止安装"
        fi
    fi

    start_docker

    info "安装环境确认正常"
}

check_depend

container_id=$(docker ps -n 1 --filter name=.*safeline-mgt.* --format '{{.ID}}')
safeline_path=$(docker inspect --format '{{index .Config.Labels "com.docker.compose.project.working_dir"}}' $container_id)

while [ -z "$safeline_path" ]; do
    echo -e -n "\033[34m[SafeLine] 未发现正在运行的雷池，请输入雷池安装路径 (留空则为 '/data/safeline'): \033[0m"
    read input_path
    [[ -z "$input_path" ]] && input_path='/data/safeline'

    if [[ ! $input_path == /* ]]; then
        warning "'$input_path' 不是合法的绝对路径"
        continue
    fi

    safeline_path=$input_path
done

cd "$safeline_path"

grep COLLIE .env >/dev/null 2>&1
if [ $? -eq "0" ]; then
    abort "检测到你的环境通过牧云主机助手安装，请使用牧云主机助手-应用市场进行升级."
fi

compose_name=$(ls docker-compose.yaml compose.yaml 2>/dev/null)
compose_path=$safeline_path/$compose_name

mv $compose_name $compose_name.old || true

curl "https://waf-ce.chaitin.cn/release/latest/compose.yaml" -sSLk -o $compose_name
curl "https://waf-ce.chaitin.cn/release/latest/reset_tengine.sh" -sSLk -o reset_tengine.sh

if [ $? -ne "0" ]; then
    abort "下载 compose.yaml 脚本失败"
fi
info "下载 compose.yaml 脚本成功"

if [ $STREAM -eq 1 ]; then
    sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest-stream/g" ".env"
else
    sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" ".env"
fi

grep "SAFELINE_DIR" ".env" >/dev/null || echo "SAFELINE_DIR=$(pwd)" >>".env"

if [ $STREAM -eq 1 ]; then
    grep "IMAGE_TAG" ".env" >/dev/null || echo "IMAGE_TAG=latest-stream" >>".env"
else
    grep "IMAGE_TAG" ".env" >/dev/null || echo "IMAGE_TAG=latest" >>".env"
fi

grep "MGT_PORT" ".env" >/dev/null || echo "MGT_PORT=9443" >>".env"
grep "POSTGRES_PASSWORD" ".env" >/dev/null || echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >>".env"
grep "SUBNET_PREFIX" ".env" >/dev/null || echo "SUBNET_PREFIX=172.22.222" >>".env"

if [ -z "$CDN" ]; then
    if [[ $(curl -s ipinfo.io/country) == "CN" ]]; then
        CDN=1
    else
        CDN=0
    fi
fi

if [ $CDN -eq 0 ]; then
    sed -i "s/IMAGE_PREFIX=.*/IMAGE_PREFIX=chaitin/g" ".env"
else
    sed -i "s/IMAGE_PREFIX=.*/IMAGE_PREFIX=swr.cn-east-3.myhuaweicloud.com\/chaitin-safeline/g" ".env"
fi

if [ $CDN -eq 0 ]; then
    grep "IMAGE_PREFIX" ".env" >/dev/null || echo "IMAGE_PREFIX=chaitin" >>".env"
else
    grep "IMAGE_PREFIX" ".env" >/dev/null || echo "IMAGE_PREFIX=swr.cn-east-3.myhuaweicloud.com/chaitin-safeline" >>".env"
fi

info "升级 .env 脚本成功"

info "即将开始下载新版本 Docker 镜像"

$compose_command pull
if [ $? -ne "0" ]; then
    abort "下载新版本 Docker 镜像失败"
fi
info "下载新版本 Docker 镜像成功"

info "即将开始替换 Docker 容器"

# 升级到 3.14.0 版本时，移除了 safeline-redis 容器，需要删除容器，否则无法启动新 compose 网络
docker rm -f safeline-redis &>/dev/null

$compose_command down --remove-orphans && $compose_command up -d
if [ $? -ne "0" ]; then
    abort "替换 Docker 容器失败"
fi

qrcode

check_container_health safeline-pg
check_container_health safeline-mgt
docker exec safeline-mgt /app/mgt-cli reset-admin --once

info "雷池升级成功"

warning "雷池 WAF 社区版安装成功, 请访问以下地址访问控制台"
warning "https://0.0.0.0:9443/"
