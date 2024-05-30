---
title: "Ingres-Nginx 集成雷池"
---

# 使用方式

### 版本要求
* Safeline >= 5.6.0

### 准备工作
参考 [APISIX 联动雷池](/docs/practice/apisix#准备工作) 的准备工作。

### 准备 Safeline 配置

使用 ConfigMap 配置 Safeline 插件需要的检测引擎 host 和 port，内容如下：

```yaml
# safeline.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: safeline
  namespace: ingress-nginx
data:
    host: "detector_host" # 雷池检测引擎的地址, 参考准备工作
    port: "8000"
```

使用 ingress-nginx 创建 ConfigMap：

```shell
kubectl create namespace ingress-nginx
kubectl apply -f safeline.yaml
```

### 全新安装集成方式

使用 helm 安装 Ingress-Nginx，参考 [Ingress-Nginx 官方文档](https://kubernetes.github.io/ingress-nginx/deploy/#using-helm)

使用下面的 values.yaml 进行镜像替换和插件配置：

```yaml
# values.yaml
controller:
  kind: DaemonSet
  image:
    registry: docker.io
    image: chaitin/ingress-nginx-controller
    tag: v1.10.1
    digest: ""
  extraEnvs:
    - name: SAFELINE_HOST
      valueFrom:
        configMapKeyRef:
          name: safeline
          key: host
    - name: SAFELINE_PORT
      valueFrom:
        configMapKeyRef:
          name: safeline
          key: port
  service:
    externalTrafficPolicy: Local # 便于获取真实客户端 IP
  config:
    plugins: safeline
  admissionWebhooks:
    patch:
      image:
        registry: docker.io
        image: chaitin/ingress-nginx-kube-webhook-certgen
        tag: v1.4.1
        digest: ""
```
执行安装命令
```shell
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
  -f values.yaml
```
如果你想自行构建 Ingress-Nginx 镜像，可以参考下面的 Dockerfile：

```Dockerfile
FROM registry.k8s.io/ingress-nginx/controller:v1.10.1

USER root

RUN apk add --no-cache make gcc unzip wget

# install luaroncks
RUN wget https://luarocks.org/releases/luarocks-3.11.0.tar.gz && \
    tar zxpf luarocks-3.11.0.tar.gz && \
    cd luarocks-3.11.0 && \
    ./configure && \
    make && \
    make install && \
    cd .. && \
    rm -rf luarocks-3.11.0 luarocks-3.11.0.tar.gz

RUN luarocks install ingress-nginx-safeline && \
    ln -s /usr/local/share/lua/5.1/safeline /etc/nginx/lua/plugins/safeline

USER www-data
```

### 已有 Ingress-Nginx 集成方式

#### Step 1. 安装 Safeline 插件

参考上面的 Dockerfile 使用 luarocks 安装 Safeline 插件到 nginx 的默认插件目录。

#### Step 2. 配置 Safeline 插件

使用上面的 safeline.yaml 创建 ConfigMap：

```shell
kubectl apply -f safeline.yaml
```
在 Ingress-Nginx 的插件配置中启用 Safeline ：

```yaml
# ingress-nginx-controller-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: ingress-nginx-controller
  namespace: ingress-nginx
data:
  plugins: "safeline"
```

#### Step 3. 注入 Safeline 环境变量
在 Ingress-Nginx 的 Deployment/DaemonSet 中添加环境变量便于插件读取：

```yaml
# ingress-nginx-controller-deployment.yaml
...
env:
  - name: SAFELINE_HOST
    valueFrom:
      configMapKeyRef:
        name: safeline
        key: host
  - name: SAFELINE_PORT
    valueFrom:
      configMapKeyRef:
        name: safeline
        key: port
```

#### Step 4. [可选] 获取真实客户端 IP

配置 nginx service 的 externalTrafficPolicy 为 Local，以便获取真实客户端 IP

 ### 测试 Safeline 插件

通过构造恶意请求测试 Safeline 插件是否生效，例如：

```shell
curl http://localhost:80/ -H "Host: example.com" -H "User-Agent: () { :; }; echo; echo; /bin/bash -c 'echo hello'"
```