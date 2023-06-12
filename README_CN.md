<p align="center">
  <img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_logo.png" width="120">
</p>
<h1 align="center">雷池 SafeLine 社区版</h1>
<h3 align="center">不让黑客越雷池半步</h3>
<br>
<p align="center">
  <img src="https://img.shields.io/badge/SafeLine-BEST_WAF-blue">
  <img src="https://img.shields.io/github/release/chaitin/safeline.svg?color=blue" />
  <img src="https://img.shields.io/github/release-date/chaitin/safeline.svg?color=blue&label=update" />
  <img src="https://img.shields.io/docker/v/chaitin/safeline-mgt-api?color=blue">
  <img src="https://img.shields.io/github/stars/chaitin/safeline?style=social">
</p>

<p align="center"> <a href="https://waf-ce.chaitin.cn/">官方网站</a> </p>
<p align="center"> 中文文档 | <a href="README.md">English</a> </p>

一款简单、好用的 WAF 工具。基于[长亭科技](https://www.chaitin.cn)王牌的 🤖️智能语义分析算法🤖️ 打造，专为社区设计。

## ✨ Demo

### 🔥🔥🔥 体验地址：https://demo.waf-ce.chaitin.cn:9443/

有一台运行在环境上 `http://127.0.0.1:8889` 的服务可以作为上游服务器测试使用。

![](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_detect_log.gif)

![](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

## 🚀 安装

### 1. 确保机器上正确安装 [Docker](https://docs.docker.com/engine/install/) 和 [Compose V2](https://docs.docker.com/compose/install/)
```shell
docker info # >= 20.10.6
docker compose version # >= 2.0.0
```

### 2. 部署安装

```shell
mkdir -p safeline && cd safeline
# 下载并执行 setup
curl -kfLsS https://waf-ce.chaitin.cn/release/latest/setup.sh | bash

# 运行
sudo docker compose up -d
```

#### 升级

##### 自动一键更新

**WARN: 雷池 SafeLine 服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。**

```shell
# 请到 compose.yaml 同级目录下执行下面脚本
curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash

# 如果有需要，可以手动删除旧版本镜像
docker rmi $(docker images | grep "safeline" | grep "none" | awk '{print $3}')
```
**有部分环境的默认 SafeLine 安装路径是在 `/data/safeline-ce`，安装之后可能会发现需要重新绑定 OTP、配置丢失等情况，可以修改 .env 的 `SAFELINE_DIR` 变量，指向 `/data/safeline-ce`**

##### 手动更新镜像

**适用于 docker hub 拉取镜像失败的场景，手动更新镜像，注意还是要执行 `upgrade.sh` 来处理 `.env` 的更新，否则有可能会因为缺少参数而启动失败。**

###### 1. 在一台能够从 docker hub 拉取镜像的机器上执行

```shell

# 拉取镜像
docker pull chaitin/safeline-tengine:latest
docker pull chaitin/safeline-mgt-api:latest
docker pull chaitin/safeline-mario:latest
docker pull chaitin/safeline-detector:latest
docker pull postgres:15.2

# 打包镜像
docker save -o image.tar chaitin/safeline-tengine:latest chaitin/safeline-mgt-api:latest chaitin/safeline-mario:latest chaitin/safeline-detector:latest postgres:15.2

# 传输到 SafeLine 要部署的目标服务器
# scp image.tar <target-server>:/root/
```
###### 2. 在目标服务器 load 镜像

```shell
docker load -i image.tar

curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash
```

## 🕹️ 快速使用

### 1. 登录

浏览器打开后台管理页面 `https://<waf-ip>:9443`。根据界面提示，使用 **支持 TOTP 的认证软件** 扫描二维码，然后输入动态口令登录：

![safeline_login.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_login.gif)

### 2. 添加站点

![safeline_website.gif](https://ctstack-oss.oss-cn-beijing.aliyuncs.com/veinmind/safeline-assets/safeline_website.gif)

<font color=grey>💡 TIPS: 添加后，执行 `curl -H "Host: <域名>" http://<WAF IP>:<端口>` 应能获取到业务网站的响应。</font>

### 3. 将网站流量切到雷池

- 若网站通过域名访问，则可将域名的 DNS 解析指向雷池所在设备
- 若网站前有 nginx 、负载均衡等代理设备，则可将雷池部署在代理设备和业务服务器之间，然后将代理设备的 upstream 指向雷池

### 4. 开始防护👌

试试这些攻击方式：

- 浏览器访问 `http://<IP或域名>:<端口>/webshell.php`
- 浏览器访问 `http://<IP或域名>:<端口>/?id=1%20AND%201=1`
- 浏览器访问 `http://<IP或域名>:<端口>/?a=<script>alert(1)</script>`

## 📖 FAQ

有任何问题请先查阅我们的 [FAQ 文档](FAQ.md)。

比如：
- [docker compose or docker-compose?](FAQ.md#docker-compose-还是-docker-compose)
- [站点如何配置](FAQ.md#站点配置问题)
- [配置完成之后，还是没有成功访问到上游服务器](FAQ.md#配置完成之后还是没有成功访问到上游服务器)

## 🏘️ 联系我们
1. 您可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议。
2. 扫描下方二维码可以加入雷池社区版用户讨论群进行详细讨论

<img src="https://waf-ce.chaitin.cn/images/wechat-light.png" width="30%" />

## ✨ CTStack
<img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/CT%20Stack-2.png" width="30%" />

雷池 SafeLine 现已加入 [CTStack](https://stack.chaitin.com/tool/detail?id=717) 社区

## Star History <a name="star-history"></a>

<a href="https://github.com/chaitin/safeline/stargazers">
        <img width="500" alt="Star History Chart" src="https://api.star-history.com/svg?repos=chaitin/safeline&type=Date">
      </a> 
