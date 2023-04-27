---
title: "升级雷池"
category: "上手指南"
weight: 6
---

# 升级雷池

**WARN: 雷池 SafeLine 服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。**

```shell
# 请到 compose.yaml 同级目录下执行下面脚本
curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash
```
**有部分环境的默认 SafeLine 安装路径是在 `/data/safeline-ce`，安装之后可能会发现需要重新绑定 OTP、配置丢失等情况，可以修改 .env 的 `SAFELINE_DIR` 变量，指向 `/data/safeline-ce`**