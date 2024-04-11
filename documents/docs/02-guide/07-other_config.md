---
title: "配置其他"
---

# 配置其他

其他配置项介绍

## 防护站点

### 证书管理

管理需要使用的证书，点击添加证书添加

### 代理配置

#### 源 ip 获取方式

1.使用默认的方式获取源 IP

2.自定义获取源 IP 的 header 

#### 其他代理配置

1.如果配置站点需要 http 自动转为 https 功能时，需要手动开启

2.支持使用 HTTP2

3.雷池支持开启 IPv6 的访问

4.代理增加信息，方便数据分析

注：开启后并不会遵循源请求的信息，雷池会覆盖，为防止客户端伪造

## 防护配置

### 黑白名单

黑名单：拦截

白名单：放通，优先级大于黑名单

注意：条件 AND 是指同时符合，如果希望多个匹配条件需要增加多条黑名单或者白名单

![Alt text](/images/docs/guide_config/other_config1.png)

### 频率限制

通过开启频率限制功能封锁恶意 IP

对已经限制的ip可以在限频日志页面点击解除限制进行放行

![Alt text](/images/docs/guide_config/other_config2.png)

### 人机验证

人机验证的有效时间默认是一个小时，未来可能会支持配置，敬请期待。

详情查看 [人机验证 2.0](/about/challenge)

### 语义分析

详情查看 [语义分析检测算法](/about/syntaxanalysis)

### 补充规则（专业版）

补充规则能在语义分析的基础上，针对一些特殊的业务漏洞、框架漏洞的利用行为进行防护。

社区版默认进行平衡防护，专业版可进一步配置防护模式。

![Alt text](/images/docs/guide_config/other_config3.png)

### 身份认证

可以通过添加认证规则，对雷池保护的站点额外增加身份认证功能。

![Alt text](/images/docs/guide_config/other_config4.png)

如图，触发身份认证规则后需要使用账户密码登录后继续访问网站。

![Alt text](/images/docs/guide_config/other_config5.png)

### 通用配置

#### IP 组

1.支持自定义 IP 组

2.长亭社区恶意 IP 情报，需要加入 IP 情报共享计划才可以使用


#### 拦截页面

1.自定义拦截页面的提示信息

2.自定义拦截页面（专业版），可以修改替或换拦截页面

#### 攻击告警（专业版）

支持钉钉、飞书、企业微信

支持发送（攻击、频率限制）告警到钉钉

#### IP 情报共享计划

默认加入共享计划，加入后将共享攻击 IP 信息到社区，并可使用 IP 组 “长亭社区恶意 IP 情报”。

## 系统设置

#### 雷池控制台登录设置

用于配置登录雷池管理端的方式

低于5.0.0版本升级上来的，shell会显示初始密码，忘记可以手动重置

社区版支持单用户，**专业版**支持多用户管理

管理员固定为 admin，非管理员不能修改其他用户配置

#### 雷池控制台证书

存放默认证书，可以自定义证书

#### Syslog 设置

让雷池发送syslog到设置的服务器，**当前只支持UDP协议**

![Alt text](/images/docs/guide_config/other_config6.png)

保存信息后可以点击测试按钮测试，收到测试信息表示配置成功

雷池发现攻击事件后，会发送事件的syslog信息

![Alt text](/images/docs/guide_config/other_config7.png)

```
#syslog 内容详情
{
  "scheme": "http",                 // 请求协议为 HTTP
  "src_ip": "12.123.123.123",       // 源 IP 地址
  "src_port": 53008,                // 源端口号
  "socket_ip": "10.2.71.103",       // Socket IP 地址
  "upstream_addr": "10.2.34.20",    // 上游地址
  "req_start_time": 1712819316749,  // 请求开始时间
  "rsp_start_time": null,           // 响应开始时间
  "req_end_time": 1712819316749,    // 请求结束时间
  "rsp_end_time": null,             // 响应结束时间
  "host": "safeline-ce.chaitin.net",// 主机名
  "method": "GET",                  // 请求方法为 GET
  "query_string": "",               // 查询字符串
  "event_id": "32be0ce3ba6c44be9ed7e1235f9eebab",            // 事件 ID
  "session": "",                    // 会话
  "site_uuid": "35",                // 站点 UUID
  "site_url": "http://safeline-ce.chaitin.net:8083",         // 站点 URL
  "req_detector_name": "1276d0f467e4",                       // 请求检测器名称
  "req_detect_time": 286,           // 请求检测时间
  "req_proxy_name": "16912fe30d8f", // 请求代理名称
  "req_rule_id": "m_rule/9bf31c7ff062936a96d3c8bd1f8f2ff3",  // 请求规则 ID
  "req_location": "urlpath",        // 请求位置为 URL 路径
  "req_payload": "",                // 请求负载为空
  "req_decode_path": "",            // 请求解码路径
  "req_rule_module": "m_rule",      // 请求规则模块为 m_rule
  "req_http_body_is_truncate": 0,   // 请求 HTTP 主体
  "rsp_http_body_is_truncate": 0,   // 响应 HTTP 主体
  "req_skynet_rule_id_list": [      // 请求 Skynet 规则 ID 列表
    65595,
    65595
  ],
  "http_body_is_abandoned": 0,      // HTTP 主体
  "country": "US",                  // 国家
  "province": "",                   // 省份
  "city": "",                       // 城市
  "timestamp": 1712819316,          // 时间戳
  "payload": "",  
  "location": "urlpath",            // 位置为 URL 路径
  "rule_id": "m_rule/9bf31c7ff062936a96d3c8bd1f8f2ff3",     / 规则 ID
  "decode_path": "",                // 解码路径
  "cookie": "sl-session=Z0WLa8mjGGZPki+QHX+HNQ==",          // Cookie
  "user_agent": "PostmanRuntime/7.28.4",                    // 用户代理
  "referer": "",                    // 引用页
  "timestamp_human": "2024-04-11 15:08:36",                 // 时间戳
  "resp_reason_phrase": "",         // 响应
  "module": "m_rule",               // 模块为 m_rule
  "reason": "",                     // 原因
  "proxy_name": "16912fe30d8f",     // 代理名称
  "node": "1276d0f467e4",           // 节点
  "dest_port": 8083,                // 目标端口号
  "dest_ip": "10.2.34.20",          // 目标 IP 地址
  "urlpath": "/webshell.php",       // URL 路径
  "protocol": "http",               // 协议为 HTTP
  "attack_type": "backdoor",        // 攻击类型
  "risk_level": "high",             // 风险级别
  "action": "deny",                 // 动作
  "req_header_raw": "GET /webshell.php HTTP/1.1\r\nHost: safeline-ce.chaitin.net:8083\r\nUser-Agent: PostmanRuntime/7.28.4\r\nAccept: */*\r\nAccept-Encoding: gzip, deflate, br\r\nCache-Control: no-cache\r\nCookie: sl-session=Z0WLa8mjGGZPki+QHX+HNQ==\r\nPostman-Token: 8e67bec1-6e79-458c-8ee5-0498f3f724db\r\nX-Real-Ip: 12.123.123.123\r\nSL-CE-SUID: 35\r\n\r\n",                      // 请求头原始内容
  "body": "",                       // 主体
  "req_block_reason": "web",        // 请求阻止原因
  "req_attack_type": "backdoor",    // 请求攻击类型
  "req_risk_level": "high",         // 请求风险级别
  "req_action": "deny"              // 动作
}
```

#### 系统信息

显示系统版本和设备机器码

## 常见配置问题

请参考 [其他问题](/faq/other)
