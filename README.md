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
  <a target="_blank" href="https://discord.gg/wyshSVuvxC">Discord</a> |
  <a target="_blank" href="https://waf-ce.chaitin.cn/">中文版</a>
</p>

SafeLine is a simple, lightweight, locally deployable WAF, it is the best waf for webmaster.

It serves as a reverse proxy access to protect your website from network attacks that including OWASP attacks, zero-day attacks, web crawlers, vulnerability scanning, vulnerability exploit, http flood and so on.

- Cumulative installations exceed **130,000** units
- Protecting websites over **1,000,000**
- Processing HTTP requests over **30,000,000,000** times per day
- Intercepting attacks over **50,000,000** times per day

<img src="./images/safeline_en.png" />


## Installation

> Recommended

Use the following command to start the automated installation of SafeLine. (This process requires root privileges)

```bash
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh)"
```

After the command is executed, it means the installation is successfully. Please go to "Use Web UI" directly.


## Mannually Deploy

to see [Documentation](https://docs.waf.chaitin.com/en/toturials/install)

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
