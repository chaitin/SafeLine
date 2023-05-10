#! /bin/bash
set -eE

abort()
{
    echo $1
    exit 1
}

echo "
  ____             __          _       _                
 / ___|    __ _   / _|   ___  | |     (_)  _ __     ___ 
 \___ \   / _\` | | |_   / _ \ | |     | | | '_ \   / _ \\
  ___) | | (_| | |  _| |  __/ | |___  | | | | | | |  __/
 |____/   \__,_| |_|    \___| |_____| |_| |_| |_|  \___|
"

if [[ "$#" -ne "0" ]]; then
    echo "Usage: run "$0" to set up Safeline CE at current working directory."
    exit 0
fi

command -v docker > /dev/null || abort "docker not found, unable to deploy"
compose_plugin=true
compose_command="docker compose"
docker --help | grep compose | grep v2 > /dev/null || compose_plugin=false || compose_command="docker-compose"

if [[ "x${compose_plugin}" = "xfalse" ]]; then
    command -v docker-compose > /dev/null && docker-compose --version | grep v2 > /dev/null || abort "docker compose v2 not found, unable to deploy"
fi

COMPOSE_YAML="compose.yaml"
wget https://waf-ce.chaitin.cn/release/latest/compose.yaml --no-check-certificate -O ${COMPOSE_YAML}

ENV_FILE=".env"
if [[ ! -f ${ENV_FILE} ]]; then
    echo "SAFELINE_DIR=$(pwd)" >> .env
    echo "IMAGE_TAG=latest" >> .env
    echo "MGT_PORT=9443" >> .env
    echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
    echo "SUBNET_PREFIX=169.254.0" >> .env
fi

echo "Setup success!"
echo "Run '$compose_command up -d' to start SafeLine."
echo "And then visit https://<SafeLine-IP>:9443."
