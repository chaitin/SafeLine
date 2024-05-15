---
title: "Kong 集成雷池"
---

# Kong 集成雷池

[Kong](https://github.com/Kong/kong) 是一个云原生、快速、可扩展和分布式的微服务抽象层（也称为 API 网关或 API 中间件）。它通过插件提供了丰富的流量控制、安全、监控和运维功能。

# 使用方式

### 版本要求
* Kong >= 2.6.x
* Safeline >= 5.6.0

### 准备工作

参考 [APISIX 联动雷池](/docs/practice/apisix#准备工作) 的准备工作。

### 安装 Kong 插件

自定义插件可以通过 LuaRocks 安装。Lua 插件以 .rock 格式分发，这是一个自包含的包，可以从本地或远程服务器安装。

如果您使用了官方的 Kong Gateway 安装包，则 LuaRocks 实用程序应该已经安装在您的系统中。

1. 安装 safeline 插件
```shell
luarocks install kong-safeline
```
2. 启用 safeline 插件，在 kong.conf 配置文件中添加以下配置：
```shell
plugins = bundled,safeline              # Comma-separated list of plugins this node
                                        # should load. By default, only plugins
                                        # bundled in official distributions are
                                        # loaded via the `bundled` keyword.

```
3. 重启 Kong Gateway
```shell
kong restart
```

### 使用 Kong 插件
在某个 service 上启用 safeline 插件：
```shell
curl -X POST http://localhost:8001/services/{service}/plugins \
    --data "name=safeline" \
    --data "config.safeline_host=<your_detector_host>" \
    --data "config.safeline_port=<your_detector_port>"
```

### 测试防护效果
模拟简单的 SQL 注入攻击访问 kong ，如果返回 403 Forbidden，说明防护生效。
```shell
$ curl -X POST http://localhost:8000?1=1%20and%202=2

# you will receive a 403 Forbidden response
{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "8b41a021ea9541c89bb88f3773b4da24"}
```
打开雷池的控制台界面，可以看到雷池记录了完整的攻击信息。