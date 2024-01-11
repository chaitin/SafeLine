---
title: "其他问题"
---

# 其他问题

记录不常见的其他问题

## 源 IP 显示不正确

雷池默认会通过 Socket 连接获取请求者的源 IP，如果请求在到达雷池之前，还经过了其他代理设备（如：反代、LB、CDN、AD 等），这种情况会影响雷池获取正确的源 IP 信息。

通常，代理设备都会将真实源 IP 通过 HTTP Header 的方式传递给下一跳设备。如下方的 HTTP 请求，在 `X-Forwarded-For` 和 `X-Real-IP` 两个 Header 中都包含了源 IP：

```
GET /path HTTP/1.1
Host: waf-ce.chaitin.cn
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
X-Forwarded-For: 110.123.66.233, 10.10.3.15
X-Real-IP: 110.123.66.233
```

> `X-Forwarded-For` 是链式结构，若请求经过了多级代理，这里将会按顺序记录每一跳的客户端 IP。

如果请求中没有包含存在源 IP 的相关 Header，可以通过修改前方代理设备的配置来解决。例如，Nginx 可以增加如下配置来传递 `X-Real-IP` 给后方设备：

```
location /xxx {
	proxy_pass http://xxx.xxx;
	...
	proxy_set_header X-Real-IP $remote_addr;
	...
}
```

遇到这种情况，打开雷池控制台的 “通用配置” 页面，将选项 “源 IP 获取方式” 的内容修改为 “从 HTTP Header 中获取”，并在对应的输入框中填入 `X-Real-IP` 即可。

![get_source_ip.png](/images/docs/get_source_ip.png)

## 清理数据库中的统计信息和检测日志

**_注意：该操作会清除所有日志信息，且不可恢复_**

```shell
docker exec safeline-mgt cleanlogs
```

## 将雷池的日志导出到 XXX

雷池社区版自发布以来经常有用户询问如何将拦截日志通过 syslog 转发至目标地址，接下来我们将尝试使用 `fluentd` 来实现这个需求。

首先，我们编写 `fluent.conf`，我们将读取 `mgt_detect_log_basic` 中的数据，并通过配置 syslog 转发出去。下面是 input 部分，match 部分可以参考参考文档中的 syslog 部分。

```
<source>
  @type sql

  host safeline-pg // 默认数据库地址，如果在 compose.yml 中改过，请使用改后值
  port 5432
  database safeline-ce // 数据库名
  adapter postgresql
  username safeline-ce // 默认用户名，如果在 compose.yml 中改过，请使用改后值
  password POSTGRES_PASSWORD // 数据库密码，见安装目录下 .env

  select_interval 60s  # optional
  select_limit 500     # optional

  state_file /var/run/fluentd/sql_state

  <table>
    table mgt_detect_log_basic
    update_column timestamp
    time_column timestamp  # optional
  </table>

  # detects all tables instead of <table> sections
  #all_tables
</source>
```

之后，来编写我们的 fluentd 的 Dockerfile

```
FROM fluent/fluentd:v1.16-1

# Use root account to use apk
USER root

# below RUN includes plugin as examples elasticsearch is not required
# you may customize including plugins as you wish
RUN apk add --no-cache --update --virtual .build-deps \
        sudo build-base ruby-dev \
 && apk add libpq-dev \
 && sudo gem install pg --no-document \
 && sudo gem install fluent-plugin-remote_syslog \
 && sudo gem sources --clear-all \
 && apk del .build-deps libpq-dev \
 && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

COPY fluent.conf /fluentd/etc/fluent.conf

USER fluent
```

最后，编译完成后，我们将容器跑起来，参考命令

```
echo "" > ./sql-state
docker run -d --restart=always --name safeline-fluentd --net safeline-ce -v ./sql-state:/var/run/fluentd/sql_state safeline-flunetd:latest
```

参考文档
[SQL input plugin for Fluentd event collector](https://github.com/fluent/fluent-plugin-sql)
[fluent-plugin-remote_syslog](https://github.com/fluent-plugins-nursery/fluent-plugin-remote_syslog)

## 有多个防护站点监听在同一个端口上，匹配顺序是怎么样的

如果域名处填写的分别为 IP 与域名，那么当使用进行 IP 请求时，则将会命中第一个配置的站点
![server_index02.png](/images/docs/server_index02.png)
以上图为例，如果用户使用 IP 访问，命中 example.com。

如果域名处填写的分别为域名与泛域名，除非准确命中域名，否则会命中泛域名，不论泛域名第几个配置。
![server_index01.png](/images/docs/server_index01.png)
以上图为例，如果用户使用 a.example.com 访问，命中 a.example.com。 如果用户使用 b.example.com，命中 \*.example.com。

## 自定义站点 nginx conf

雷池每次修改站点或者重启服务时，都会重新生成 **resources/nginx/sites-enabled/** 下的 nginx conf 文件。因为没法“智能”合并用户自定义的配置和自动生成的配置。但是也还是有方式能持久化地添加一些 nginx conf，不会被覆盖。

每个 `IF_backend_XXX` 的 location 中都有 `include proxy_params;` 这一行配置，且 `resources/nginx/proxy_params` 这个文件不会被修改站点、重启服务等动作覆盖。2.1.0 版本之后支持 `include custom_params/backend_XXX;` 可以自定义站点级的 nginx location 配置。

```shell
server {
    location ^~ / {
        proxy_pass http://backend_1;
        include proxy_params;
        include custom_params/IF_backend_1;
        # ...
    }
}
```

所以只需要根据需求修改对应的文件就可以了。比如在 `resources/nginx/proxy_params` 里面增加如下配置，即可支持 `X-Forwarded-Proto`：

```shell
proxy_set_header X-Forwarded-Proto $scheme;
```

修改完成后运行命令检查配置文件

```shell
docker exec safeline-tengine nginx -t
```

检查应显示

```shell
nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

最后应用配置文件

```shell
docker exec safeline-tengine nginx -s reload
```

## 攻击日志中的域名不是我的网站

攻击日志中显示的域名字段是取的 HTTP Header 中的 Host，如果这个字段不存在，则默认使用目的 IP 作为域名。如果客户端修改了 HTTP Header 的 Host，那么这里显示的就是修改之后的。

放一张截图更容易理解，注意下面「请求报文」中的 Host 字段：

![fake_host.jpg](/images/docs/fake_host.jpg)

## 上游服务器获取到的全都是雷池 WAF 的 IP，如何获取到真实 IP？

雷池默认透传了源 IP，放在 HTTP Header 中的 `X-Forwarded-For` 里面。

如果上游服务器是 NGINX，添加如下配置就可以。如果不是，需要自行配置解析 XFF

```
set_real_ip_from 0.0.0.0/0;
real_ip_header X-Forwarded-For;
```

## 是否支持 WebSocket ？

默认支持


## 问题无法解决

1. 通过右上角搜索检索其他页面

2. 通过社群（官网首页加入微信讨论组）寻求帮助或者 Github issue 提交反馈，并附上排查的过程和截图
