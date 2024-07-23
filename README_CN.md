<p align="center">
  <img src="/images/banner.png" width="400" />
</p>

<h4 align="center">
  SafeLine - 雷池 - 不让黑客越过半步
</h4>

<p align="center">
  <a target="_blank" href="https://waf-ce.chaitin.cn/">🏠 官网</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://docs.waf-ce.chaitin.cn/">📖 文档</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://demo.waf-ce.chaitin.cn:9443/">🔍 演示环境</a> &nbsp; | &nbsp;
  <a target="_blank" href="/images/wechat.png">🙋‍♂️ 社区微信群</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://github.com/chaitin/SafeLine">国际版</a>
</p>

## 👋 项目介绍

SafeLine，中文名 "雷池"，是一款简单好用, 效果突出的 **`Web 应用防火墙(WAF)`**，可以保护 Web 服务不受黑客攻击。

雷池通过过滤和监控 Web 应用与互联网之间的 HTTP 流量来保护 Web 服务。可以保护 Web 服务免受 `SQL 注入`、`XSS`、 `代码注入`、`命令注入`、`CRLF 注入`、`ldap 注入`、`xpath 注入`、`RCE`、`XXE`、`SSRF`、`路径遍历`、`后门`、`暴力破解`、`CC`、`爬虫` 等攻击。

#### 💡 工作原理

<img src="/images/how-it-works.png" width="800" />

雷池通过阻断流向 Web 服务的恶意 HTTP 流量来保护 Web 服务。雷池作为反向代理接入网络，通过在 Web 服务前部署雷池，可在 Web 服务和互联网之间设置一道屏障。

雷池的核心功能如下:

- 防护 Web 攻击
- 防爬虫, 防扫描
- 前端代码动态加密
- 基于源 IP 的访问速率限制
- HTTP 访问控制

#### ⚡️ 项目截图

| <img src="./images/screenshot-1.png" width=370 /> | <img src="./images/screenshot-2.png" width=370 /> |
| ------------------------------------------------- | ------------------------------------------------- | 
| <img src="./images/screenshot-3.png" width=370 /> | <img src="./images/screenshot-4.png" width=370 /> | 

查看 [演示环境](https://demo.waf-ce.chaitin.cn:9443/)

## 🔥 核心能力

对于你的网站而言, 雷池可以实现如下效果:

- **`阻断 Web 攻击`**
  - 可以防御所有的 Web 攻击，例如 `SQL 注入`、`XSS`、`代码注入`、`操作系统命令注入`、`CRLF 注入`、`XXE`、`SSRF`、`路径遍历` 等等。
- **`限制访问频率`**
  - 限制用户的访问速率，让 Web 服务免遭 `CC 攻击`、`暴力破解`、`流量激增` 和其他类型的滥用。
- **`人机验证`**
  - 互联网上有来自真人用户的流量，但更多的是由爬虫, 漏洞扫描器, 蠕虫病毒, 漏洞利用程序等自动化程序发起的流量，开启雷池的人机验证功能后真人用户会被放行，恶意爬虫将会被阻断。
- **`身份认证`**
  - 雷池的 "身份认证" 功能可以很好的解决 "未授权访问" 漏洞，当用户访问您的网站时，需要输入您配置的用户名和密码信息，不持有认证信息的用户将被拒之门外。
- **`动态防护`**
  - 在用户浏览到的网页内容不变的情况下，将网页赋予动态特性，对 HTML 和 JavaScript 代码进行动态加密，确保每次访问时这些代码都以随机且独特的形态呈现。

#### 🧩 核心能力展示

|                               | Legitimate User                                     | Malicious User                                                   |
| ----------------------------- | --------------------------------------------------- | ---------------------------------------------------------------- | 
| **`阻断 Web 攻击`**            | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-attack-detected.png" width=270 /> |
| **`限制访问频率`**             | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-access-too-fast.png" width=270 /> |
| **`人机验证`**                 | <img src="./images/captcha-1.gif" width=270 />      | <img src="./images/captcha-2.gif" width=270 />                     |
| **`身份认证`**                 | <img src="./images/auth-1.gif" width=270 />         | <img src="./images/auth-2.gif" width=270 />                        |
| **`HTML 动态防护`**            | <img src="./images/dynamic-html-1.png" width=270 /> | <img src="./images/dynamic-html-2.png" width=270 />              |
| **`JS 动态防护`**              | <img src="./images/dynamic-js-1.png" width=270 />   | <img src="./images/dynamic-js-2.png" width=270 />                | 

## 🚀 上手指南

#### 📦 安装

查看 [安装雷池](https://docs.waf-ce.chaitin.cn/zh/%E4%B8%8A%E6%89%8B%E6%8C%87%E5%8D%97/%E5%AE%89%E8%A3%85%E9%9B%B7%E6%B1%A0)

#### ⚙️ 配置防护站点

查看 [快速配置](https://docs.waf-ce.chaitin.cn/zh/%E4%B8%8A%E6%89%8B%E6%8C%87%E5%8D%97/%E5%BF%AB%E9%80%9F%E9%85%8D%E7%BD%AE)

## 📋 更多信息

#### 防护效果测试

| Metric            | ModSecurity, Level 1 | CloudFlare           | 雷池, 平衡            | 雷池, 严格             |
| ----------------- | -------------------- | -------------------- | ---------------------- | --------------------- |
| 样本数量          | 33669                | 33669                | 33669                  | 33669                 |
| **检出率**        | 69.74%               | 10.70%               | 71.65%                 | **76.17%**            |
| **误报率**        | 17.58%               | 0.07%                | **0.07%**              | 0.22%                 |
| **准确率**        | 82.20%               | 98.40%               | **99.45%**             | 99.38%                |


#### 雷池可以投入生产使用吗

是的，已经有不少用户将雷池投入生产使用，截至目前

- 全球累计装机量已超过 18 万台
- 防护的网站数量超过 100 万个
- 每天清洗 HTTP 请求超过 300 亿次

#### 🙋‍♂️ 用户社区

欢迎加入雷池 [社区微信群](/images/wechat.png) 进行技术交流。

也可以加入雷池 [Discord](https://discord.gg/SVnZGzHFvn) 来获取更多社区支持。

<p align="left">
  <a target="_blank" href="/images/wechat.png"><img src="https://img.shields.io/badge/WeChat-07C160?style=flat&logo=wechat&logoColor=white"></a>
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn"><img src="https://img.shields.io/badge/Discord-5865F2?style=flat&logo=discord&logoColor=white"></a> &nbsp;
  <a target="_blank" href="https://x.com/safeline_waf"><img src="https://img.shields.io/badge/X.com-000000?style=flat&logo=x&logoColor=white"></a> &nbsp;
</p>

#### 💪 专业版

查看 [社区版 vs 专业版](https://waf-ce.chaitin.cn/version)
