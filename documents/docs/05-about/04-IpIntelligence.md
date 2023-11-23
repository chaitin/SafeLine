---
title: "免费黑 IP 情报源"
---

# 免费黑 IP 情报源

长亭免费开放 WAF 黑 IP 情报源

## IP 情报赋能 WAF

在过去的几年，各种大小重保活动期间，总会看到甲方的安全工程师们会拉一个群，用于共享攻击源 IP，这些 IP 也是重保防护能力的重要支撑。

威胁情报对于 Web 攻击防护的作用毋庸置疑，可以精准识别 Bot、C2、VPN、僵尸网络，对于防 0Day 攻击、防自动化攻击、降低误报漏报等方面都有着非常不错的效果。

## 情报从哪来

安装雷池社区版后，若选择 “加入 IP 情报共享计划”，雷池将定时自动聚合攻击 IP 数据（只有攻击 IP，不涉及任何敏感信息）发送到长亭百川云平台。

长亭百川云平台收到来自五湖四海的社区版兄弟们共享的 IP 情报，使用算法进行汇总和优选，生成雷池社区黑 IP 库，最终再回馈给社区。

### 使用方式

1. 在 “通用配置” 页面选择 “加入 IP 情报共享计划”：
   ![join_ip_intelligence](/images/docs/join_ip_intelligence.png)

2. IP 组页面中将会出现 “长亭社区恶意 IP 情报”：
   ![malicious_ip_intelligence](/images/docs/malicious_ip_intelligence.png)

3. 在 “黑白名单” 页面增加黑名单规则，条件设置为对应的 IP 组：
   ![ip_intelligence_blacklist](/images/docs/ip_intelligence_blacklist.png)

4. 坐等 IP 情报拦截。
