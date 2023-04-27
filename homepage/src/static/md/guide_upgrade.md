---
title: "升级雷池"
category: "上手指南"
weight: 6
---

# 升级雷池

**WARN: 雷池 SafeLine 服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。**

```shell
curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash

# 根据环境情况自行使用 `docker compose` 或者 `docker-compose`
docker compose down && docker compose pull && docker compose up -d
```
