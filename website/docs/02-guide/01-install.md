---
title: "安装雷池"
---

# 安装雷池

## 配置需求

- 操作系统：Linux
- 指令架构：x86_64
- 软件依赖：Docker 20.10.6 版本以上
- 软件依赖：Docker Compose 2.0.0 版本以上
- 最小化环境：1 核 CPU / 1 GB 内存 / 5 GB 磁盘

可以逐行执行以下命令来确认服务器配置

```shell
uname -m                 # 查看指令架构
docker version           # 查看 Docker 版本
docker compose version   # 查看 Docker Compose 版本
docker-compose version   # 同上（兼容老版本 Docker Compose）
cat /proc/cpuinfo        # 查看 CPU 信息
cat /proc/meminfo        # 查看内存信息
df -h                    # 查看磁盘信息

lscpu | grep ssse3       # 确认 CPU 是否支持 ssse3 指令集
```

有三种安装方式供选择

- [在线安装](#在线安装) : 推荐安装方式
- [离线安装](#离线安装) : 服务器无法连接 Docker Hub 时选择
- [一键安装](#使用牧云助手安装) : 最简单的安装方式

## 在线安装

**_如果服务器可以访问互联网环境，推荐使用该方式_**

执行以下命令，即可开始安装

```
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh)"
```

> 如果连接 Docker Hub 网络不稳，导致镜像下载失败，可以采用 [离线安装](#离线安装) 方式

经过以上步骤，你的雷池已经安装好了，下一步请参考 [登录雷池](/docs/guide/login)

## 离线安装

如果你的服务器无法连接互联网环境，或连接 Docker Hub 网络不稳，可以使用镜像包安装方式

> 这里忽略 Docker 安装的过程

首先，下载 [雷池社区版镜像包](https://demo.waf-ce.chaitin.cn/image.tar.gz) 并传输到需要安装雷池的服务器上，执行以下命令加载镜像

```
cat image.tar.gz | gzip -d | docker load
```

执行以下命令创建并进入雷池安装目录

```
mkdir -p safeline        # 创建 safeline 目录
cd safeline              # 进入 safeline 目录
```

下载 [编排脚本](https://waf-ce.chaitin.cn/release/latest/compose.yaml) 和 [seccomp](https://waf-ce.chaitin.cn/release/latest/seccomp.json) 并传输到 safeline 目录中

执行以下命令，生成雷池运行所需的相关环境变量

```
echo "SAFELINE_DIR=$(pwd)" >> .env
echo "IMAGE_TAG=latest" >> .env
echo "MGT_PORT=9443" >> .env
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "SUBNET_PREFIX=172.22.222" >> .env
```

执行以下命令启动雷池

```
docker compose up -d
```

经过以上步骤，你的雷池已经安装好了，下一步请参考 [登录雷池](/docs/guide/login)

## 使用牧云助手安装

也可以使用 [牧云主机管理助手](https://collie.chaitin.cn/) 进行一键安装

![](/images/docs/guide_install/collie_apps.png)

参考视频教程 [用 “白嫖的云主机” 一键安装 “开源的 Web 防火墙”](https://www.bilibili.com/video/BV1sh4y1t7Pk/)

## 常见安装问题

请参考 [安装问题](/docs/faq/install)
