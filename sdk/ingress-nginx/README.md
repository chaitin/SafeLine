# ingress-nginx-safeline

[Ingress-nginx](https://kubernetes.github.io/ingress-nginx/) plugin for Chaitin SafeLine Web Application Firewall (WAF). This plugin is used to protect your API from malicious requests. It can be used to block requests that contain malicious content in the request body, query parameters, headers, or URI.

## Usage

### Step 1: Install the plugin

way 1: Build your own ingress-nginx/controller image with the plugin installed.

```dockerfile
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

way 2: Use the chaitin ingress-nginx-controller image.

replace image ingress-nginx-controller with `docker.io/chaitin/ingress-nginx-controller:v1.10.1` in your deployment.

### Step 2: Configure the plugin

use a ConfigMap to configure the plugin

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: safeline
  namespace: default
data:
  host: "YOUR_DETECTOR_HOST"
  port: "YOUR_DETECTOR_PORT"
```

### Step 3: Configure the ingress-controller

inject env `SAFELINE_HOST` and `SAFELINE_PORT` to the ingress-controller deployment

```yaml
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
...            

```

### Step 3: Enable the plugin

enable safeline plugin in configmap

```yaml
apiVersion: v1
data:
  allow-snippet-annotations: "false"
  plugins: "safeline"
kind: ConfigMap
metadata:
  name: ingress-nginx-controller
  namespace: default
```

### Step 4: Set externalTrafficPolicy to Local
by default, the ingress-nginx-controller service is of type LoadBalancer, which means the source IP of the request will be the IP of the LoadBalancer. If you want to get the real source IP, you can set the externalTrafficPolicy to Local.

### Step 5: Test the plugin

use a simple http sql injection test

```bash
curl -X POST http://localhost/ -d "select * from users where id=1 or 1=1" 
```

you should get a 403 response.
```bash
{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "b53eb5b95796475699c52a019abb8e6a"}
```