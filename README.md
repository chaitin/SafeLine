
# SafeLine, make your web apps secure

<img src="/images/403.svg" align="right" width="200" />

SafeLine is a self-hosted **web application firewall** to protect your web apps from attacks and exploits.

It defenses for all of web attacks, such as sql injection, code injection, os command injection, CRLF injection, ldap injection, xpath injection, rce, xss, xxe, ssrf, path traversal, backdoor, bruteforce, http-flood, bot abused and so on.

<p align="left">
  <a target="_blank" href="https://waf.chaitin.com/">🏠Home</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://docs.waf.chaitin.com/">📖Documentation</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://demo.waf.chaitin.com:9443/dashboard">🔍Live Demo</a> &nbsp; | &nbsp;
  <a target="_blank" href="https://waf-ce.chaitin.cn/">中文版</a>
</p>

<p align="left">
  <a target="_blank" href="https://discord.gg/SVnZGzHFvn"><img src="https://img.shields.io/badge/Discord-5865F2?style=flat&logo=discord&logoColor=white"></a> &nbsp;
  <a target="_blank" href="/images/wechat-230825.png"><img src="https://img.shields.io/badge/WeChat-07C160?style=flat&logo=wechat&logoColor=white"></a>
</p>

# Table of Contents

- [Introduction](#Intruction)
- [Security Features](#Security-Features)
- [Quickstart](#Quickstart)
- [Community](#Community)
- [More Informations](#More-Informations)

# Introduction

SafeLine is a self-hosted **web application firewall** to protect your web apps from attacks and exploits.

It defenses for all of web attacks, such as sql injection, code injection, os command injection, CRLF injection, ldap injection, xpath injection, rce, xss, xxe, ssrf, path traversal, backdoor, bruteforce, http-flood, bot abused and so on.

## How It Works

<img src="/images/safeline-as-proxy.png" align="right" width=400 />

SafeLine is developed based on nginx, it serves as a reverse proxy middleware to detect and cleans web attacks, its core capabilities include:

- Defenses for web attacks
- Proactive bot abused defense 
- HTML & JS code encryption
- IP-based rate limiting
- Web Access Control List

## Screenshots

| <img src="./images/screenshot-1.png" width=600 /> | <img src="./images/screenshot-2.png" width=600 /> |
| ------------------------------------------------ | ------------------------------------------------ | 
| <img src="./images/screenshot-3.png" width=600 /> | <img src="./images/screenshot-4.png" width=600 /> | 



# Security Features

## Web Attacks

SafeLine uses a non-rule detection algorithm based on syntax analysis, and uses the context-free grammar commonly used in programming languages to replace the regular grammar used by traditional WAFs, which greatly improves the accuracy and recall rate of the detection algorithm.


## Rate Limiting

Defend your applications and APIs against abuse by throttling traffic that exceeds defined limits

Rate Limiting protects against denial-of-service attacks, brute force login attempts, traffic surges, and other types of abuse targeting APIs and applications.

Choose IP-based Rate Limiting to protect unauthenticated endpoints, limit the number of requests from specific IP addresses, and handle abuse from repeat offenders.


## Captcha Challenge

CAPTCHA challenges to protect your website from bot attacks, humen users will be allowed, crawlers and bots will be blocked.


## **Authentication Challenge**

when athentication turned on, visitors need to enter the username and password information you configured below, users who do not hold the password will be blocked.


## Dynamic Protection

When dynamic protection turned on, the html and javascript codes in your website will be dynamically encrypted into different random result each time you visit, it could effectively block crawlers and attack automated exploit programs.

After the html code passes through SafeLine's dynamic protection, it will be randomly encrypted and decrypted automatically when used in the browser. Please see the example below.

The left side is before encrypted, and the right side is after encrypted.

![Untitled](https://prod-files-secure.s3.us-west-2.amazonaws.com/2f02162d-df95-4b24-9ef0-39f24f015615/077e72f8-ffd9-461f-9a56-2ad3ec03637b/Untitled.png)

![Untitled](https://prod-files-secure.s3.us-west-2.amazonaws.com/2f02162d-df95-4b24-9ef0-39f24f015615/e4082371-682c-4c66-9d16-59389afe1106/Untitled.png)

## Web ACL


# Quickstart

**中国大陆用户安装国际版可能会导致无法连接云服务，请查看** [中文版安装文档](https://docs.waf-ce.chaitin.cn/zh/%E4%B8%8A%E6%89%8B%E6%8C%87%E5%8D%97/%E5%AE%89%E8%A3%85%E9%9B%B7%E6%B1%A0)

## Installing

Information on how to install SafeLine can be found in the [Install Guide](https://docs.waf.chaitin.com/en/tutorials/install)

## Protecting Web Apps

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

## Attack Simulation

Now, your website is protected by SafeLine, let’s try to attack it and see what happens.

There are some testcases for common attacks:

- SQL Injection: `https://example.com/?id=1+and+1=2+union+select+1`
- XSS: `https://example.com/?id=<img+src=x+onerror=alert()>`
- Path Traversal: `https://example.com/?id=../../../../etc/passwd`
- Code Injection: `https://example.com/?id=phpinfo();system('id')`

Replace `example.com` in the above cases with your website domain name and try to access it. Then you will see that these attacks will be blocked by SafeLine.

# More Informations

## Is SafeLine Production-Ready?

Yes, SafeLine is production-ready.

- Over 180,000 installations worldwide
- Protecting over 1,000,000 Websites
- Handling over 30,000,000,000 HTTP Requests Daily

## Pro Version

## Stargazers Over Time

<a href="https://starchart.cc/chaitin/SafeLine"><img src="https://starchart.cc/chaitin/SafeLine.svg?variant=light" width=800></a>

## Related Repo
<p >
  <a href="https://github.com/chaitin/yanshi">Automaton Generator</a> | 
  <a href="https://github.com/chaitin/safeline-open-platform">Lua Plugin</a> | 
  <a href="https://github.com/chaitin/lua-resty-t1k">T1K Protocol</a> |
  <a href="https://github.com/chaitin/blazehttp">WAF Test Tool</a>
</p>
