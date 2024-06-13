#!/bin/bash

echo "
  ____             __          _       _                
 / ___|    __ _   / _|   ___  | |     (_)  _ __     ___ 
 \___ \   / _\` | | |_   / _ \ | |     | | | '_ \   / _ \\
  ___) | | (_| | |  _| |  __/ | |___  | | | | | | |  __/
 |____/   \__,_| |_|    \___| |_____| |_| |_| |_|  \___|
"

export STREAM=${STREAM:-0}
export CDN=${CDN:-1}

qrcode() {
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

command_exists() {
	command -v "$1" 2>&1
}

check_container_health() {
    local container_name=$1
    local max_retry=30
    local retry=0
    local health_status="unhealthy"
    info "Waiting for $container_name to be healthy"
    while [[ "$health_status" == "unhealthy" && $retry -lt $max_retry ]]; do
        health_status=$(docker inspect --format='{{.State.Health.Status}}' $container_name 2>/dev/null || info 'unhealthy')
        sleep 5
        retry=$((retry+1))
    done
    if [[ "$health_status" == "unhealthy" ]]; then
        abort "Container $container_name is unhealthy"
    fi
    info "Container $container_name is healthy"
}

space_left() {
    dir="$1"
    while [ ! -d "$dir" ]; do
        dir=`dirname "$dir"`;
    done
    echo `df -h "$dir" --output='avail' | tail -n 1`
}

local_ips() {
    if [ -z `command_exists ip` ]; then
        ip_cmd="ip addr show"
    else
        ip_cmd="ifconfig -a"
    fi

    echo $($ip_cmd | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | awk '{print $2}')
}

ips=`local_ips`
subnets="172.22.222 169.254.222 192.168.222"

for subnet in $subnets; do
    if [[ $ips != *$subnet* ]]; then
        SUBNET_PREFIX=$subnet
        break
    fi
done

start_docker() {
    systemctl start docker && systemctl enable docker
}

trap 'onexit' INT
onexit() {
    echo
    abort "用户手动结束安装"
}

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

safeline_path='/data/safeline'

if [ -z "$BASH" ]; then
    abort "请用 bash 执行本脚本，请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ ! -t 0 ]; then
    abort "STDIN 不是标准的输入设备，请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ "$#" -ne "0" ]; then
    abort "当前脚本无需任何参数，请参考最新的官方技术文档 https://waf-ce.chaitin.cn/"
fi

if [ "$EUID" -ne "0" ]; then
    abort "请以 root 权限运行"
fi
info "脚本调用方式确认正常"

install_docker() {
    if [[ $(curl -s ipinfo.io/country) == "CN" ]]; then
        sources=(
            "https://mirrors.aliyun.com/docker-ce"
            "https://mirrors.tencent.com/docker-ce"
            "https://mirrors.163.com/docker-ce"
            "https://mirrors.cernet.edu.cn/docker-ce"
        )

        docker_install_scripts=(
            "https://get.docker.com"
            "https://testingcf.jsdelivr.net/gh/docker/docker-install@master/install.sh"
            "https://cdn.jsdelivr.net/gh/docker/docker-install@master/install.sh"
            "https://fastly.jsdelivr.net/gh/docker/docker-install@master/install.sh"
            "https://gcore.jsdelivr.net/gh/docker/docker-install@master/install.sh"
            "https://raw.githubusercontent.com/docker/docker-install/master/install.sh"
        )

        get_average_delay() {
            local source=$1
            local total_delay=0
            local iterations=3

            for ((i = 0; i < iterations; i++)); do
                delay=$(curl -o /dev/null -s -w "%{time_total}\n" "$source")
                total_delay=$(awk "BEGIN {print $total_delay + $delay}")
            done

            average_delay=$(awk "BEGIN {print $total_delay / $iterations}")
            echo "$average_delay"
        }

        min_delay=${#sources[@]}
        selected_source=""

        for source in "${sources[@]}"; do
            average_delay=$(get_average_delay "$source")

            if (( $(awk 'BEGIN { print '"$average_delay"' < '"$min_delay"' }') )); then
                min_delay=$average_delay
                selected_source=$source
            fi
        done

        if [ -n "$selected_source" ]; then
            echo "选择延迟最低的源 $selected_source，延迟为 $min_delay 秒"
            export DOWNLOAD_URL="$selected_source"
            
            for alt_source in "${docker_install_scripts[@]}"; do
                echo "尝试从备选链接 $alt_source 下载 Docker 安装脚本..."
                if curl -fsSL --retry 2 --retry-delay 3 --connect-timeout 5 --max-time 10 "$alt_source" -o get-docker.sh; then
                    echo "成功从 $alt_source 下载安装脚本"
                    break
                else
                    echo "从 $alt_source 下载安装脚本失败，尝试下一个备选链接"
                fi
            done
            
            if [ ! -f "get-docker.sh" ]; then
                echo "所有下载尝试都失败了。您可以尝试手动安装 Docker，运行以下命令："
                echo "bash <(curl -sSL https://linuxmirrors.cn/docker.sh)"
                exit 1
            fi

            sh get-docker.sh

            echo "启动 docker"
            systemctl enable docker
            systemctl daemon-reload
            systemctl start docker

            docker_config_folder="/etc/docker"
            if [[ ! -d "$docker_config_folder" ]]; then
                mkdir -p "$docker_config_folder"
            fi

            docker version >/dev/null 2>&1
            if [[ $? -ne 0 ]]; then
                echo "Docker 安装失败，您可以尝试手动安装 Docker，运行以下命令："
                echo "bash <(curl -sSL https://linuxmirrors.cn/docker.sh)"
                exit 1
            else
                echo "docker 安装成功"
            fi
        else
            echo "无法选择源进行安装"
            exit 1
        fi
    else
        echo "非中国大陆地区，无需更改源"
        export DOWNLOAD_URL="https://download.docker.com"
        curl -fsSL "https://get.docker.com" -o get-docker.sh
        sh get-docker.sh

        echo "启动 docker"
        systemctl enable docker
        systemctl daemon-reload
        systemctl start docker

        docker_config_folder="/etc/docker"
        if [[ ! -d "$docker_config_folder" ]]; then
            mkdir -p "$docker_config_folder"
        fi

        docker version >/dev/null 2>&1
        if [[ $? -ne 0 ]]; then
            echo "Docker 安装失败，您可以尝试手动安装 Docker，运行以下命令："
            echo "bash <(curl -sSL https://linuxmirrors.cn/docker.sh)"
            exit 1
        else
            echo "docker 安装成功"
        fi
    fi
}
if [ -z `command_exists docker` ]; then
    warning "缺少 Docker 环境"
    if confirm "是否需要自动安装 Docker"; then
        install_docker

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
            install_docker
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
        warning "'$input_path' 路径已经存在，请换一个"
        continue
    fi

    safeline_path=$input_path

    if confirm "目录 '$safeline_path' 当前剩余存储空间为 `space_left \"$safeline_path\"` ，雷池至少需要 5G，是否确定"; then
        break
    fi
done

mkdir -p "$safeline_path"
if [ $? -ne "0" ]; then
    abort "创建安装目录 '$safeline_path' 失败"
fi
info "创建安装目录 '$safeline_path' 成功"
cd "$safeline_path"

curl "https://waf-ce.chaitin.cn/release/latest/compose.yaml" -sSLk -o compose.yaml
curl "https://waf-ce.chaitin.cn/release/latest/reset_tengine.sh" -sSLk -o reset_tengine.sh

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

if [ $STREAM -eq 1 ]; then
    echo "IMAGE_TAG=latest-stream" >>".env"
else
    echo "IMAGE_TAG=latest" >>".env"
fi

echo "MGT_PORT=9443" >> .env
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "SUBNET_PREFIX=$SUBNET_PREFIX" >> .env

if [ $CDN -eq 0 ]; then
    echo "IMAGE_PREFIX=chaitin" >>".env"
else
    echo "IMAGE_PREFIX=swr.cn-east-3.myhuaweicloud.com/chaitin-safeline" >>".env"
fi

info "即将开始下载 Docker 镜像"

$compose_command up -d

if [ $? -ne "0" ]; then
    abort "启动 Docker 容器失败"
fi

qrcode

check_container_health safeline-pg
check_container_health safeline-mgt
docker exec safeline-mgt /app/mgt-cli reset-admin --once

warning "雷池 WAF 社区版安装成功，请访问以下地址访问控制台"
warning "https://0.0.0.0:9443/"
for ip in $ips; do
    warning https://$ip:9443/
done
