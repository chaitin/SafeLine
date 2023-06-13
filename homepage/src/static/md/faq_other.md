---
title: "其他问题"
group: "常见问题排查"
order: 10
---

# 其他问题
### 为什么我的检测日志中的攻击IP显示的是负载均衡的IP
因为雷池的源 IP 获取方式默认为**从 Socket 中获取**，如需更改请至 通用配置 -> 源 IP 获取方式 处更改
![get_source_ip.png](/images/docs/get_source_ip.png)

### 如何清理数据库中的统计信息和检测日志

***注意：该操作会清除所有日志信息，且不可恢复***

```shell
docker exec -it safeline-mgt-api cleanlogs
```


### 如何记录所有访问雷池的请求
默认情况下雷池是并不会保存请求记录的，如果需要保存请求记录，可以修改安装路径下的**resources/nginx/nginx.conf**
![config_access_log.png](/images/docs/config_access_log.png)
如图所示，去掉文件第99行的注释，删除第100行的内容，保存后运行命令检查配置文件
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

***注意：该操作会会加快对硬盘的消耗，请定时清理访问日志***

### 如何在业务服务器上获取用户真实地址
雷池默认将上一跳的IP加入至X-Forwarded-For中，如果雷池前，没有CDN或者nginx等代理服务，那么可以认为该IP为用户的真实IP

### 如果配置完成后，测试时返回400 Request Header Or Cookie Too Large
请麻烦检查是否形成了环路，即：雷池将请求转发给上游服务器后，上游服务器又将请求转发回雷池。

### 如何将雷池的日志导出到XXX
由于雷池的日志是存在Postgres 数据库中，用户可以通过logstash将数据库中的数据导出，并且利用大量的logstash output 插件导入至目标数据库中

参考资料: https://www.elastic.co/guide/en/logstash/current/plugins-inputs-jdbc.html

### 如何开启监听ipv6
雷池默认不开启ipv6, 如果需要开启ipv6，需手动修改安装路径下的**resources/nginx/sites-enabled/** 文件夹下对应域名的配置文件

如需同时监听ipv4与ipv6，则
```shell
server {
    listen [::]:80;
    server_name example.com;
}
```

如只需监听ipv6，则
```shell
server {
    listen [::]:80 ipv6only=on;
    server_name example.com;
}
```
***注意：页面上编辑当前站点会覆盖配置***
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