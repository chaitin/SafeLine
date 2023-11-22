---
title: "安装雷池"
---
# 安装雷池
根据实际情况选择安装方式，支持一键安装

- [环境检测](#环境检测) : 查看环境是否符合安装要求
- [在线安装](#在线安装) : 推荐方式，一行命令完成安装
- [离线安装](#离线安装) : 下载离线安装包，轻松完成安装
- [其他方式安装](#使用牧云助手安装) : 使用牧云助手，点击即可完安装

## 在线安装（推荐）

**_如果服务器可以访问互联网环境，推荐使用该方式_**

复制以下命令执行，即可完成安装

```
bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh)"
```

**若安装失败，请参考 [安装问题](/docs/faq/install)**

### 在线安装演示

<iframe src="//player.bilibili.com/player.html?aid=236214137&bvid=BV1Je411f7hQ&cid=1339309164&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"
style={{ width: '100%', height: '350px' }}
> 
</iframe>


## 离线安装

**_如果服务器不可以访问互联网环境，推荐使用该方式_**

> 离线安装前需完成[环境检测](#环境检测),默认已完成docker环境准备

首先，下载 [雷池社区版镜像包](https://demo.waf-ce.chaitin.cn/image.tar.gz) 并传输到需要安装雷池的服务器上，执行以下命令加载镜像

```
cat image.tar.gz | gzip -d | docker load
```

执行以下命令创建并进入雷池安装目录

```
mkdir -p safeline        # 创建 safeline 目录
cd safeline              # 进入 safeline 目录
```

下载 [编排脚本](https://waf-ce.chaitin.cn/release/latest/compose.yaml) 并传输到 safeline 目录中

执行以下命令，生成雷池运行所需的相关环境变量

```
echo "SAFELINE_DIR=$(pwd)" >> .env
echo "IMAGE_TAG=latest" >> .env
echo "MGT_PORT=9443" >> .env
echo "POSTGRES_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "REDIS_PASSWORD=$(LC_ALL=C tr -dc A-Za-z0-9 </dev/urandom | head -c 32)" >> .env
echo "SUBNET_PREFIX=172.22.222" >> .env
```

执行以下命令启动雷池

```
docker compose up -d
```

**若安装失败，请参考 [安装问题](/docs/faq/install)**

### 离线安装演示

<iframe src="//player.bilibili.com/player.html?aid=236214137&bvid=BV1Je411f7hQ&cid=1339309164&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"
style={{ width: '100%', height: '350px' }}
> 
</iframe>





## 使用牧云助手安装

使用 [牧云主机管理助手](https://collie.chaitin.cn/) 进行一键安装

直接点击应用市场的安装按钮安装

![](/images/docs/guide_install/collie_apps.png)


### 助手安装演示

<iframe src="//player.bilibili.com/player.html?aid=613778738&bvid=BV1sh4y1t7Pk&cid=1134834926&p=1"  scrolling="no" border="0" frameBorder="no" framespacing="0" allowFullScreen="{true}"
style={{ width: '100%', height: '350px' }}
> </iframe>


## 环境检测

最低配置需求
- 操作系统：Linux
- 指令架构：x86_64
- 软件依赖：Docker 20.10.14 版本以上
- 软件依赖：Docker Compose 2.0.0 版本以上
- 最小化环境：1 核 CPU / 1 GB 内存 / 5 GB 磁盘

可以逐行执行以下命令来确认服务器配置

```shell
uname -m                                    # 查看指令架构
docker version                              # 查看 Docker 版本
docker compose version                      # 查看 Docker Compose 版本
docker-compose version                      # 老版本查看Compose 版本
cat /proc/cpuinfo| grep "processor"         # 查看 CPU 信息
free -h                                     # 查看内存信息
df -h                                       # 查看磁盘信息
lscpu | grep ssse3                          # 确认CPU是否支持 ssse3 指令集
```



### 配置检测演示

<iframe src="//player.bilibili.com/player.html?aid=918634668&bvid=BV1Uu4y1L7Ko&cid=1339439164&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"
style={{ width: '100%', height: '350px' }}
></iframe>


## 常见安装问题

请参考 [安装问题](/docs/faq/install)

下一步请参考 [登录雷池](/docs/guide/login)