<p align="center">
  <img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_logo.png" width="120">
</p>
<h1 align="center">SafeLine Community Edition</h1>
<h3 align="center">Keep hackers at bay</h3>
<br>
<p align="center">
  <img src="https://img.shields.io/badge/SafeLine-BEST_WAF-blue">
  <img src="https://img.shields.io/github/release/chaitin/safeline.svg?color=blue" />
  <img src="https://img.shields.io/github/release-date/chaitin/safeline.svg?color=blue&label=update" />
  <img src="https://img.shields.io/docker/v/chaitin/safeline-mgt-api?color=blue">
  <img src="https://img.shields.io/github/license/chaitin/safeline?color=blue">
  <img src="https://img.shields.io/github/stars/chaitin/safeline?style=social">
</p>

<p align="center"> English | <a href="README_CN.md">中文文档</a> </p>

A simple and easy to use WAF tool. Built on [Chaitin Technology](https://www.chaitin.cn/en/)'s ace 🤖️Intelligent Semantic Analysis algorithm🤖️, designed for the community.

## ✨ Demo

### 🔥🔥🔥 Online Demo: https://47.104.188.110:9443/

There is a simple http server, listened on `http://127.0.0.1:8889`, can be used as for testing.

![](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_detect_log.gif)

![](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

## 🚀 Installation

### 1. Make sure [Docker](https://docs.docker.com/engine/install/) and [Compose V2](https://docs.docker.com/compose/install/) are installed correctly on the machine 
```
docker info
docker compose version
```

### 2. Download startup script

```shell
# download startup script
git clone git@github.com:chaitin/safeline.git
# or use `wget https://github.com/chaitin/safeline/releases/download/v1.0.3/safeline.zip -O safeline.zip && unzip safeline.zip`

# execute `./safeline-ce.sh` to generate the initial env params. WAF will be installed by default in the `/data/safeline-ce/` directory
cd safeline
./safeline-ce.sh

# launch
sudo docker compose up -d
```

## 🕹️ Quick Start

### 1. Login

Open admin page `https://<waf-ip>:9443` and scan qrcode with any authenticator Apps that support TOTP, enter the code to login.

![safeline_login.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_login.gif)

### 2. Create website

![safeline_website.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

<font color=grey>💡 TIPS: After creating website，execute `curl -H "Host: <Domain>" http://<WAF IP>:<Port>` to check if you can get correct response from web server.</font>

### 3. Deploy your website to SafeLine

- If your website is hosted by DNS, just modify your DNS record to WAF
- If your website is behind any reverse-proxy like nginx, you can modify your nginx conf and set upstream to WAF

### 4. Protected!👌

Try these:

- `http://<IP or Domain>:<Port>/webshell.php`
- `http://<IP or Domain>:<Port>/?id=1%20AND%201=1`
- `http://<IP or Domain>:<Port>/?a=<script>alert(1)</script>`

## 🏘️ Contact Us

1. You can make bug feedback and feature suggestions directly through GitHub Issues.
2. By scanning the QR code below (use wechat or qq), you can join the discussion group of SafeLine users for detailed discussions.

<img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_wx_light.jpg" width="30%" />
<img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_qq_light.jpg" width="30%" />

## ✨ CTStack
<img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/CT%20Stack-2.png" width="30%" />

SafeLine has already joined [CTStack](https://stack.chaitin.com/tool/detail?id=717) community.
