---
title: "安装雷池"
group: "上手指南"
order: 2
---

# 安装雷池

## 机器运行的最低配置

最低 1C1G 能运行，具体需要多少配置取决于你的业务流量特征，比如 QPS、网络吞吐等等，暂时没有详细的 datasheet 性能参考。

### 1. 确保机器上正确安装 [Docker](https://docs.docker.com/engine/install/) 和 [Compose V2](https://docs.docker.com/compose/install/)

```shell
docker info # >= 20.10.6
docker compose version # >= 2.0.0
```

### 2. 部署安装

```shell
mkdir -p safeline && cd safeline
# 下载并执行 setup
curl -fLsS https://waf-ce.chaitin.cn/release/latest/setup.sh | bash

# 运行
sudo docker compose up -d
```
