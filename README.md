<p align="center">
  <img src="https://waf-ce.chaitin.cn/images/403.svg" width="120">
</p>
<h1 align="center">雷池 - 广受好评的社区 WAF</h1>
<br>
<p align="center">
  <img src="https://img.shields.io/badge/SafeLine-BEST_WAF-blue">
  <img src="https://img.shields.io/github/release/chaitin/safeline.svg?color=blue" />
  <img src="https://img.shields.io/github/release-date/chaitin/safeline.svg?color=blue&label=update" />
  <img src="https://img.shields.io/docker/v/chaitin/safeline-mgt-api?color=blue">
  <img src="https://img.shields.io/github/stars/chaitin/safeline?style=social">
</p>

<p align="center">
  <a href="https://waf-ce.chaitin.cn/">官方网站</a> | 
  <a href="https://demo.waf-ce.chaitin.cn:9443/dashboard">在线 Demo</a> | 
  <a href="https://waf-ce.chaitin.cn/posts/guide_introduction">技术文档</a> | 
  <a href="README_EN.md">For English</a>
</p>

一款足够简单、足够好用、足够强的免费 WAF。基于业界领先的语义引擎检测技术，作为反向代理接入，保护你的网站不受黑客攻击。

核心检测能力由智能语义分析算法驱动，专为社区而生，不让黑客越雷池半步。

<img src="https://waf-ce.chaitin.cn/images/album/0.png" />

<h4 align="center">相关源码仓库</h4>
<p align="center">
  <a href="https://github.com/chaitin/yanshi">语义分析自动机引擎</a> | 
  <a href="https://github.com/chaitin/safeline-open-platform">流量分析插件</a> | 
  <a href="https://github.com/chaitin/lua-resty-t1k">T1K 协议</a> |
  <a href="https://github.com/chaitin/blazehttp">测试工具</a>
</p>

## 相关特性

#### 便捷性

采用容器化部署，一条命令即可完成安装，0 成本上手。安全配置开箱即用，无需人工维护，可实现安全躺平式管理。

#### 安全性

首创业内领先的智能语义分析算法，精准检测、低误报、难绕过。语义分析算法无规则，面对未知特征的 0day 攻击不再手足无措。

#### 高性能

无规则引擎，线性安全检测算法，平均请求检测延迟在 1 毫秒级别。并发能力强，单核轻松检测 2000+ TPS，只要硬件足够强，可支撑的流量规模无上限。

#### 高可用

流量处理引擎基于 Nginx 开发，性能与稳定性均可得到保障。内置完善的健康检查机制，服务可用性高达 99.99%。

## 🚀 安装

### 配置需求

- 操作系统：Linux
- 指令架构：x86_64
- 软件依赖：Docker 20.10.6 版本以上
- 软件依赖：Docker Compose 2.0.0 版本以上
- 最小化环境：1 核 CPU / 1 GB 内存 / 10 GB 磁盘


### 一键安装

```
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh)"
```

> 更多安装方式请参考 <a href="https://waf-ce.chaitin.cn/posts/guide_install">安装雷池</a>

## 🕹️ 快速使用

### 登录

浏览器打开后台管理页面 `https://<waf-ip>:9443`。根据界面提示，使用 **支持 TOTP 的认证软件** 扫描二维码，然后输入动态口令登录：

![login.gif](https://waf-ce.chaitin.cn/images/gif/login.gif)

### 配置防护站点

雷池以反向代理方式接入，优先于网站服务器接收流量，对流量中的攻击行为进行检测和清洗，将清洗过后的流量转发给网站服务器。

![config.gif](https://waf-ce.chaitin.cn/images/gif/config_site.gif)

<font color=grey>💡 TIPS: 添加后，执行 `curl -H "Host: <域名>" http://<WAF IP>:<端口>` 应能获取到业务网站的响应。</font>

### 测试效果

使用以下方式尝试模拟黑客攻击，看看雷池的防护效果如何

- 浏览器访问 `http://<IP或域名>:<端口>/?id=1%20AND%201=1`
- 浏览器访问 `http://<IP或域名>:<端口>/?a=<script>alert(1)</script>`

![log.gif](https://waf-ce.chaitin.cn/images/gif/detect_log.gif)

> 如果你需要进行深度测试，请参考 <a href="https://waf-ce.chaitin.cn/posts/guide_test">测试防护效果</a>

### FAQ

- [安装问题](https://waf-ce.chaitin.cn/posts/faq_install)
- [登录问题](https://waf-ce.chaitin.cn/posts/faq_login)
- [网站无法访问](https://waf-ce.chaitin.cn/posts/faq_access)
- [配置问题](https://waf-ce.chaitin.cn/posts/faq_config)
- [其他问题](https://waf-ce.chaitin.cn/posts/faq_other)

## 🏘️ 联系我们

1. 可以通过 GitHub Issue 直接进行 Bug 反馈和功能建议
2. 可以扫描下方二维码加入雷池社区版用户讨论群

<img src="https://waf-ce.chaitin.cn/images/wechat-230717.png" width="30%" />

## Star History <a name="star-history"></a>

<a href="https://github.com/chaitin/safeline/stargazers">
    <img width="500" alt="Star History Chart" src="https://api.star-history.com/svg?repos=chaitin/safeline&type=Date">
</a> 
