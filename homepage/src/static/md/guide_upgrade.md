---
title: "升级雷池"
group: "上手指南"
order: 6
---

# 升级雷池

## 自动一键更新

**WARN: 雷池 SafeLine 服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。**

```shell
# 请到 compose.yaml 同级目录下执行下面脚本
curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash
```
**有部分环境的默认 SafeLine 安装路径是在 `/data/safeline-ce`，安装之后可能会发现需要重新绑定 OTP、配置丢失等情况，可以修改 .env 的 `SAFELINE_DIR` 变量，指向 `/data/safeline-ce`**

## 手动更新镜像

适用于 docker hub 拉取镜像失败的场景，手动更新镜像，注意还是要执行 `upgrade.sh` 来处理 `.env` 的更新，否则有可能会因为缺少参数而启动失败。

### 1. 在一台能够从 docker hub 拉取镜像的机器上执行

```shell

# 拉取镜像
docker pull chaitin/safeline-tengine:latest
docker pull chaitin/safeline-mgt-api:latest
docker pull chaitin/safeline-mario:latest
docker pull chaitin/safeline-detector:latest
docker pull postgres:15.2

# 打包镜像
docker save -o image.tar chaitin/safeline-tengine:latest chaitin/safeline-mgt-api:latest chaitin/safeline-mario:latest chaitin/safeline-detector:latest postgres:15.2

# 传输到 SafeLine 要部署的目标服务器
# scp image.tar <target-server>:/root/
```
### 2. 在目标服务器 load 镜像

```shell
docker load -i image.tar

curl -kfLsS https://waf-ce.chaitin.cn/release/latest/upgrade.sh | bash
```