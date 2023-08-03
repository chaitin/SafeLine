---
title: "测试防护效果"
group: "上手指南"
order: 5
---

# 测试防护效果

## 确认网站可以正常访问

根据雷池 WAF 配置的网站参数访问你的网站

打开浏览器访问 `http://<IP或域名>:<端口>/`

> 网站协议默认是 http，勾选 ssl 则为 https   
> 主机名可以是雷池的 IP，也可以是网站的域名（确保域名已经解析到雷池）  
> 端口是你在雷池页面中配置的网站端口  

若网站访问不正常，请参考 [网站无法访问](/docs/04-常见问题排查/03-faq_access.md)

## 尝试手动模拟攻击

打开浏览器，访问以下地址即可模拟出对应的攻击：

- 模拟 SQL 注入，请访问 `http://<IP或域名>:<端口>/?id=1%20AND%201=1`
- 模拟 XSS，请访问 `http://<IP或域名>:<端口>/?html=<script>alert(1)</script>`

通过浏览器，你将会看到雷池已经发现并阻断了攻击请求。

若请求没有被阻断，请参考 [防护不生效](/docs/04-常见问题排查/05-faq_other.md)

## 自动化测试防护效果

两条请求当然无法完整的测试雷池的防护效果，可以使用 blazehttp 自动化工具进行批量测试

#### 下载测试工具

- [Windows 版本](https://waf-ce.chaitin.cn/blazehttp/blazehttp_windows.exe)
- [Mac 版本(x64)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_mac_x64)
- [Mac 版本(M1)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_mac_m1)
- [Linux 版本(x64)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_linux_x64)
- [Linux 版本(ARM)](https://waf-ce.chaitin.cn/blazehttp/blazehttp_linux_arm64)
- [源码仓库](https://github.com/chaitin/blazehttp)

#### 准备测试样本

- [测试样本](https://waf-ce.chaitin.cn/blazehttp/testcases.zip)

下载请求样本后解压到 `testcases` 目录

#### 开始测试

1. 将测试工具 `blazehttp` 和测试样本 `testcases` 放在同一个目录下
2. 进入对应的目录
3. 使用以下请求开始测试

```
./blazehttp -t http://<IP或域名>:<端口> -g './testcases/**/*.http'
```

#### 测试效果展示

```
# 测试请求
./blazehttp -t http://192.168.0.1:8080 -g './testcases/**/*.http'

sending 100% |██████████████████████████████████████████| (18/18, 86 it/s)        
Total http file: 18, success: 18 failed: 0
Stat http response code

Status code: 403 hit: 16
Status code: 200 hit: 2

Stat http request tag

tag: sqli hit: 1
tag: black hit: 16
tag: file_include hit: 1
tag: file_upload hit: 1
tag: java_unserialize hit: 1
tag: php_unserialize hit: 1
tag: cmdi hit: 1
tag: ssrf hit: 1
tag: xslti hit: 1
tag: xss hit: 1
tag: xxe hit: 1
tag: asp_code hit: 1
tag: white hit: 2
tag: ognl hit: 1
tag: shellshock hit: 1
tag: ssti hit: 1
tag: directory_traversal hit: 1
tag: php_code hit: 1
```
