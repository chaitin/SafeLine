---
title: "网站无法访问"
group: "常见问题排查"
order: 3
---

# 网站无法访问

为了方便讨论问题，我们假设: </br>

在没有 SafeLine 的时候，假设小明的域名 `xiaoming.com` 通过 DNS 解析到自己主机 `192.168.1.111`，上面在 `:8888` 端口监听了自己的服务（网站/博客/靶场）等等。

小明通过 `http://xiaoming.com:8888` 或者 `192.168.1.111:8888` 来访问自己的服务。

## 如果返回502
![DNS.png](/images/docs/tengine_502.png)
如果出现类似相应，请检查上游服务器地址是否配置正确或者雷池是否能够访问能访问到上游服务器

## 请求返回缓慢
1. 首先检查服务器负载情况
2. 检查雷池服务器与上游服务器的网络状况，检查命令`curl -H "Host: <SafeLine-IP>" -vv -o /dev/null -s -w 'time_namelookup: %{time_namelookup}\ntime_connect: %{time_connect}\ntime_starttransfer: %{time_starttransfer}\ntime_total: %{time_total}\n' http://xiaoming.com:8888` </br>
如果 `time_namelookup` 时间过大，请检查dns server配置 </br>
如果 `time_connect` 时间过大，请检查与上游服务器之间的网络状态 </br>
如果 `time_starttransfer` 时间过大，请检查上游服务器状态，是否出现资源过载情况 </br>