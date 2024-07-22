
# SafeLine, make your web apps secure

<img src="/images/403.svg" align="right" width="200" />

SafeLine is a web security gateway to protect your websites from attacks and exploits.

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
- [Screenshots](#Screenshots)
- [How It Works](.#How-It-Works)
- [Security Features](#Security-Features)
- [Quickstart](#Quickstart)
- [Community](#Community)
- [More Informations](#More-Informations)

# Screenshots

| <img src="./images/safeline_en.png" width=600 /> | <img src="./images/safeline_en.png" width=600 /> |
| ------------------------------------------------ | ------------------------------------------------ | 
| <img src="./images/safeline_en.png" width=600 /> | <img src="./images/safeline_en.png" width=600 /> | 


# How It Works

<img src="/images/safeline-as-proxy.png" align="right" width=400 />

SafeLine is developed based on nginx, it serves as a reverse proxy middleware to detect and cleans web attacks, its core capabilities include:

- Defenses for web attacks
- Proactive bot abused defense 
- HTML & JS code encryption
- IP-based rate limiting
- Web Access Control List

# Security Features


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

If https://chaitin.com is a website protected by SafeLine, here are some test cases for common attacks:

- SQL Injection: `https://chaitin.com/?id=1+and+1=2+union+select+1`
- XSS: `https://chaitin.com/?id=<img+src=x+onerror=alert()>`
- Path Traversal: `https://chaitin.com/?id=../../../../etc/passwd`
- Code Injection: `https://chaitin.com/?id=phpinfo();system('id')`
- XXE: `https://chaitin.com/?id=<?xml+version="1.0"?><!DOCTYPE+foo+SYSTEM+"">`

Replace `chaitin.com` in the above cases with your website domain name and try to access it.

<img src="/images/blocked.png" width=400>

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
