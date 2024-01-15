---
title: "升级雷池"
---

# 升级雷池

**注意**: 升级雷池时服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。

[版本更新记录](/about/changelog)

## 在线升级

执行以下命令进行升级，升级不会清除历史数据。

```sh
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/upgrade.sh)"
```

[可选] 执行以下命令删除旧版本 Docker 镜像，释放磁盘空间。

```sh
docker rmi $(docker images | grep "safeline" | grep "none" | awk '{print $3}')
```

> 有部分环境的默认 SafeLine 安装路径是在 `/data/safeline-ce`，安装之后可能会发现需要重新绑定 OTP、配置丢失等情况，可以修改 .env 的 `SAFELINE_DIR` 变量，指向 `/data/safeline-ce`

如果需要使用华为云加速，可使用
```sh
CDN=1 bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/upgrade.sh)"
```
如果需要升级到最新版本流式检测模式，可使用
```sh
STREAM=1 bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/upgrade.sh)"
```

## 离线镜像

适用于 docker hub 拉取镜像失败的场景，手动更新镜像。

```sh
# cd /path/to/safeline

mv compose.yaml compose.yaml.old
wget "https://waf-ce.chaitin.cn/release/latest/compose.yaml" --no-check-certificate -O compose.yaml

wget "https://waf-ce.chaitin.cn/release/latest/seccomp.json" --no-check-certificate -O seccomp.json

sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" ".env"

grep "SAFELINE_DIR" ".env" > /dev/null || echo "SAFELINE_DIR=$(pwd)" >> ".env"
grep "IMAGE_TAG" ".env" > /dev/null || echo "IMAGE_TAG=latest" >> ".env"
grep "MGT_PORT" ".env" > /dev/null || echo "MGT_PORT=9443" >> ".env"
grep "POSTGRES_PASSWORD" ".env" > /dev/null || echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> ".env"
grep "SUBNET_PREFIX" ".env" > /dev/null || echo "SUBNET_PREFIX=172.22.222" >> ".env"
grep "IMAGE_PREFIX" ".env" >/dev/null || echo "IMAGE_PREFIX=chaitin" >>".env"
```

下载 [雷池社区版镜像包](https://demo.waf-ce.chaitin.cn/image.tar.gz) 并传输到需要安装雷池的服务器上，执行以下命令加载镜像

```
docker load -i image.tar.gz
```

执行以下命令替换 Docker 容器

```
docker compose down --remove-orphans
docker compose up -d
```

## 常见升级问题

请参考 [升级问题](/faq/upgrade)
