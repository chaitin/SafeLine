---
title: "雷池技术架构"
---

# 雷池技术架构

查看雷池的服务架构图。最上面虚线框住的是数据流，也就是访问业务服务器的流量数据的流动情况。中间框起来的部分是雷池的各个服务。

![framework](/images/docs/framework.png)

各个容器和服务说明：

| 名称              | 定义         | 详情                                                    |
| ----------------- | ------------ | ------------------------------------------------------- |
| safeline-mgt  | 管理容器     | 接收管理后台行为，向其他服务或容器推送消息              |
| safeline-detector | 检测容器     | 执行检测的容器，从 Tengine 进入的流量会转发到该节点检测 |
| safeline-mario    | 日志容器     | 记录与统计恶意行为的节点                                |
| safeline-tengine  | 网关         | 转发网关，有简单的过滤功能                              |
| safeline-pg | 关系型数据库 | 存储攻击日志、保护站点、黑白名单配置的数据库            |

对于后台管理人员，可以直接通信的节点为管理服务 `safeline-mgt`，该节点负责：

- 向 Tengine 网关推送自定义配置并利用 NGINX 命令进行 reload 热更新
- 自定义检测规则（黑白名单等）并向检测引擎 `safeline-detector` 推送
- 直接读取 `postgres` 数据库，向后台管理人员返回日志、统计、当前配置等

## 各个配置文件说明

### .env 文件

用于设置 `compose.yaml` 要引用的环境变量

```bash
echo "SAFELINE_DIR=$(pwd)" >> .env  # 设置当前路径为雷池社区版的根路径
echo "IMAGE_TAG=latest" >> .env  # 设置镜像的 tag
echo "MGT_PORT=9443" >> .env  # 管理容器服务使用的端口
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env  # /dev/urandom是一个很长的随机数文本，tr -dc 命令用于删除非字母、非数字的字符，用于生成随机的 postgres 密码
echo "SUBNET_PREFIX=172.22.222" >> .env  # 定义 docker 虚拟网卡的子网前缀
```

### compose.yml 文件

用于启动多个容器

```yml
# 基于3.10 版本进行说明，最新配置文件可能存在部分变动
networks:
  safeline-ce:
    name: safeline-ce # 定义该子网名称
    driver: bridge # 该子网为网桥模式
    ipam:
      driver: default
      config:
        - gateway: ${SUBNET_PREFIX:?SUBNET_PREFIX required}.1 # 定义网关为 SUBNET_PREFIX.1，若按上文设置，此处为 172.22.222.1
          subnet: ${SUBNET_PREFIX}.0/24
    driver_opts:
      com.docker.network.bridge.name: safeline-ce
services:
  postgres:
    container_name: safeline-postgres
    restart: always # 容器启动失败或崩溃时自动重启
    image: postgres:15.2
    volumes: # 开启的映射文件夹
      - ${SAFELINE_DIR}/resources/postgres/data:/var/lib/postgresql/data
      - /etc/localtime:/etc/localtime:ro
    environment:
      - POSTGRES_USER=safeline-ce
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:?postgres password required}
    networks: # 使用上文的 safeline-ce 网络，IP 为 172.22.222.2
      safeline-ce:
        ipv4_address: ${SUBNET_PREFIX}.2
    cap_drop:
      - net_raw
    command: [postgres, -c, max_connections=200] # 设置 postgres 的最大连接数
  management:
    container_name: safeline-mgt-api
    restart: always
    image: chaitin/safeline-mgt-api:${IMAGE_TAG:?image tag required}
    volumes:
      - ${SAFELINE_DIR?safeline dir required}/resources/management:/resources/management
      - ${SAFELINE_DIR}/resources/nginx:/resources/nginx
      - ${SAFELINE_DIR}/logs:/logs
      - /etc/localtime:/etc/localtime:ro
    ports:
      - ${MGT_PORT:-9443}:1443
    environment:
      - MANAGEMENT_RESOURCES_DIR=/resources/management
      - NGINX_RESOURCES_DIR=/resources/nginx
      - DATABASE_URL=postgres://safeline-ce:${POSTGRES_PASSWORD}@127.0.0.1/safeline-ce
      - MANAGEMENT_LOGS_DIR=/logs/management
    networks:
      safeline-ce: # 使用上文的 safeline-ce 网络，IP 为 172.22.222.4
        ipv4_address: ${SUBNET_PREFIX}.4
    cap_drop:
      - net_raw
  detector:
    container_name: safeline-detector
    restart: always
    image: chaitin/safeline-detector:${IMAGE_TAG}
    volumes:
      - ${SAFELINE_DIR}/resources/detector:/resources/detector
      - ${SAFELINE_DIR}/logs/detector:/logs/detector
      - /etc/localtime:/etc/localtime:ro
    environment:
      - LOG_DIR=/logs/detector
    networks:
      safeline-ce: # 使用上文的 safeline-ce 网络，IP 为 172.22.222.5
        ipv4_address: ${SUBNET_PREFIX}.5
    cap_drop:
      - net_raw
  mario:
    container_name: safeline-mario
    restart: always
    image: chaitin/safeline-mario:${IMAGE_TAG}
    volumes:
      - ${SAFELINE_DIR}/resources/mario:/resources/mario
      - ${SAFELINE_DIR}/logs/mario:/logs/mario
      - /etc/localtime:/etc/localtime:ro
    environment:
      - LOG_DIR=/logs/mario
      - GOGC=100
      - DATABASE_URL=postgres://safeline-ce:${POSTGRES_PASSWORD}@safeline-postgres/safeline-ce
    networks:
      safeline-ce: # 使用上文的 safeline-ce 网络，IP 为172.22.222.6
        ipv4_address: ${SUBNET_PREFIX}.6
    cap_drop:
      - net_raw
  tengine:
    container_name: safeline-tengine
    restart: always
    image: chaitin/safeline-tengine:${IMAGE_TAG}
    volumes:
      - ${SAFELINE_DIR}/resources/nginx:/etc/nginx
      - ${SAFELINE_DIR}/resources/management:/resources/management
      - ${SAFELINE_DIR}/resources/detector:/resources/detector
      - ${SAFELINE_DIR}/logs/nginx:/var/log/nginx
      - /etc/localtime:/etc/localtime:ro
      - ${SAFELINE_DIR}/resources/cache:/usr/local/nginx/cache
      - /etc/resolv.conf:/etc/resolv.conf
    environment:
      - MGT_ADDR=${SUBNET_PREFIX}.4:9002 # 配置 mgt-api 的 grpc 服务器地址，用于与 mgt-api 容器通信
    ulimits:
      nofile: 131072
    network_mode: host # Tengine 直接使用宿主机网络
```

### 各个服务的运行日志

```bash

docker ps  # 查看各容器状态
docker logs -f <comtainer_name>  # 输出容器 std 日志

# 一些服务运行中的持久化也会存储在磁盘上，目录结构如下

root@user:/path/to/safeline-ce/logs# tree
.
├── detector
│   └── snserver.log    #检测容器的输出日志
├── management
│   ├── nginx.log     # Tengine 容器中 NGINX 日志输出
│   └── webserver.log    # safeline-mgt-api 容器的日志输出
├── mario
│   └── mario.log    # 流量日志输出
└── nginx
    ├── error.log  # NGINX 的错误日志
    └── tcd.log  # tcd 是 Tengine 用于与 safeline-mgt-api 通信的网络代理进程，该文件存储了两者的通信日志
```
