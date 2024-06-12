<h1 align="center">SafeLine, The Best Free WAF For Webmaster</h1>

<p align="center">
  <img src="https://raw.githubusercontent.com/chaitin/SafeLine/main/documents/static/images/403.svg" width="160">
</p>
<br>
<p align="center">
  <img src="https://img.shields.io/badge/SafeLine-BEST_WAF-blue">
  <img src="https://img.shields.io/github/release/chaitin/safeline.svg?color=blue" />
  <img src="https://img.shields.io/github/release-date/chaitin/safeline.svg?color=blue&label=update" />
  <img src="https://img.shields.io/docker/v/chaitin/safeline-mgt-api?color=blue">
  <img src="https://img.shields.io/github/stars/chaitin/safeline?style=social">
</p>

<p align="center">
  <a target="_blank" href="https://waf.chaitin.com/">Home</a> | 
  <a target="_blank" href="https://demo.waf.chaitin.com:9443/dashboard">Demo</a> | 
  <a target="_blank" href="https://docs.waf.chaitin.com/">Docs</a> | 
  <a target="_blank" href="https://waf-ce.chaitin.cn/">中文版</a>
</p>

<p align="center">
  <a target="_blank" href="https://discord.gg/wyshSVuvxC">Discord</a> |
  <a target="_blank" href="https://x.com/safeline_waf">X</a> |
  <a target="_blank" href="https://t.me/safeline_waf">Telegram</a> |
  <a target="_blank" href="/documents/static/images/wechat-230825.png">WeChat</a>
</p>

SafeLine is a simple, lightweight, locally deployable WAF, it is the best waf for webmaster.

It serves as a reverse proxy access to protect your website from network attacks that including OWASP attacks, zero-day attacks, web crawlers, vulnerability scanning, vulnerability exploit, http flood and so on.

- Cumulative installations exceed **130,000** units
- Protecting websites over **1,000,000**
- Processing HTTP requests over **30,000,000,000** times per day
- Intercepting attacks over **50,000,000** times per day

<img src="./images/safeline_en.png" />


## Installation

**中国大陆用户安装国际版可能会导致无法连接云服务，请查看 [中文版安装文档](https://waf-ce.chaitin.cn/docs/guide/install)**

> Recommended

Use the following command to start the automated installation of SafeLine. (This process requires root privileges)

```bash
bash -c "$(curl -fsSLk https://waf.chaitin.com/release/latest/setup.sh)"
```

After the command is executed, it means the installation is successfully. Please go to "Use Web UI" directly.


## Mannually Deploy

to see [Documentation](https://docs.waf.chaitin.com/en/tutorials/install)

## Use Web UI

Open the web console page `https://<safeline-ip>:9443/` in the browser, then you will see below.


<img width="400" src="/images/login.png">

Execute the following command to get administrator account

```bash
docker exec safeline-mgt /app/mgt-cli reset-admin --once
```

After the command is successfully executed, you will see the following content

> Please must remember this content

```text
[SafeLine] Initial username：admin
[SafeLine] Initial password：**********
[SafeLine] Done
```

Enter the password in the previous step and you will successfully logged into SafeLine.

## Protecting a website

### How SafeLine works

SafeLine is a web application firewall developed based on nginx, designed to help websites defend against network attacks.

Its principle is to act as an http/https reverse proxy, receive network traffic for the original website, then clean the malicious attack traffic and forward the safe and reliable traffic to the original website.

<img src="/images/safeline-as-proxy.png" width=400>

### Proxy a website in SafeLine

Log into the SafeLine Web Admin Console, go to the "Site" -> "Website" page and click the "Add Site" button in the upper right corner.

<img src="/images/add-site-1.png" width=800>

In the next dialog box, enter the information to the original website.    

- **Domain**: domain name of your original website, or hostname, or ip address, for example: `www.chaitin.com`
- **Port**: port that SafeLine will listen, such as 80 or 443. (for `https` websites, please check the `SSL` option)
- **Upstream**: real address of your original website, through which SafeLine will forward traffic to it

After completing the above settings, please resolve the domain name you just entered to the IP address of the server where SafeLine is located.

<img src="/images/add-site-2.png" width=400>

Then you can access the website protected by the SafeLine through the domain name like this.

<img src="/images/safeline-as-proxy-2.png" width=400>

## Try to attack your website

Now, your website is protected by SafeLine, let’s try tp attack it and see what happens.

If https://chaitin.com is a website protected by SafeLine, here are some test cases for common attacks:

- SQL Injection: `https://chaitin.com/?id=1+and+1=2+union+select+1`
- XSS: `https://chaitin.com/?id=<img+src=x+onerror=alert()>`
- Path Traversal: `https://chaitin.com/?id=../../../../etc/passwd`
- Code Injection: `https://chaitin.com/?id=phpinfo();system('id')`
- XXE: `https://chaitin.com/?id=<?xml+version="1.0"?><!DOCTYPE+foo+SYSTEM+"">`

Replace `chaitin.com` in the above cases with your website domain name and try to access it.

<img src="/images/blocked.png" width=400>

Check the web console of SafeLine to see the attack list

<img src="/images/log-list.png" width=800>

To view the specific details of the attack, click "detail"

<img src="/images/log-detail.png" width=600>

## Core Capabilities

#### Defenses For OWASP Attacks

SafeLine use as an important tool to defense against OWASP Top 10 Attack, such as SQL injection, XSS, Insecure deserialization etc.

#### Defenses For 0-Day Attacks

SafeLine use intelligent rule-free detection algorithm to against 0-Day attacks with unknown attack signatures.

#### Proactive Bot defense

SafeLine uses advanced algorithms to send capthcha challenge for suspicious users to against automated robot attacks.

#### In-Browser Code Encryption

SafeLine can dynamically encrypt and obfuscate static code in the browser (such as HTML, JavaScript) to against reverse engineering.

#### Web Authentication

SafeLine prompting the user for authentication to web apps that lacks valid authentication credentials, Illegal users will be blocked.

#### Web Access Control List

SafeLine offering fine-grained control over traffic allows you to define a set of rules that determine which requests are allowed or denied.

## Features

#### Easy To Use

Deployed by Docker, one command can complete the installation, and you can get started at 0 cost.

The security configuration is ready to use, no manual maintenance is required, and safe lying management can be achieved.

#### High Security Efficacy

The first intelligent semantic analysis algorithm in the industry, accurate detection, low false alarm, and difficult to bypass.

The semantic analysis algorithm has no rules, and you are no longer at a loss when facing 0-day attacks with unknown features.

#### High Performance

Ruleless engine, linear security detection algorithm, average request detection delay at 1 millisecond level.

Strong concurrency, single core easily detects 2000+ TPS, as long as the hardware is strong enough, there is no upper limit to the traffic scale that can be supported.

#### High Availability

The traffic processing engine is developed based on Nginx, and both performance and stability can be guaranteed.

Built-in complete health check mechanism, service availability is as high as 99.99%.

## Star History <a name="star-history"></a>

<a href="https://github.com/chaitin/safeline/stargazers">
    <img width="500" alt="Star History Chart" src="https://api.star-history.com/svg?repos=chaitin/safeline&type=Date">
</a> 

## Related Repo
<p >
  <a href="https://github.com/chaitin/yanshi">Automaton Generator</a> | 
  <a href="https://github.com/chaitin/safeline-open-platform">Lua Plugin</a> | 
  <a href="https://github.com/chaitin/lua-resty-t1k">T1K Protocol</a> |
  <a href="https://github.com/chaitin/blazehttp">WAF Test Tool</a>
</p>
