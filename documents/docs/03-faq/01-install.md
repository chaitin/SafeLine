---
title: "安装问题"
---

# 安装问题

记录常见的安装问题

## 在线安装失败

1. 检查是否手动关闭防火墙

2. 检测配置是否符合最低的配置要求

> 参考 [环境检测](/guide/install#环境检测) 方式

3. 如果连接 Docker Hub 网络不稳，导致镜像下载失败（超时）:

> docker hub 默认使用国外节点拉取镜像，可以自行搜索配置国内镜像加速源

> 采用 [离线安装](/guide/install#离线安装) 方式

## 安装时遇到报错处理方法

#### 报错：ERROR: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

需要安装 docker。尝试 `curl -fLsS https://get.docker.com/ | sh` 或者 [Install Docker Engine](https://docs.docker.com/engine/install/)。

#### 报错：docker not found, unable to deploy

failed to create network safeline-ce
safeline-ce 是雷池部署时候创建的 network，出现类似报错，先重启下 dockerd 之后重试

需要启动 docker daemon 才能执行相关的命令。尝试 `systemctl start docker`。

#### 报错：docker compose v2 not found, unable to deploy

需要安装 `docker compose v2`。尝试 `[Install Docker Compose](https://docs.docker.com/compose/install/)`。

#### 报错： `failed to create network safeline-ce`

safeline-ce 是雷池部署时候创建的 network，出现类似报错，先重启下 dockerd 之后重试。

#### 报错： safeline-tengine 出现 Address already in use

`docker logs -f safeline-tengine` 容器日志中看到 `Address already in use` 信息。

端口冲突，根据报错信息中的端口号，排查是哪个服务占用了，手动处理冲突。

#### 报错：safelint-mgt-api 出现 Operation not permitted

`docker logs -f safelint-mgt-api` 容器日志中看到 `runtime/cgo: pthread_create failed: Operation not permitted` 报错，这个错误一般会在 docker 20.10.9 及以下发生。

- 最推荐的方式是升级 docker 到最新版本尝试解决这个问题。
- 或您的系统支持配置 seccomp （执行 `grep CONFIG_SECCOMP= /boot/config-$(uname -r)` 输出 `CONFIG_SECCOMP=y` 则为支持）,
  则可以在雷池工作目录下载 [seccomp](https://waf-ce.chaitin.cn/release/latest/seccomp.json) 并且编辑 compose.yaml 文件，
  在 management 下加入如下配置项，然后执行 `docker compose down && docker compose up -d` 来尝试解决这个问题:

```yaml
security_opt:
  - seccomp=./seccomp.json
```

#### 报错：safeline-pg 出现 Operation not permitted

`docker logs -f safeline-pg` 容器日志中看到 `Operation not permitted` 报错。

可能是您的 docker 版本过低，升级 docker 到最新版本尝试一下。

#### 其他奇怪的报错比如：It does not belong to any of this network's subnets...等等

查看[如何卸载](#如何卸载) ，卸载以后重新安装一次

## 如何自定义 SafeLine 安装路径？

基于最新的 `compose.yaml`，你可以手动修改 `.env` 文件的 `SAFELINE_DIR` 变量。

## 雷池和业务服务可以部署到同一台机器中吗？

可以，但是不建议，机器负载将高于分开部署。

## MacOS/Windows 是否支持安装雷池

社区版暂不支持，如有需求咨询企业版。

## docker compose 还是 docker-compose？

属于两个版本，推荐使用 docker compose

参考资料：https://stackoverflow.com/questions/66514436/difference-between-docker-compose-and-docker-compose

## 如何修改 SafeLine 后台管理的默认端口？比如：本机 `:9443` 已经被别的服务占用

使用 `ss -antp|grep LISTEN` 确认端口使用情况，找到未被占用端口

修改在安装目录(默认 safeline)下的隐藏文件`.env` 文件，你可以手动添加 `MGT_PORT` 变量到 `.env` 文件。

文件修改后，需要等重启才会生效。

在安装目录(默认 safeline)下执行 `docker compose down && docker compose up -d`

## 如何卸载

在安装目录(默认 safeline)下

根据本地的compose版本，执行 `docker compose down` 或者 `docker-compose down`



## 问题无法解决

1. 通过右上角搜索检索其他页面

2. 通过社群（官网首页加入微信讨论组）寻求帮助或者 Github issue 提交反馈，并附上排查的过程和截图
