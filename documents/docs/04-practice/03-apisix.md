---
title: "APISIX 联动雷池"
---

# APISIX 联动雷池

Apache APISIX 是一个动态、实时、高性能的云原生 API 网关，提供了负载均衡、动态上游、灰度发布、服务熔断、身份认证、可观测性等丰富的流量管理功能。

雷池是由长亭科技开发的 WAF 系统，提供对 HTTP 请求的安全请求，提供完整的 API 管理和防护能力。

自 APISIX 3.5.0 之后的版本将内置长亭雷池 WAF 插件，在启用 chaitin-waf 插件后，流量将被转发给长亭 WAF 服务，用以检测和防止各种 Web 应用程序攻击，以保护应用程序和用户数据的安全。

## 开源仓库

apisix：https://github.com/apache/apisix

## 使用方式

### 版本要求
* APISIX >= 3.5.0
* Safeline >= 5.6.0

### 准备工作

社区版雷池的检测引擎默认以 unix socket 的方式提供服务，我们需要把他修改为 tcp 方式，供 t1k 插件调用。

进入雷池检测引擎的配置目录：

```
cd /data/safeline/resources/detector/
```

用文本编辑器打开目录里的 detector.yml 文件，我们需要将 bind 方式从 unix socket 改为 tcp，添加如下配置：

```
bind_addr: 0.0.0.0
listen_port: 8000
```

detector配置的属性值将覆盖容器内默认配置文件的同名属性值。这样我们就把雷池引擎的服务监听到了 8000 端口。

现在只需要把容器内的 8000 端口映射到宿主机即可，首先进入雷池的安装目录
```
cd /data/safeline/
```
然后用文本编辑器打开目录里的 compose.yaml 文件，为 detector 容器增加 ports 字段，暴露其 8000 端口，参考如下：

```
......

detect:
    ......
    ports:
    - 8000:8000

......
```

OK，改好了，在雷池安装目录下执行以下命令重启雷池即可生效。

```
docker compose down
docker compose up -d
```

### 修改雷池的默认端口

雷池和 apisix 默认都监听 9443 端口，如果在同一台机器上安装，需要修改雷池的默认端口。

在雷池的安装目录下，有一个名为 .env 的隐藏文件，其中的 MGT_PORT 字段，修改这里后使用上面的方法再重启雷池即可生效。

### 安装 APISIX

> 注意，要使用 APISIX 3.5.0 及以上版本

本文使用 apisix 的 docker 版本来做演示，克隆 apisix-docker 仓库，运行以下命令来安装：

```
git clone <https://github.com/apache/apisix-docker>
cd apisix-docker/compose
echo 'APISIX_DOCKER_TAG=3.5.0-debian' >> .env
docker compose -f docker-compose-release.yaml up -d
```

业务地址：http://127.0.0.1:9080/

管理地址：http://127.0.0.1:9180/

### 在 apisix 里绑定雷池

调用 apisix 的 api，设置雷池检测引擎的地址，供 apisix 调用，参考以下请求：

192.168.99.11 是我本地雷池的地址，替换为你的 IP 即可

```
curl <http://127.0.0.1:9180/apisix/admin/plugin_metadata/chaitin-waf> -H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' -X PUT -d '
{
  "nodes":[
     {
       "host": "192.168.99.11",
       "port": 8000
     }
  ]
}'
```

调用 apisix 的 api，设置一条路由，参考以下请求：

> 192.168.99.12:80 是上游服务器的地址，apisix 会将请求反向代理到这个地址。

```
curl <http://127.0.0.1:9180/apisix/admin/routes/1> -H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' -X PUT -d '
{
   "uri": "/*",
   "plugins": {
       "chaitin-waf": {}
    },
   "upstream": {
       "type": "roundrobin",
       "nodes": {
           "192.168.99.12:80": 1
       }
   }
}'
```

### 测试防护效果

经过以上步骤，雷池 + apisix 基本配置完成，可以试试效果了，请求 9080 端口，可以看到 apisix 成功代理了上游服务器的页面：

```sh
curl '<http://127.0.0.1:9080/>'
```

在请求中加入一个 a 参数，模拟 SQL 注入攻击：

```sh
curl '<http://127.0.0.1:9080/>' -d 'a=1 and 1=1'
```

返回了 HTTP 403 错误，从错误消息中可以看出，雷池成功抵御了此次攻击。

```json
{
  "code": 403,
  "success": false,
  "message": "blocked by Chaitin SafeLine Web Application Firewall",
  "event_id": "18e0f220f7a94127acb21ad3c1b4ac47"
}
```

打开雷池的控制台界面，可以看到雷池记录了完整的攻击信息

### 让 WAF 本身的站点拦截同时工作

由于更改了 detector 的监听方式，通过雷池控制台界面配置的站点将无法正常拦截恶意请求，需要在 nginx 配置中将检测引擎地址改为 http 方式。

首先拷贝一份 nginx 配置文件

```shell
cp /data/safeline/resources/nginx/safeline_unix.conf /data/safeline/resources/nginx/safeline_http.conf
```

编辑 nginx.conf 文件，将里面的 `include /etc/nginx/safeline_unix.conf;` 注释掉，添加下面的一行

```shell
# nginx.conf
include /etc/nginx/conf.d/*.conf;
include /etc/nginx/sites-enabled/generated;
# include /etc/nginx/safeline_unix.conf; 
include /etc/nginx/safeline_http.conf;
```

通过 docker inspect 获取 detector 的 ip

```shell
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' safeline-detector
```

编辑 safeline_http.conf 文件做出如下修改
```yaml
# safeline_http.conf
  upstream detector_server {
    keepalive   256;
    #server      unix:/resources/detector/snserver.sock;
    server      detector_ip:8000; # 这里填写上面获取到的 detector 的 ip
}
```

### 问题答疑

如果在使用过程中遇到问题，可以在加入 SDK 讨论群

![雷池SDK讨论群](/images/docs/sdk_chat.png)