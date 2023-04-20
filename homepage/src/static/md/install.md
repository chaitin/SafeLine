<!-- ---
title: "快速部署"
category: "上手指南"
weight: 2
date: 2023-04-16T01:00:38+08:00
draft: true
--- -->

### 1. 确保机器上正确安装 [Docker](https://docs.docker.com/engine/install/) 和 [Compose V2](https://docs.docker.com/compose/install/)

```shell
docker info # >= 20.10.6
docker compose version # >= 2.0.0
```

### 2. 部署安装

```shell
mkdir -p safeline && cd safeline
# 下载并执行 setup
curl -fLsS https://waf-ce.chaitin.cn/release/latest/setup.sh | bash

# 运行
sudo docker compose up -d
```

## 🕹️ 快速使用

### 1. 登录

浏览器打开后台管理页面 `https://<waf-ip>:9443`。根据界面提示，使用 **支持 TOTP 的认证软件** 扫描二维码，然后输入动态口令登录：

![safeline_login.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_login.gif)

### 2. 添加站点

![safeline_website.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

`<font color=grey>`💡 TIPS: 添加后，执行 `curl -H "Host: <域名>" http://<WAF IP>:<端口>` 应能获取到业务网站的响应。`</font>`

### 3. 将网站流量切到雷池

- 若网站通过域名访问，则可将域名的 DNS 解析指向雷池所在设备
- 若网站前有 nginx 、负载均衡等代理设备，则可将雷池部署在代理设备和业务服务器之间，然后将代理设备的 upstream 指向雷池

### 4. 开始防护👌

试试这些攻击方式：

- 浏览器访问 `http://<IP或域名>:<端口>/webshell.php`
- 浏览器访问 `http://<IP或域名>:<端口>/?id=1%20AND%201=1`
- 浏览器访问 `http://<IP或域名>:<端口>/?a=<script>alert(1)</script>`
