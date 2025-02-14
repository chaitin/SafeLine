# Kong Safeline Plugin
Kong plugin for Chaitin SafeLine Web Application Firewall (WAF). This plugin is used to protect your API from malicious requests. It can be used to block requests that contain malicious content in the request body, query parameters, headers, or URI.

## Safeline Prepare
The detection engine of the SafeLine provides services by default via Unix socket. We need to modify it to use TCP, so it can be called by the t1k plugin.

1.Navigate to the configuration directory of the SafeLine detection engine:
```shell
cd /data/safeline/resources/detector/
```
2.Open the `detector.yml` file in a text editor. Modify the bind configuration from Unix socket to TCP by adding the following settings:
```yaml
bind_addr: 0.0.0.0
listen_port: 8000
```
These configuration values will override the default settings in the container, making the SafeLine engine listen on port 8000.

3.Next, map the container’s port 8000 to the host machine. First, navigate to the SafeLine installation directory:
```shell
cd /data/safeline
```

4.Open the compose.yaml file in a text editor and add the ports field to the detector container to expose port 8000:
```yaml
...
detect:
  ports:
    - 8000:8000
...
```

5.Save the changes and restart SafeLine with the following commands:
```shell
docker-compose down
docker-compose up -d
```
This will apply the changes and activate the new configuration.

## Plugin Installation
To install the plugin, run the following command in your Kong server:

```shell
$ luarocks install kong-safeline
```

## Plugin Configuration
You can add the plugin to your API by making the following request:

```shell
# if your detector is running on tcp port
$ curl -X POST http://kong:8001/services/{name}/plugins \
    --data "name=safeline" \
    --data "config.host=your_detector_host" \
    --data "config.port=your_detector_port"
# if your detector is running on unix socket
$ curl -X POST http://kong:8001/services/{name}/plugins \
    --data "name=safeline" \
    --data "config.host=unix:/path/to/your/unix/socket"
```

## Test
You can test the plugin by sending a request to your API with malicious content. If the request is blocked, you will receive a `403 Forbidden` response.

```shell
$ curl -X POST http://kong:8000?1=1%20and%202=2

# you will receive a 403 Forbidden response
{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "8b41a021ea9541c89bb88f3773b4da24"}
```
