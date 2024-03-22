---
title: "登录雷池"
---

# 登录雷池

> TOTP (Time-based One-Time Password algorithm) 将密钥与当前时间进行组合，通过哈希算法产生一次性密码，已被采纳为 RFC 6238，被用于许多双因素身份验证系统。

## 登录流程

1.浏览器打开后台管理页面 `https://<waf-ip>:9443`。

2.输入初始的admin密码

完成安装后在shell会自动输出密码。

![Alt text](/images/docs/guide_config/login_1.png)

若忘记查看，需手生执行命令重置获得初始密码

`docker exec safeline-mgt resetadmin`

![Alt text](/images/docs/guide_config/login_2.png)

3.根据界面提示，使用 **支持 TOTP 的认证软件或者小程序** 扫描二维码，然后输入动态口令登录：

<iframe src="//player.bilibili.com/player.html?aid=748637002&bvid=BV1wC4y177zN&cid=1339420830&p=1&autoplay=0" scrolling="no" border="0" frameBorder="no" framespacing="0" allowFullScreen='{true}'
style={{ width: '100%', height: '350px' }}
></iframe>

### 注意事项：

1.服务器和 totp 应用的**时间必须保持一致**，否则无法验证通过

2.完成首次登录后，**无法回退查看二维码**，使用页面提供的方法重置

## 常见登录问题

请参考 [登录问题](/faq/login)
