use Test::Nginx::Socket 'no_plan';

our $MainConfig = <<'_EOC_';
env DETECTOR_IP;
_EOC_

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
    lua_socket_log_errors off;
_EOC_

run_tests();

__DATA__

=== TEST 1: integration test blocked
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_headers
Content-Type: application/json
--- response_body_like eval
'^{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": ".*"}$'
--- error_code: 403
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 2: integration test blocked internal handle
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_headers
Content-Type: application/json
--- response_body_like eval
'^{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": ".*"}$'
--- error_code: 403
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 3: integration test blocked http2
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- http2
--- request
GET /t/shell.php
--- response_headers
Content-Type: application/json
--- response_body_like eval
'^{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": ".*"}$'
--- error_code: 403
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 4: integration test blocked http2 internal handle
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- http2
--- request
GET /t/shell.php
--- response_headers
Content-Type: application/json
--- response_body_like eval
'^{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": ".*"}$'
--- error_code: 403
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 5: integration test monitor
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_body
passed
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 6: integration test monitor internal handle
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_body
passed
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 7: integration test monitor http2
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- http2
--- request
GET /t/shell.php
--- response_body
passed
--- no_error_log
[error]
--- error_log eval
["lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000", "skip blocking"]
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 8: integration test monitor http2 internal handle
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = os.getenv("DETECTOR_IP"),
                port = 8000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- http2
--- request
GET /t/shell.php
--- response_body
passed
--- no_error_log
[error]
--- error_log eval
"lua-resty-t1k: successfully connected to t1k server $ENV{DETECTOR_IP}:8000"
--- log_level: debug
--- skip_eval
4: not exists($ENV{DETECTOR_IP})



=== TEST 9: integration test disabled
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "off",
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: t1k is not enabled
skip blocking
--- log_level: debug



=== TEST 10: integration test configuration priority
--- main_config eval: $::MainConfig
--- http_config eval: $::HttpConfig
--- config
    access_by_lua_block {
        local t1k = require "resty.t1k"

        local t = {
            mode = "block",
            host = os.getenv("DETECTOR_IP"),
            port = 8000,
            connect_timeout = 1000,
            send_timeout = 1000,
            read_timeout = 1000,
            req_body_size = 1024,
            keepalive_size = 16,
            keepalive_timeout = 10000,
        }

        local ok, err, result = t1k.do_access(t)
        if not ok then
            ngx.log(ngx.ERR, err)
            return
        end

        if t.mode ~= "block" then
            ngx.log(ngx.DEBUG, "skip blocking")
            return
        end

        ok, err = t1k.do_handle(result)
        if not ok then
            ngx.log(ngx.ERR, err)
            return
        end
    }

    header_filter_by_lua_block {
        local t1k = require "resty.t1k"
        t1k.do_header_filter()
    }

    location /pass {
        access_by_lua_block {
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }

    location /block {
        content_by_lua_block {
            ngx.say("there must be a problem when you see this line")
        }
    }
--- request eval
["GET /pass/shell.php", "GET /block/shell.php"]
--- response_body_like eval
["passed", '^{"code": 403, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": ".*"}$']
--- error_code eval
[200, 403]
--- no_error_log
[error]
--- skip_eval
6: not exists($ENV{DETECTOR_IP})



=== TEST 11: integration test blocked extra headers
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

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

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_headers
Content-Type: application/json
--- response_body
{"code": 405, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "c0c039a7c348486eaffd9e2f9846b66b"}
--- error_code eval
"405"
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 12: integration test blocked extra headers internal handle
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

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

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_headers
Content-Type: application/json
--- response_body
{"code": 405, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "c0c039a7c348486eaffd9e2f9846b66b"}
--- error_code eval
"405"
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 13: integration test blocked extra headers with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_headers
Content-Type: application/json
--- response_body
{"code": 405, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "c0c039a7c348486eaffd9e2f9846b66b"}
--- error_code eval
"405"
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 14: integration test passed extra headers
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

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

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00.\xa3\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a"
--- request
GET /t
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 15: integration test passed extra headers internal handle
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

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

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00.\xa3\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a"
--- request
GET /t
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 16: integration test passed extra headers with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = "unix:t1k.sock",
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\x41\x01\x00\x00\x00.\xa3\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a"
--- request
GET /t
--- response_headers
k1: v1
k2: v2
k3: v3
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
--- log_level: debug



=== TEST 17: integration test monitor extra headers
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- raw_response_headers_unlike eval
'.*k1: v1\r\n.*'
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
skip blocking
--- log_level: debug



=== TEST 18: integration test monitor extra headers internal handle
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = "127.0.0.1",
                port = 18000,
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t, true)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- raw_response_headers_unlike eval
'.*k1: v1\r\n.*'
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server 127.0.0.1:18000
--- log_level: debug



=== TEST 19: integration test monitor extra headers with unix domain socket
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "monitor",
                host = "unix:t1k.sock",
                connect_timeout = 1000,
                send_timeout = 1000,
                read_timeout = 1000,
                req_body_size = 1024,
                keepalive_size = 16,
                keepalive_timeout = 10000,
            }

            local ok, err, result = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end

            if t.mode ~= "block" then
                ngx.log(ngx.DEBUG, "skip blocking")
                return
            end

            ok, err = t1k.do_handle(result)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        header_filter_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- tcp_listen: t1k.sock
--- tcp_reply eval
"\x41\x01\x00\x00\x00?\x02\x03\x00\x00\x00405\x23\x12\x00\x00\x00k1:v1\x0ak2:v2\x0ak3:v3\x0a\xa4\x33\x00\x00\x00<!-- event_id: c0c039a7c348486eaffd9e2f9846b66b -->"
--- request
GET /t/shell.php
--- raw_response_headers_unlike eval
'.*k1: v1\r\n.*'
--- response_body
passed
--- no_error_log
[error]
--- error_log
lua-resty-t1k: successfully connected to t1k server unix:t1k.sock
skip blocking
--- log_level: debug
