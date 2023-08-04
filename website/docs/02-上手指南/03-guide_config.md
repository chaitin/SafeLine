---
title: "配置防护站点"
---

# 配置防护站点


![config_site.gif](https://waf-ce.chaitin.cn/images/gif/config_site.gif)

💡 TIPS: 添加后，执行 `curl -H "Host: <域名>" http://<WAF IP>:<端口>` 应能获取到业务网站的响应。

## 将网站流量切到雷池

- 若网站通过域名访问，则可将域名的 DNS 解析指向雷池所在设备
![DNS.png](/images/docs/DNS.png)
- 若网站前有 nginx 、负载均衡等代理设备，则可将雷池部署在代理设备和业务服务器之间，然后将代理设备的 upstream 指向雷池
![DNS.png](/images/docs/LoadBlance.png)

## 如何配置https
![safeline_https_website.gif](/images/docs/safeline_https_website.gif)

## 测试防护效果

[测试防护效果](/docs/上手指南/guide_test)
