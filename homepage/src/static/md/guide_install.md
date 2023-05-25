---
title: "安装雷池"
group: "上手指南"
order: 2
---

# 安装雷池

## 配置需求

- 操作系统：Linux
- 指令架构：x86_64
- 软件依赖：Docker 20.10.6 版本以上
- 软件依赖：Docker Compose 2.0.0 版本以上
- 最小化环境：1 核 CPU / 1 GB 内存 / 10 GB 磁盘

可以逐行执行以下命令来确认服务器配置

```shell
uname -m                 # 查看指令架构
docker version           # 查看 Docker 版本
docker compose version   # 查看 Docker Compose 版本
docker-compose version   # 同上（兼容老版本 Docker Compose）
cat /proc/cpuinfo        # 查看 CPU 信息
cat /proc/meminfo        # 查看内存信息
df -h                    # 查看磁盘信息
```

有三种安装方式供选择

- [在线安装](#在线安装) : 推荐安装方式
- [离线安装](#离线安装) : 服务器无法联网时选择
- [一键安装](#使用牧云助手安装) : 最简单的安装方式

## 在线安装

***如果服务器可以访问互联网环境，推荐使用在线安装方式***

### 在线安装 Docker

执行以下命令来安装 Docker 和 Docker Compose

```bash
curl -fsSLk https://get.docker.com/ | bash
```

安装命令结束后，可以执行以下命令来确认 Docker 和 Docker Compose 是否安装成功

```
docker version           # 查看 Docker 版本
docker compose version   # 查看 Docker Compose 版本
```

### 在线安装雷池

执行以下命令创建并进入雷池安装目录

```
mkdir -p safeline        # 创建 safeline 目录
cd safeline              # 进入 safeline 目录
```

执行以下命令，将会自动下载镜像，并完成环境的初始化

```
curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh | bash
```

执行以下命令启动雷池

```
docker compose up -d
```

经过以上步骤，你的雷池已经安装好了，下一步请参考 [登录雷池](/posts/guide_login)

## 离线安装 

如果你的服务器无法连接互联网环境，可以使用离线安装方式

> 这里忽略 Docker 安装的过程

首先，你要找一台能够访问互联网的服务器，执行以下命令拉取镜像，打包到 image.tar 文件中
```
docker pull chaitin/safeline-tengine:latest
docker pull chaitin/safeline-mgt-api:latest
docker pull chaitin/safeline-mario:latest
docker pull chaitin/safeline-detector:latest
docker pull postgres:15.2

docker save -o image.tar chaitin/safeline-tengine:latest chaitin/safeline-mgt-api:latest chaitin/safeline-mario:latest chaitin/safeline-detector:latest postgres:15.2
```

将 image.tar 文件传输到需要安装雷池的服务器上，执行以下命令加载镜像

```
docker load -i image.tar
```

执行以下命令创建并进入雷池安装目录

```
mkdir -p safeline        # 创建 safeline 目录
cd safeline              # 进入 safeline 目录
```

下载 [编排脚本](https://waf-ce.chaitin.cn/release/latest/compose.yaml) 并传输到 safeline 目录中


执行以下命令，生成雷池运行所需的相关环境变量

```
echo "SAFELINE_DIR=$(pwd)" > .env
echo "IMAGE_TAG=latest" >> .env
echo "MGT_PORT=9443" >> .env
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "SUBNET_PREFIX=169.254.0" >> .env
```

执行以下命令启动雷池

```
docker compose up -d
```

经过以上步骤，你的雷池已经安装好了，下一步请参考 [登录雷池](/posts/guide_login)

## 使用牧云助手安装

也可以使用 [牧云主机管理助手](https://collie.chaitin.cn/) 进行一键安装

![](/images/docs/guide_install/collie_apps.png)

参考视频教程 [用 “白嫖的云主机” 一键安装 “开源的Web防火墙”](https://www.bilibili.com/video/BV1sh4y1t7Pk/)

## 常见安装问题

请参考 [安装问题](/posts/faq_install)
