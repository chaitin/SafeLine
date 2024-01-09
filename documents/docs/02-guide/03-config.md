---
title: "配置站点"
---

# 配置站点

根据指导，完成站点配置

## 工作原理

雷池社区版主要以 **反向代理** 的方式工作，类似nginx。

**让网站流量先抵达雷池，经过雷池检测和过滤后，再转给原来的网站业务。**

## 配置界面

![config_site.gif](/images/gif/config_site.gif)

## 在单独的服务器部署雷池时配置（推荐）

### 开始配置

```shell
环境信息:
网站服务器:IPA，对外端口80，域名‘example.com’
雷池服务器:IPB

步骤:
1.将原网站流量指向雷池的IPB（必须）。例如修改域名解析服务，将域名解析到IPB
2.参考配置如下图
3.禁止网站服务器上，除雷池之外的访问。例如配置防火墙
```

![Alt text](/images/docs/guide_config/config_site2.png)

### 配置完成

如果浏览器访问`example.com:80`能获取到业务网站的响应，并且数据统计页的 “今日请求数” 增加，代表配置成功。


效果大致如下：

![Alt text](/images/docs/guide_config/deploy_on_separate_server.svg)

## 在网站服务器上部署雷池时配置

提示：不建议，因为这样单机的负载更高、设备宕机的概率更大。非纯净的环境还会提高升级失败的概率，故障排查更困难。

### 开始配置

```shell
环境信息:
网站服务器:IPA，对外端口80，域名‘example.com’

步骤:
1.需要原网站的监听修改为端口A，使80端口变成未使用状态，再进行配置
2.具体配置参考下图
```

![Alt text](/images/docs/guide_config/config_site1.png)

<!-- ### 参考视频

<video width="640" height="360" controls id="mp4" src="https://chaitin-marketing-public.oss-cn-beijing.aliyuncs.com/chaitin-website/safeline.mp4" type="video/mp4">

</video> -->

### 配置完成

如果浏览器访问`example.com:80`能获取到业务网站的响应，并且数据统计页的 “今日请求数” 增加，代表配置成功。

效果大致如图：

![Alt text](/images/docs/guide_config/deploy_on_web_server.svg)

## 和其他反代设备一起部署时配置

雷池作为反代设备，可以在任意位置接入主链路。

将接入位置的流量指向雷池，并在雷池的 “上游服务器” 处填写请求的下一跳服务器地址即可。

### 开始配置

```shell
环境信息:
网站服务器:IPA
雷池服务器:IPB
上游服务器:IPC，端口C
下游服务器:IPD，域名‘example.com’

步骤:
1.将下游nginx的流量指向雷池的IPC，访问端口指向80。
2.具体配置参考下图
```

![Alt text](/images/docs/guide_config/config_site3.png)

### 配置完成

如果浏览器访问`example.com:80`能获取到业务网站的响应，并且数据统计页的 “今日请求数” 增加，代表配置成功。

效果大致如图：

![Alt text](/images/docs/guide_config/deploy_with_other_server.png)

## 常见配置问题

请参考 [配置问题](/faq/config)
