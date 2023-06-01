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
默认情况下雷池是并不会保存请求记录的，如果需要保存请求记录，可以修改**resources/nginx/nginx.conf**
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
配置生效后，访问日志将会保存至**logs/nginx**

***注意：该操作会会加快对硬盘的消耗，请定时清理访问日志***