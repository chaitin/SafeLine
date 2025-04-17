<p align="center">
  <img src="/images/banner.png" width="400" />
</p>

<h4 align="center">
  SafeLine - Make your web apps secure
</h4>

<p align="center">
  <a target="_blank" href="https://ly.safepoint.cloud/laA8asp">🏠 Website</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://ly.safepoint.cloud/w2AeHhb">📖 Docs</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://ly.safepoint.cloud/hSMd4SH">🔍 Live Demo</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn">🙋‍♂️ Discord</a> &nbsp; | &nbsp;
  <a target="_blank" href="/README_CN.md">中文版</a>
</p>

## 👋 INTRODUCTION

SafeLine is a self-hosted **`WAF(Web Application Firewall)`** to protect your web apps from attacks and exploits.

A web application firewall helps protect web apps by filtering and monitoring HTTP traffic between a web application and the Internet. It typically protects web apps from attacks such as `SQL injection`, `XSS`, `code injection`, `os command injection`, `CRLF injection`, `ldap injection`, `xpath injection`, `RCE`, `XXE`, `SSRF`, `path traversal`, `backdoor`, `bruteforce`, `http-flood`, `bot abused`, among others.

#### 💡 How It Works

<img src="/images/how-it-works.png" width="800" />

By deploying a WAF in front of a web application, a shield is placed between the web application and the Internet. While a proxy server protects a client machine’s identity by using an intermediary, a WAF is a type of reverse-proxy, protecting the server from exposure by having clients pass through the WAF before reaching the server.

A WAF protects your web apps by filtering, monitoring, and blocking any malicious HTTP/S traffic traveling to the web application, and prevents any unauthorized data from leaving the app. It does this by adhering to a set of policies that help determine what traffic is malicious and what traffic is safe. Just as a proxy server acts as an intermediary to protect the identity of a client, a WAF operates in similar fashion but acting as a reverse proxy intermediary that protects the web app server from a potentially malicious client.

its core capabilities include:

- Defenses for web attacks
- Proactive bot abused defense 
- HTML & JS code encryption
- IP-based rate limiting
- Web Access Control List

#### ⚡️ Screenshots

| <img src="./images/screenshot-1.png" width=370 /> | <img src="./images/screenshot-2.png" width=370 /> |
| ------------------------------------------------- | ------------------------------------------------- | 
| <img src="./images/screenshot-3.png" width=370 /> | <img src="./images/screenshot-4.png" width=370 /> | 

Get [Live Demo](https://demo.waf.chaitin.com:9443/)

## 🔥 FEATURES

List of the main features as follows:

- **`Block Web Attacks`**
  - It defenses for all of web attacks, such as `SQL injection`, `XSS`, `code injection`, `os command injection`, `CRLF injection`, `XXE`, `SSRF`, `path traversal` and so on.
- **`Rate Limiting`**
  - Defend your web apps against `DoS attacks`, `bruteforce attempts`, `traffic surges`, and other types of abuse by throttling traffic that exceeds defined limits.
- **`Anti-Bot Challenge`**
  - Anti-Bot challenges to protect your website from `bot attacks`, humen users will be allowed, crawlers and bots will be blocked.
- **`Authentication Challenge`**
  - When authentication challenge turned on, visitors need to enter the password, otherwise they will be blocked.
- **`Dynamic Protection`**
  - When dynamic protection turned on, html and js codes in your web server will be dynamically encrypted by each time you visit.

#### 🧩 Showcases

|                               | Legitimate User                                     | Malicious User                                                   |
| ----------------------------- | --------------------------------------------------- | ---------------------------------------------------------------- | 
| **`Block Web Attacks`**       | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-attack-detected.png" width=270 /> |
| **`Rate Limiting`**           | <img src="./images/skeleton.png" width=270 />       | <img src="./images/blocked-for-access-too-fast.png" width=270 /> |
| **`Anti-Bot Challenge`**       | <img src="./images/captcha-1.gif" width=270 />      | <img src="./images/captcha-2.gif" width=270 />                     |
| **`Auth Challenge`**          | <img src="./images/auth-1.gif" width=270 />         | <img src="./images/auth-2.gif" width=270 />                        |
| **`HTML Dynamic Protection`** | <img src="./images/dynamic-html-1.png" width=270 /> | <img src="./images/dynamic-html-2.png" width=270 />              |
| **`JS Dynamic Protection`**   | <img src="./images/dynamic-js-1.png" width=270 />   | <img src="./images/dynamic-js-2.png" width=270 />                | 

## 🚀 Quickstart

> [!WARNING]
> 中国大陆用户安装国际版可能会导致无法连接云服务，请查看 [中文版安装文档](https://docs.waf-ce.chaitin.cn/zh/%E4%B8%8A%E6%89%8B%E6%8C%87%E5%8D%97/%E5%AE%89%E8%A3%85%E9%9B%B7%E6%B1%A0)

#### 📦 Installing

Information on how to install SafeLine can be found in the [Install Guide](https://docs.waf.chaitin.com/en/GetStarted/Deploy)

#### ⚙️ Protecting Web Apps

to see [Configuration](https://docs.waf.chaitin.com/en/GetStarted/AddApplication)

## 📋 More Informations

#### Effect Evaluation

| Metric            | ModSecurity, Level 1 | CloudFlare, Free     | SafeLine, Balance      | SafeLine, Strict      |
| ----------------- | -------------------- | -------------------- | ---------------------- | --------------------- |
| Total Samples     | 33669                | 33669                | 33669                  | 33669                 |
| **Detection**     | 69.74%               | 10.70%               | 71.65%                 | **76.17%**            |
| **False Positive**| 17.58%               | 0.07%                | **0.07%**              | 0.22%                 |
| **Accuracy**      | 82.20%               | 98.40%               | **99.45%**             | 99.38%                |


#### Is SafeLine Production-Ready?

Yes, SafeLine is production-ready.

- Over 180,000 installations worldwide
- Protecting over 1,000,000 Websites
- Handling over 30,000,000,000 HTTP Requests Daily

#### 🙋‍♂️ Community

Join our [Discord](https://discord.gg/SVnZGzHFvn) to get community support, the core team members are identified by the STAFF role in Discord.

- channel [#feedback](https://discord.com/channels/1243085666485534830/1243120292822253598): for new features discussion.
- channel [#FAQ](https://discord.com/channels/1243085666485534830/1263761679619981413): for FAQ.
- channel [#general](https://discord.com/channels/1243085666485534830/1243115843919806486): for any other questions.

Several contact options exist for our community, the primary one being Discord. These are in addition to GitHub issues for creating a new issue.

<p align="left">
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn"><img src="https://img.shields.io/badge/Discord-5865F2?style=flat&logo=discord&logoColor=white"></a> &nbsp;
  <a target="_blank" href="https://x.com/safeline_waf"><img src="https://img.shields.io/badge/X.com-000000?style=flat&logo=x&logoColor=white"></a> &nbsp;
  <a target="_blank" href="/images/wechat.png"><img src="https://img.shields.io/badge/WeChat-07C160?style=flat&logo=wechat&logoColor=white"></a>
</p>

#### 💪 PRO Edition

Coming soon!

#### 📝 License

See [LICENSE](/LICENSE.md) for details.
