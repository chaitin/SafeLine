# FAQ

## docker compose or docker-compose?

tl;dr The `docker compose` (with a space), aka v2, is a newer project to migrate compose to Go. The original python project `docker-compose`, aka v1 has now been deprecated.

refers also: [https://stackoverflow.com/questions/66514436/difference-between-docker-compose-and-docker-compose](https://stackoverflow.com/a/66516826)

## Setup and Deploy

### What specs should I deploy SafeLine to protect my web services? / Minimum hardware specs requirements

1C1G is ok for running docker containers of SafeLine. But the specific hardware specs depends on the traffic characteristics your services, such as qps, network io and so on. For now, no detailed datasheet is provided for references.

### ERROR: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

As shown, you shall start docker first. Try `systemctl start docker`.

### docker not found, unable to deploy

As shown, you shall install `docker` first. Try `curl -fLsS https://get.docker.com/ | sh` or [Install Docker Engine](https://docs.docker.com/engine/install/)

### docker compose v2 not found, unable to deploy

As shown, you shall install `docker compose v2`. Try `[Install Docker Compose](https://docs.docker.com/compose/install/)`

### safeline-pg: Operation not permitted

`docker logs -f safeline-pg` with error `Operation not permitted`

Upgrade your docker engine and retry.

### safeline-tengine: Address already in use

`docker logs -f safeline-tengine` with error `Address already in use`

Port conflicts. Based on the port number in the error message, troubleshoot which service is occupying it and handle the conflict manually.

### How to customize safeline-ce installation path?

With the latest `compose.yaml`, you can manually modify the `SAFELINE_DIR` variable in the `.env` file after `setup.sh`.

### How to change the default port of SafeLine management, as for `:9443` is already in used by some other services?

With the latest `compose.yaml`, you can add `MGT_PORT` variable in `.env` file.

## Login

### TOTP login failed with correct code

TOTP is calculated and verified according to time. So check your server time.

# 常见问题

## docker compose 还是 docker-compose？

`docker compose`（带空格）是 V2 版本，Go 写的。`docker-compose` 是 V1 版本，Python 写的，已经不维护了。

我们推荐使用 V2 版本的 `docker compose`，V1 可能会有兼容性等问题。

[docker/compose](https://github.com/docker/compose/) 中提到：

> For a smooth transition from legacy docker-compose 1.xx, please consider installing [compose-switch](https://github.com/docker/compose-switch) to translate `docker-compose ...` commands into Compose V2's `docker compose ....` . Also check V2's `--compatibility` flag.

其他参考：[https://stackoverflow.com/questions/66514436/difference-between-docker-compose-and-docker-compose](https://stackoverflow.com/a/66516826)

## 安装部署

### 机器运行的最低配置

最低 1C1G 能运行，具体需要多少配置取决于你的业务流量特征，比如 QPS、网络吞吐等等，暂时没有详细的 datasheet 性能参考。

### 镜像下载缓慢甚至连接超时

这个是因为 docker hub 默认使用位于美西节点拉取镜像，可以自行配置国内镜像加速源

### ERROR: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

如描述，你需要启动 docker daemon 才能执行相关的命令。尝试 `systemctl start docker`

As shown, you shall start docker first. Try `systemctl start docker`.

### docker not found, unable to deploy

如描述，你需要安装 `docker`。尝试 `curl -fLsS https://get.docker.com/ | sh` 或者 [Install Docker Engine](https://docs.docker.com/engine/install/)

### docker compose v2 not found, unable to deploy

如描述，你需要安装 `docker compose v2`。尝试 `[Install Docker Compose](https://docs.docker.com/compose/install/)`

### safeline-tengine 出现 Address already in use

`docker logs -f safeline-tengine` 容器日志中看到 `Address already in use` 信息。

端口冲突，根据报错信息中的端口号，排查是哪个服务占用了，手动处理冲突。

### safeline-pg 出现 Operation not permitted

`docker logs -f safeline-pg` 容器日志中看到 `Operation not permitted` 报错

可能是您的 docker 版本过低，升级 docker 到最新版本尝试一下。

### 如何自定义 SafeLine 安装路径？

基于最新的 `compose.yaml`，你可以手动修改 `.env` 文件的 `SAFELINE_DIR` 变量。

### 如何修改 SafeLine 后台管理的默认端口？本机 `:9443` 已经被别的服务占用了

基于最新的 `compose.yaml`，你可以手动添加 `MGT_PORT` 变量到 `.env` 文件。

## 登录问题

### OTP 认证码登录失败

TOTP 是基于时间生成和校验的，请检查你的服务器时间是否同步。

## 站点配置问题

在没有 SafeLine 的时候，假设小明的域名 `xiaoming.com` 通过 DNS 解析到自己主机 `192.168.1.111`，上面在 `:8888` 端口监听了自己的服务（网站/博客/靶场）等等。

小明通过 `http://xiaoming.com:8888` 或者 `192.168.1.111:8888` 来访问自己的服务。

### 我该如何配置？/ 域名填什么？/ 端口怎么写？/ 上游服务器是什么？

目前社区版 SafeLine 支持的是反向代理的方式接入站点，也就是类似于一台 nginx 服务。这时候小明需要做的就是让流量先抵达 SafeLine，然后经过 SafeLine 检测之后，再转发给自己原先的业务。

小明只需要按照如下方式创建站点即可：
- `xiaoming.com` 填入页面的「域名」
- `:7777` 填入「端口」；或者别的任意非 `:8888`和 `:9443`（被 SafeLine 后台管理页面占用）端口
- `http://192.168.1.111:8888` 填入「上游服务器」

创建之后，就可以通过 `http://xiaoming.com:7777` 或者 `192.168.1.111:7777` 访问自己的服务了，这时候请求到 `http://xiaoming.com:7777` 的流量都会被 SafeLine 检测。经过 SafeLine 过滤后，安全的流量会被透传到原先的 `:8888` 业务服务器（即上游服务器）。

**注：直接访问 `http://xiaoming.com:8888` 的流量，仍然不会被 SafeLine 检测，因为流量并没有经过 SafeLine，而是绕过 SafeLine 直接打到了上游服务器上**

**如果按照如上配置，还是无法成功访问到我的上游服务器，接着往下看，尝试逐项进行问题排查。**

## 配置完成之后，还是没有成功访问到上游服务器

下面例子都还按照上面小明的环境情况介绍。

#### 1. netstat/ss/lsof 查看端口占用情况

先确认下 `0.0.0.0:7777` 端口是否有服务在监听。SafeLine 使用 Tengine 来作为代理服务，所以正常来说，应该有一个 nginx 进程监听在 `:7777` 端口。如果没有的话，可能是 SafeLine 的问题，请通过社群或者 Github issue 提交反馈。

如果有的话，继续往下排查。

#### 2. 是否是被非 SafeLine 的 nginx 监听

基于第一步，已经能确认 `:7777` 是被某个 nginx 进程监听了，但是并不能确认是被 SafeLine 自己的 nginx 监听。排查是否自己原先有 nginx conf 中配置了 server 监听 `:7777`。如果有的话，手动解决冲突。要么修改自己原先的 nginx conf，要么修改 SafeLine 的站点配置。

也可以直接通过 `docker logs -f safeline-tengine` 确认 SafeLine 是否有 nginx 报错说端口冲突。

*常见的情况就是自己原先有一个服务监听在 `:80`，SafeLine 上配置了站点也监听 `:80` 端口，就产生了冲突。*

如果没有的话，继续往下排查。

#### 3. 是否被防火墙拦截

有操作系统本身的防火墙，还有可能是云服务商的防火墙。根据实际情况逐项排查，配置开放端口的 TCP 访问。

出现如下情况，可能就是被中间某防火墙拦截了：
1. 在 `192.168.1.111` 上 curl -vv `127.0.0.1:7777` 能访问到业务，有 HTTP 返回码。
2. 在本机 curl -vv `192.168.1.111:7777` 不通，没有 HTTP 响应；`telnet 192.168.1.111 7777` 返回 `Unable to connect to remote host: Connection refused`

#### 4. SafeLine 是否能访问到上游服务器

小明的情况是 SafeLine 和业务在同一台机器，一般不会有不同机器之间的网络问题，但是也建议在 SafeLine 部署的机器上测试一下。如果是两台机器的情况下，需要考虑是否互相之间能正常通信。

直接 `curl -H "Host: <SafeLine-IP>" -vv http://xiaoming.com:8888` 测一下是否能访问到。如果不行，需要自行排查为什么 SafeLine 的机器没法访问到。

注：这里需要 -H 指定 Host `Host: <SafeLine-IP>` 进行连通性测试。收到比较多的反馈，在 WAF 上直接配置上游服务器为 HTTPS 的域名，比如 `https://xiaoming.com`。实际场景是希望先测试 WAF 能力正常后再把域名解析切到 WAF 进行上线。这种本地测试的场景，需要修改本机 host，把 `xiaoming.com` 解析到 `SafeLine-IP`，否则可能会无法成功代理。因为 SafeLine 向上游服务器转发时，代理请求中的 Host 使用的是原始 HTTP 请求中的 Host，此时需要自行判断上游业务服务器能够正确处理该代理请求（例如上游业务服务器在 Host 没有匹配自己的站点名称时，是否能够处理）

#### 5. 其他情况

如果执行了 1-4：
1. 确认有 nginx 进程监听了 SafeLine 机器的 `0.0.0.0:7777` 端口
2. 确认 SafeLine tengine 无端口冲突报错
3. 确认主机和云服务商的防火墙都没有限制 `:7777` 端口的 TCP 访问
4. 确认在 SafeLine 上能访问到「上游服务器」

问题还是没有解决，可能是 SafeLine 产品的问题，请通过社群或者 Github issue 提交反馈。
