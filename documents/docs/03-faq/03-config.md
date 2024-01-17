---
title: "配置问题"
---

# 配置问题

记录常见的配置问题

## 配置后网站无法访问，排查思路

如果按照指引配置了站点，但配置的网站无法访问

梳理问题可能存在的几个原因：

1. 配置站点错误，ip 错误、端口冲突等

2. 雷池端与配置的站点网络不通

3. 访问雷池配置的站点端口网络不通，对于雷池端已配置的端口没有被开放访问（防火墙、安全组等）

4. 同时存在其他错误的配置可能会导致新的配置一直不生效，检查有没有存在其他错误的配置

5. 雷池本身的状态不正常，使用 docker ps 检查容器状态 

## 排查步骤

1. 明确 “网站无法访问” 的具体表现：

   - 如果 `502 Bad Gateway tengine`：

      大概率是是雷池的上游服务器配置不正确，或者雷池无法访问到上游服务器，请继续按下面步骤排查。

      ![Alt text](/images/docs/guide_config/tengine_502.png)

   - 如果请求能够返回但是十分缓慢

     - 确认服务器负载是否正常，检查服务器的 CPU、内存、带宽使用情况

     - 在客户端执行命令，检查雷池服务器与上游服务器的网络：`curl -H "Host: <雷池 IP>" -vv -o /dev/null -s -w 'time_namelookup: %{time_namelookup}\ntime_connect: %{time_connect}\ntime_starttransfer: %{time_starttransfer}\ntime_total: %{time_total}\n' http://<上游服务器地址>`

       - 如果 time_namelookup 时间过大，请检查 dns server 配置
       - 如果 time_connect 时间过大，请检查雷池与上游服务器之间的网络状态
       - 如果 time_starttransfer 时间过大，请检查上游服务器状态，是否出现资源过载情况

   - 如果不是以上情况，继续下一步

2. 在客户端执行 `curl -v -H "Host: <域名或者IP>" http://<雷池 IP>:<雷池监听端口>` 。如能获取到业务网站的响应，如图，并且站点的 “今日访问量” +1，说明雷池配置正确，网络正常

   ![Alt text](/images/docs/guide_config/check_the_site1.png)

   如果浏览器无法访问，但这一步正常获取到响应，大概率是因为：

     - 网站域名还没有切到雷池，浏览器测试时访问的是 `http(s)://<雷池 IP>`，恰好业务服务上有 Host 验证，所以拒绝了该请求。这种情况需要修改本机 host，把域名解析到雷池 IP，再访问 `http(s)://<域名>`，才能准确测试
     - 网站业务做了其他一些特殊处理。例如访问后 301 跳转到了其他地址，需要具体排查网站业务的响应内容
     - 如果不能获取到响应，继续下一步

3. 在雷池设备上执行 `curl -v -H "Host: <域名或者IP>" http://<雷池 IP>:<雷池监听端口>`。如能获取到业务网站的响应，并且站点上 “今日访问量” +1，说明雷池配置正确

   - 如果步骤 2 失败而这里成功，说明客户端到雷池之间的网络存在问题。请排查网络，保证客户端可访问到雷池，检测防火墙、端口开放等
   - 如果不能获取到响应，继续下一步

4. 在雷池设备上执行 `curl -H "Host: <域名或者IP>" http://127.0.0.1:<雷池监听端口>`。如能获取到业务网站的响应，并且站点的 “今日访问量” +1，说明雷池配置正确

   - 如果步骤 3 失败而这里成功，且 `telnet <雷池 IP> <雷池监听端口>` 返回 `Unable to connect to remote host: Connection refused`，可能是被雷池设备上的防火墙拦截了。
   - 排查操作系统本身的防火墙，还有可能是云服务商的防火墙。请根据实际情况逐项排查，开放雷池监听端口的访问
   - 如果不能获取到响应，继续下一步

5. 在雷池设备上执行 `netstat -anp | grep <雷池监听端口>` 确认端口监听情况。正常情况下，应该有一个 nginx 进程监听在 `0.0.0.0:<雷池监听端口>`。

   - 没有的话请通过社群或者 Github issue 提交反馈，附上排查过程。有的话继续下一步

   ![Alt text](/images/docs/guide_config/check_listen_port.png)

6. 在雷池设备上 `curl -H "Host: <域名或者IP>" <上游服务器地址>`。如能获取到业务网站的响应，说明雷池设备和站点网络没有问题

   ![Alt text](/images/docs/guide_config/check_the_site2.png)

   - 如果步骤 4 失败而这里成功,可能是配置错误,查看配置站点教程确认配置是否正确，如无法解决，请通过社群或者 Github issue 提交反馈，附上排查过程

   - 如果这步失败，说明雷池和上游服务器之间的网络存在问题。请排查网络，确保雷池可以访问到上游服务器

## 配置完成后，测试时返回 400 Request Header Or Cookie Too Large

请麻烦检查是否形成了环路，即：雷池将请求转发给上游服务器后，上游服务器又将请求转发回雷池。

## 不同版本关闭防火墙的命令
   
Ubuntu 18.04 LTS 、 Ubuntu 20.04 LTS 、 Ubuntu 22.04 LTS

Debian 9 (Stretch)、Debian 10 (Buster)、Debian 11 (Bullseye)
```
关闭防火墙命令（UFW）：sudo ufw disable
注：Debian 默认可能不安装 UFW，依赖于 iptables。
```

CentOS 7、CentOS 8、RHEL 7、 RHEL 8、Fedora 32、 Fedora 33、Fedora 34
```
关闭防火墙命令（Firewalld）：sudo systemctl stop firewalld && sudo systemctl disable firewalld
```
openSUSE Leap 15.2、openSUSE Leap 15.3
```
关闭防火墙命令（通常是 SuSEfirewall2 或 firewalld）：
1.SuSEfirewall2, 使用 sudo SuSEfirewall2 stop
2.firewalld, 使用 sudo systemctl stop firewalld && sudo systemctl disable firewalld
```

## 如何对站点开启强制hppts访问、开启IPV6监听、使用HTTP/2

根据站点需求开启

开启路径：防护配置-通用配置-其他-站点通用配置 

 ![Alt text](/images/docs/guide_config/check_the_site3.png)


## 问题无法解决

1. 通过右上角搜索检索其他页面

2. 通过社群（官网首页加入微信讨论组）寻求帮助或者 Github issue 提交反馈，并附上排查的过程和截图
