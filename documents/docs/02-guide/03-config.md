---
title: "配置站点"
---

# 配置站点

根据指导，完成站点配置

## 工作原理

雷池社区版主要以 **反向代理** 的方式工作，类似于一台 nginx 服务。

**部署时，需要让网站流量先抵达雷池，经过雷池检测和过滤后，再转给原来的网站业务。**

建议优先熟悉反向代理概念再继续配置

## 配置界面

![config_site.gif](/images/gif/config_site.gif)

## 在单独的服务器部署雷池时配置（推荐）

### 开始配置

```shell
环境信息：
网站服务器:IPA，对外端口80，域名‘example.com’
部署雷池的服务器:IPB
目的：使用雷池的80端口接受请求进行防护

步骤：
1. 必须将网站流量指向雷池的IPB。例如修改域名解析服务的配置，将域名解析到雷池IPB
2. 具体配置参考下图
3. 禁止网站服务器（IPA）上，所有除了雷池之外的访问。例如配置防火墙
```

![Alt text](/images/docs/guide_config/config_site2.png)

### 配置完成

浏览器访问`example.com:80`,若能获取到业务网站的响应，并且站点上 “今日访问量” 增加，则代表配置成功。

效果大致如下：

![Alt text](/images/docs/guide_config/deploy_on_separate_server.svg)

## 在网站服务器上部署雷池时配置

提示：不建议，因为这样单机的负载更高、设备宕机的概率更大。非纯净的环境还会提高升级失败的概率，故障排查更困难。

### 开始配置

```shell
参考场景：
网站服务器:IPA，对外端口80，域名‘example.com’
服务器上部署雷池:雷池页面端口9443
目的：继续使用网站的80端口接受请求，进行防护
步骤：
1.需要原网站的监听修改为端口A，使80端口变成未使用状态，再进行配置
2.具体配置参考下图
```

![Alt text](/images/docs/guide_config/config_site1.png)

<!-- ### 参考视频

<video width="640" height="360" controls id="mp4" src="https://chaitin-marketing-public.oss-cn-beijing.aliyuncs.com/chaitin-website/safeline.mp4" type="video/mp4">

</video> -->

### 配置完成

浏览器访问`example.com:80`,若能获取到业务网站的响应，并且站点上 “今日访问量” 增加，则代表配置成功。

效果大致如图：

![Alt text](/images/docs/guide_config/deploy_on_web_server.svg)

## 和其他反代设备一起部署时配置

雷池作为反代设备，可以在任意位置接入主链路。

将接入位置的流量指向雷池，并在雷池的 “上游服务器” 处填写请求的下一跳服务器地址即可。

### 开始配置

```shell
环境信息：
网站服务器:IPA，对外端口80，域名‘example.com’
部署雷池的服务器:IPB
上游nginx：IPC，端口C
下游nginx：IPD

目的：使用雷池的80端口接受请求进行防护

步骤：
1. 将下游nginx的流量指向雷池的IPB，访问端口指向80。
2. 具体配置参考下图
```

![Alt text](/images/docs/guide_config/config_site3.png)

### 配置完成

浏览器访问`example.com:80`,若能获取到业务网站的响应，并且站点上 “今日访问量” 增加，则代表配置成功。

效果大致如图：

![Alt text](/images/docs/guide_config/deploy_with_other_server.svg)

## 常见配置问题

请参考 [配置问题](/faq/config)
