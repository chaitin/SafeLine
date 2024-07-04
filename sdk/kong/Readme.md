# Kong Safeline Plugin
Kong plugin for Chaitin SafeLine Web Application Firewall (WAF). This plugin is used to protect your API from malicious requests. It can be used to block requests that contain malicious content in the request body, query parameters, headers, or URI.

## Installation
To install the plugin, run the following command in your Kong server:

```shell
$ luarocks install kong-safeline
```

## Configuration
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
