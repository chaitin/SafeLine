---
title: "防护问题"
---

# 防护问题

记录常见的防护问题

## 攻击测试时未记录攻击，并且访问记录为 0

检查是否直接访问了源站，没有走雷池的配置站点访问

## 内网用户如何使用在线的威胁情报 IP？加白哪个域名？

威胁情报的云服务部署在百川云平台，域名是 **`https://challenge.rivers.chaitin.cn/`**

雷池部署在内网的师傅需要加白一下，就可以正常同步情报数据了。

## 如何记录所有访问雷池的请求 （如何开启访问日志）

默认情况下雷池是并不会保存请求记录的，如果需要保存请求记录，可以修改waf的安装目录下的**resources/nginx/nginx.conf**

![config_access_log.png](/images/docs/config_access_log.png)

如图所示，去掉文件第 98 行的注释，删除第 99 行的内容，保存后运行命令检查配置文件

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

配置生效后，访问日志将会保存至安装路径下的**logs/nginx**

**_注意：该操作会加快对硬盘的消耗，请定时清理访问日志_**

## 有信任的ip进入了社区黑名单怎么办 

社区黑名单的进入条件是需要多个雷池设备上报恶意行为，自动移出条件是连续一段时间没有被上报

由于网络环境的复杂性，存在ip被利用或者规则误报导致进入社区黑名单的情况

出于安全考虑，雷池社区不会主动为任何人移出社区黑名单

误报如何处理

1.排查信任ip的情况，确认是否真的安全

2.对信任的ip针对性配置加白规则

3.等待信任ip自动被移出社区黑名单

## 问题无法解决

1. 通过右上角搜索检索其他页面

2. 通过社群（官网首页加入微信讨论组）寻求帮助或者 Github issue 提交反馈，并附上排查的过程和截图
