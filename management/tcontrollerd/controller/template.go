package controller

var nginxConfigTpl = `
upstream %s {
    %s
    keepalive 128;
    keepalive_timeout 75;
}
server {
    %s
    %s
    %s
    %s
    location = /forbidden_page {
        internal;
        root /etc/nginx/forbidden_pages;
        try_files /default_forbidden_page.html =403;
    }
    location ^~ / {
        proxy_pass %s://%s;
        include proxy_params;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Accept-Encoding "";
        t1k_append_header SL-CE-SUID %d;
        t1k_body_size 1024k;
        tx_body_size 4k;
        t1k_error_page 403 /forbidden_page;
        tx_error_page 403 /forbidden_page;
    }
}`

var upstreamAddrTpl = `server %s;`
var serverListenTpl = `listen 0.0.0.0:%s%s%s;`
var addrAnyPropertiesTpl = " default_server backlog=65536 reuseport"
var serverNameTpl = `server_name %s;`
var certTpl = `ssl_certificate /etc/nginx/certs/%s;`
var certKeyTpl = `ssl_certificate_key /etc/nginx/certs/%s;`
var backendTpl = `backend_%d`
