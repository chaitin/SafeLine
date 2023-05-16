---
title: "配置防护站点"
group: "上手指南"
order: 4
---

# 配置防护站点


![safeline_website.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

💡 TIPS: 添加后，执行 `curl -H "Host: <域名>" http://<WAF IP>:<端口>` 应能获取到业务网站的响应。

## 将网站流量切到雷池

- 若网站通过域名访问，则可将域名的 DNS 解析指向雷池所在设备
![DNS.png](https://https://waf-ce.chaitin.cn/images/docs/DNS.png)
- 若网站前有 nginx 、负载均衡等代理设备，则可将雷池部署在代理设备和业务服务器之间，然后将代理设备的 upstream 指向雷池 （）
![DNS.png](https://https://waf-ce.chaitin.cn/images/docs/LoadBlance.png)

## 测试防护效果

[测试防护效果](/posts/guide_test)
