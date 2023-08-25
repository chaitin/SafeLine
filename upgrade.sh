#! /bin/bash

echo "
  ____             __          _       _                
 / ___|    __ _   / _|   ___  | |     (_)  _ __     ___ 
 \___ \   / _\` | | |_   / _ \ | |     | | | '_ \   / _ \\
  ___) | | (_| | |  _| |  __/ | |___  | | | | | | |  __/
 |____/   \__,_| |_|    \___| |_____| |_| |_| |_|  \___|
"

echo $1

qrcode() {
    echo

    echo "█████████████████████████████████████████"
    echo "█████████████████████████████████████████"
    echo "████ ▄▄▄▄▄ █▀ █▀▀█▄▀█▀▀▄▀▄▀▄██ ▄▄▄▄▄ ████"
    echo "████ █   █ █▀ ▄ █▀█ ▄▀ ▄▀▄  ▀█ █   █ ████"
    echo "████ █▄▄▄█ █▀█ █▄█ ▀ █▀▄▀▀▀▄▄█ █▄▄▄█ ████"
    echo "████▄▄▄▄▄▄▄█▄█▄█ █ ▀▄█ █▄▀ █▄█▄▄▄▄▄▄▄████"
    echo "████▄ ▄▄ ▀▄   ▄█▄▄▀▄▄▄▀▀▀██ ▄▄▀▄█▄▀ ▀████"
    echo "████▄▄ ▄▀▀▄█  ▀ ▄█▀ ▀▄▀█▄▀ ▀█▀▄██▄▀▄█████"
    echo "████  █▀  ▄ ▀ ▀▄▀▀█▄▀▄▀▀▄▀▄ ▄ ▀▄▄▄█▀▀████"
    echo "████ ▄  █ ▄█▀▄██▀ ▄  █▄▀▀ ▄▀▀█▄█▄█▀▄█████"
    echo "████▄ ▀▄█▄▄█ █▄█▄█▄▄▄▄▄▀▄▀▀▀▄▄▀█▄▀█ ▀████"
    echo "████▄▄▄▀█▄▄█▄▀▀ ▄█▀▄▄▄▄▀█  ▀▄▀▄▄ ▄▀▄█████"
    echo "████▀█▀▀ ▄▄▀▄▄ ▄▀▀▀▄▀▄▀▀▄▀▀  ▄▀ ▄▀█ ▀████"
    echo "████ █▀█ ▀▄▀█▀██▀ ▄  █▄██  █▀▄ ▀  ▀▄█████"
    echo "████▄██▄█▄▄█▀▄▄█▄█▄▄▄▄▀ ▄▀▀▀ ▄▄▄  ▀█ ████"
    echo "████ ▄▄▄▄▄ █▄▀█ ▄█▀ ██▀█▀▀█  █▄█  ▀▄▀████"
    echo "████ █   █ █ ▀ ▄▀▀▀█▄▄▄▀▀▀▄▀ ▄▄  ▀█  ████"
    echo "████ █▄▄▄█ █ ▀▄█▀ ▄  ███  ▄▀▄▀ ▀▀ ▀██████"
    echo "████▄▄▄▄▄▄▄█▄███▄█▄▄█▄████▄▄▄█▄▄█▄█▄█████"
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

trap 'onexit' INT
onexit() {
    echo
    abort "用户手动结束升级"
}

# CPU ssse3 指令集检查
lscpu | grep ssse3 > /dev/null 2>&1
if [ $? -ne "0" ]; then
    abort "雷池需要运行在支持 ssse3 指令集的 CPU 上，虚拟机请自行配置开启 CPU ssse3 指令集支持"
fi

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

container_id=`docker ps --filter ancestor=chaitin/safeline-mgt-api --format '{{.ID}}'`
mount_path=`docker inspect --format '{{range .Mounts}}{{if eq .Destination "/logs"}}{{.Source}}{{end}}{{end}}' $container_id`
safeline_path=`dirname $mount_path`

while [ -z "$safeline_path" ]; do
    echo -e -n "\033[34m[SafeLine] 未发现正在运行的雷池，请输入雷池安装路径 (留空则为 '`pwd`'): \033[0m"
    read input_path
    [[ -z "$input_path" ]] && input_path=`pwd`

    if [[ ! $input_path == /* ]]; then
        warning "'$input_path' 不是合法的绝对路径"
        continue
    fi

    safeline_path=$input_path
done

cd "$safeline_path"

compose_name=`ls docker-compose.yaml compose.yaml 2>/dev/null`
compose_path=$safeline_path/$compose_name

if [ -f "$compose_path" ]; then
    info "发现位于 '$safeline_path' 的雷池环境"
else
    abort "没有发现位于 $safeline_path 的雷池环境"
fi


mv $compose_name $compose_name.old

curl "https://waf-ce.chaitin.cn/release/latest/compose.yaml" -sSLk -o $compose_name
if [ $? -ne "0" ]; then
    abort "下载 compose.yaml 脚本失败"
fi
info "下载 compose.yaml 脚本成功"

curl -sS "https://waf-ce.chaitin.cn/release/latest/seccomp.json" -o seccomp.json
if [ $? -ne "0" ]; then
    abort "下载 seccomp.json 脚本失败"
fi
info "下载 seccomp.json 脚本成功"

sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" ".env"

grep "SAFELINE_DIR" ".env" > /dev/null || echo "SAFELINE_DIR=$(pwd)" >> ".env"
grep "IMAGE_TAG" ".env" > /dev/null || echo "IMAGE_TAG=latest" >> ".env"
grep "MGT_PORT" ".env" > /dev/null || echo "MGT_PORT=9443" >> ".env"
grep "POSTGRES_PASSWORD" ".env" > /dev/null || echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> ".env"
grep "REDIS_PASSWORD" ".env" > /dev/null || echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> ".env"
grep "SUBNET_PREFIX" ".env" > /dev/null || echo "SUBNET_PREFIX=172.22.222" >> ".env"

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

