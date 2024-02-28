---
title: "测试防护"
---

# 测试防护

使用手工或者自动的方式测试防护效果

## 确认网站可以正常访问

根据雷池 WAF 配置的网站参数访问你的网站。

打开浏览器访问 `http://<IP或域名>:<端口>/`。

> 网站协议默认是 http，勾选 ssl 则为 https  
> 主机名可以是雷池的 IP，也可以是网站的域名（确保域名已经解析到雷池）  
> 端口是你在雷池页面中配置的网站端口

若网站访问不正常，请参考 [配置问题](/03-faq/03-config.md)。

整体监测流程参考：

![flow.png](/images/docs/flow.svg)


## 尝试手动模拟攻击

访问以下地址模拟出对应的攻击：

- 模拟 SQL 注入，请访问 `http://<IP或域名>:<端口>/?id=1%20AND%201=1`
- 模拟 XSS，请访问 `http://<IP或域名>:<端口>/?html=<script>alert(1)</script>`

通过浏览器，你将会看到雷池已经发现并阻断了攻击请求。
![Alt text](/images/docs/guide_config/protection_page.png)

若请求没有被阻断，请参考 [防护问题](/03-faq/04-test.md)

## 自动化测试防护效果

两条请求当然无法完整的测试雷池的防护效果，可以使用 blazehttp 自动化工具进行批量测试

### 下载测试工具

- [Windows 版本](https://waf-ce.chaitin.cn/blazehttp/blazehttp_windows.exe)
- [Mac 版本(x64)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_mac_x64)
- [Mac 版本(M1)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_mac_m1)
- [Linux 版本(x64)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_linux_x64)
- [Linux 版本(ARM)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_linux_arm64)
- [源码仓库](https://github.com/chaitin/blazehttp)

### 准备测试样本

- [测试样本](https://waf-ce.chaitin.cn/blazehttp/testcases.zip)

下载请求样本后解压到 `testcases` 目录

### 开始测试

1. 将测试工具 `blazehttp` 和测试样本 `testcases` 放在同一个目录下
2. 进入对应的目录
3. 使用以下请求开始测试

```
./blazehttp -t http://<IP或域名>:<端口>
```

### 测试效果展示

```sh
# 测试请求
.//blazehttp -t http://127.0.0.1:8008
sending 100% |█████████████████████████████████████████████████████████| (33669/33669, 940 it/s) [35s:0s]
总样本数量: 33669    成功: 33669    错误: 0
检出率: 71.65% (恶意样本总数: 575 , 正确拦截: 412 , 漏报放行: 163)
误报率: 0.07% (正常样本总数: 33094 , 正确放行: 33071 , 误报拦截: 23)
准确率: 99.45% (正确拦截 + 正确放行）/样本总数
平均耗时: 1.00毫秒
```

## 常见防护问题

请参考 [防护问题](/faq/test)
