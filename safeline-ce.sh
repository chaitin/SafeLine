#! /bin/bash
set -eE

installer_path=$1

version_file="VERSION.TXT"

if [[ ! -f $version_file ]]; then
    echo "Error: VERSION.TXT not found!"
    exit 1
fi

version=$(cat VERSION.TXT)

if [ -z "$installer_path" ];then
    installer_path="/data/safeline-ce"
fi

if [[ ! -e $installer_path ]]; then
    echo "WAF will be installed at $installer_path, y/N"
    read answer
    if [ "$answer" != "${answer#[Yy]}" ] ; then
        echo "Start installing..."
    else
        echo "End"
        exit 1
    fi
elif [[ ! -d $installer_path ]]; then
    echo "Error: $installer_path already exists but is not a directory"
    exit 1
fi

env_file=".env"
if [[ ! -f $env_file ]]; then
    echo -n "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)
HOST_RESOURCES_DIR=$installer_path/resources
HOST_LOGS_DIR=$installer_path/logs
IMAGE_TAG=$version
COMPOSE_PROJECT_NAME=safeline-ce
COMPOSE_FILE=compose.yaml" > $env_file
fi

mkdir -p $installer_path
