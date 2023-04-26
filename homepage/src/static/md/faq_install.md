---
title: "无法安装"
category: "常见问题排查"
weight: 1
---

# 无法安装

## docker compose 还是 docker-compose？

`docker compose`（带空格）是 V2 版本，Go 写的。`docker-compose` 是 V1 版本，Python 写的，已经不维护了。

我们推荐使用 V2 版本的 `docker compose`，V1 可能会有兼容性等问题。

[docker/compose](https://github.com/docker/compose/) 中提到：

> For a smooth transition from legacy docker-compose 1.xx, please consider installing [compose-switch](https://github.com/docker/compose-switch) to translate `docker-compose ...` commands into Compose V2's `docker compose ....` . Also check V2's `--compatibility` flag.

其他参考：[https://stackoverflow.com/questions/66514436/difference-between-docker-compose-and-docker-compose](https://stackoverflow.com/a/66516826)




## ERROR: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

如描述，你需要启动 docker daemon 才能执行相关的命令。尝试 `systemctl start docker` 或者手动启动 `Docker Desktop` （MacOS 或者 Windows 用户）

As shown, you shall start docker first. Try `systemctl start docker` or manually start your docker desktop for MacOS/Windows users.

## docker not found, unable to deploy

如描述，你需要安装 `docker`。尝试 `curl -fLsS https://get.docker.com/ | sh` 或者 [Install Docker Engine](https://docs.docker.com/engine/install/)

## docker compose v2 not found, unable to deploy

如描述，你需要安装 `docker compose v2`。尝试 `[Install Docker Compose](https://docs.docker.com/compose/install/)`

### safeline-tengine 出现 Address already in use

`docker logs -f safeline-tengine` 容器日志中看到 `Address already in use` 信息。

端口冲突，根据报错信息中的端口号，排查是哪个服务占用了，手动处理冲突。

## safeline-postgres 出现 Operation not permitted

`docker logs -f safeline-postgres` 容器日志中看到 `Operation not permitted` 报错

可能是您的 docker 版本过低，升级 docker 到最新版本尝试一下。

## 如何自定义 SafeLine 安装路径？

基于最新的 `compose.yaml`，你可以手动修改 `.env` 文件的 `SAFELINE_DIR` 变量。

## 如何修改 SafeLine 后台管理的默认端口？本机 `:9443` 已经被别的服务占用了

基于最新的 `compose.yaml`，你可以手动添加 `MGT_PORT` 变量到 `.env` 文件。


