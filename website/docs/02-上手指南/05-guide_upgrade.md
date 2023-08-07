---
title: "升级雷池"
---

# 升级雷池

**注意**: 升级雷池时服务会重启，流量会中断一小段时间，根据业务情况选择合适的时间来执行升级操作。

## 在线升级

执行以下命令即可进行升级。

```
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/upgrade.sh)"
```

[可选] 升级成功后, 可以执行以下命令删除旧版本 Docke 镜像, 以释放磁盘空间

```
docker rmi $(docker images | grep "safeline" | grep "none" | awk '{print $3}')
```

> 有部分环境的默认 SafeLine 安装路径是在 `/data/safeline-ce`，安装之后可能会发现需要重新绑定 OTP、配置丢失等情况，可以修改 .env 的 `SAFELINE_DIR` 变量，指向 `/data/safeline-ce`

## 离线镜像

适用于 docker hub 拉取镜像失败的场景，手动更新镜像，注意还是要执行 `upgrade.sh` 来处理 `.env` 的更新，否则有可能会因为缺少参数而启动失败。


下载 [雷池社区版镜像包](http://demo.waf-ce.chaitin.cn/image.tar.gz) 并传输到需要安装雷池的服务器上，执行以下命令加载镜像

```
docker load -i image.tar.gz
```

**重要**, 进入安装雷池的目录   

执行以下命令修改配置文件

```
sed -i "s/IMAGE_TAG=.*/IMAGE_TAG=latest/g" ".env"
grep "SAFELINE_DIR" ".env" > /dev/null || echo "SAFELINE_DIR=$(pwd)" >> ".env"
grep "SUBNET_PREFIX" ".env" > /dev/null || echo "SUBNET_PREFIX=172.22.222" >> ".env"
grep "REDIS_PASSWORD" ".env" > /dev/null || echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> ".env"
```

执行以下命令替换 Docker 容器

```
docker compose down
docker compose up -d
```

OK, 你已经完成了升级
