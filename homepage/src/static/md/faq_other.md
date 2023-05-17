---
title: "其他问题"
group: "常见问题排查"
order: 10
---

# 其他问题
### 为什么我的检测日志中的攻击IP显示的是负载均衡的IP
因为雷池的源 IP 获取方式默认为**从 Socket 中获取**，如需更改请至 通用配置 -> 源 IP 获取方式 处更改
![get_source_ip.png](/images/docs/get_source_ip.png)