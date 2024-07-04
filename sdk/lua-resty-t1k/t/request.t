use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3 + 14);

run_tests();

__DATA__

=== TEST 1: do_request blocked
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
            ngx.say(result["status"])
            ngx.say(result["event_id"])
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_body
?
405
c0c039a7c348486eaffd9e2f9846b66b
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 2: do_request blocked with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
            ngx.say(result["status"])
            ngx.say(result["event_id"])
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_body
?
405
c0c039a7c348486eaffd9e2f9846b66b
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 3: do_request passed
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\xc1\x01\x00\x00\x00."
--- request
GET /t
--- response_body
.
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 4: do_request passed with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\xc1\x01\x00\x00\x00."
--- request
GET /t
--- response_body
.
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 5: do_request trim request body
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\xc1\x01\x00\x00\x00."
--- request
GET /t
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
--- response_body
.
--- no_error_log
[error]
--- error_log
lua-resty-t1k: request body is too long: 123 bytes, cut to 64 bytes
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 6: do_request trim request body with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\xc1\x01\x00\x00\x00."
--- request
GET /t
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
--- response_body
.
--- no_error_log
[error]
--- error_log
lua-resty-t1k: request body is too long: 123 bytes, cut to 64 bytes
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 7: do_request refuse connection
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("result: ", result)
        }
    }
--- request
GET /t
--- response_body
result: nil
--- error_log
failed to connect to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 8: do_request refuse connection with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("result: ", result)
        }
    }
--- request
GET /t
--- response_body
result: nil
--- error_log
failed to connect to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 9: do_request timeout
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 100,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("result: ", result)
        }
    }
--- tcp_listen: 18000
--- tcp_reply_delay: 200ms
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_body
result: nil
--- error_log
failed to receive info packet from t1k server 127.0.0.1:18000: timeout
--- log_level: debug



=== TEST 10: do_request timeout with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 100,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("result: ", result)
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply_delay: 200ms
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_body
result: nil
--- error_log
failed to receive info packet from t1k server unix:t1k.sock: timeout
--- log_level: debug



=== TEST 11: do_request invalid action
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("action: ", result["action"])
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\xc1\x01\x00\x00\x00~"
--- request
GET /t
--- response_body
action: ~
--- no_error_log
[error]
--- error_log
successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 12: do_request invalid action with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 0.0625,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say("action: ", result["action"])
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\xc1\x01\x00\x00\x00~"
--- request
GET /t
--- response_body
action: ~
--- no_error_log
[error]
--- error_log
successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 13: do_request remote address
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("ngx.var.http_x_real_ip or ngx.var.remote_addr is ", utils.get_indexed_element(ngx.var.http_x_real_ip) or ngx.var.remote_addr)
            ngx.say("ngx.var.http_x_forwarded_for: 2 or ngx.var.remote_addr is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 2) or ngx.var.remote_addr)
            ngx.say("ngx.var.http_x_forwarded_for: -2 or ngx.var.remote_addr is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -2) or ngx.var.remote_addr)
            ngx.say("ngx.var.http_non_existent_header or ngx.var.remote_addr is ", utils.get_indexed_element(ngx.var.http_non_existent_header) or ngx.var.remote_addr)
        }
    }
--- request
GET /t
--- more_headers
X-Forwarded-For: 1.1.1.1, 2.2.2.2, 2001:db8:3333:4444:5555:6666:7777:8888, 3.3.3.3
X-Real-IP: 100.100.100.100
--- response_body
ngx.var.http_x_real_ip or ngx.var.remote_addr is 100.100.100.100
ngx.var.http_x_forwarded_for: 2 or ngx.var.remote_addr is 2.2.2.2
ngx.var.http_x_forwarded_for: -2 or ngx.var.remote_addr is 2001:db8:3333:4444:5555:6666:7777:8888
ngx.var.http_non_existent_header or ngx.var.remote_addr is 127.0.0.1
--- no_error_log
[error]



=== TEST 14: do_request http2
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
            ngx.say(result["status"])
            ngx.say(result["event_id"])
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- http2
--- request
GET /t/shell.php
--- tcp_query eval
qr/.*HTTP\/2.0.*/
--- response_body
?
405
c0c039a7c348486eaffd9e2f9846b66b
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 15: do_request http2 with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local request = require "resty.t1k.request"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                uds = true,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = request.do_request(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end

            ngx.say(result["action"])
            ngx.say(result["status"])
            ngx.say(result["event_id"])
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- http2
--- request
GET /t/shell.php
--- tcp_query eval
qr/.*HTTP\/2.0.*/
--- response_body
?
405
c0c039a7c348486eaffd9e2f9846b66b
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug
