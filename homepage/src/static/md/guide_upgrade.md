---
title: "升级雷池"
category: "上手指南"
weight: 6
---

# 升级雷池

**WARN: 雷池 SafeLine 服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。**

```
# 查看 `IMAGE_TAG`
cat .env | grep IMAGE_TAG
# 把 IMAGE_TAG 修改为 latest 或者某个特定版本，比如 1.1.0
sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" .env

# 检查 `compose.yaml`
# docker 镜像的命名空间调整到了 `chaitin`，部分旧版本配置是使用的 `chaitinops`
sed -i "s/chaitinops/chaitin/g" compose.yaml

# 根据环境情况自行使用 `docker compose` 或者 `docker-compose`
docker compose down && docker compose pull && docker compose up -d
```
