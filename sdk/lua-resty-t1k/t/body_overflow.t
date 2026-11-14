use Test::Nginx::Socket 'no_plan';

run_tests();

__DATA__

=== TEST 1: body above one megabyte is rejected before upstream
--- http_config
    lua_package_path "lib/?.lua;;";
--- config
    location /t {
        client_max_body_size 2m;
        client_body_buffer_size 2m;
        access_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_access({ mode = "block", host = "127.0.0.1", port = 18000,
                req_body_size = 1024 }, true)
        }
        content_by_lua_block { ngx.say("upstream") }
    }
--- request eval
"POST /t\n" . ("a" x (1024 * 1024 + 1))
--- error_code: 413
--- response_body_unlike: upstream


=== TEST 5: overflow rejection is independent of detector mode
--- http_config
    lua_package_path "lib/?.lua;;";
--- config
    location /t {
        client_max_body_size 2m;
        access_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_access({ mode = "monitor", host = "127.0.0.1", port = 18000,
                req_body_size = 0.0625 }, true)
        }
        content_by_lua_block { ngx.say("upstream") }
    }
--- request eval
"POST /t\n" . ("a" x 65)
--- error_code: 413
--- response_body_unlike: upstream


=== TEST 2: body above a smaller configured limit is rejected
--- http_config
    lua_package_path "lib/?.lua;;";
--- config
    location /t {
        client_max_body_size 2m;
        client_body_buffer_size 1k;
        access_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_access({ mode = "block", host = "127.0.0.1", port = 18000,
                req_body_size = 0.0625 }, true)
        }
        content_by_lua_block { ngx.say("upstream") }
    }
--- request eval
"POST /t\n" . ("a" x 2048)
--- error_code: 413
--- response_body_unlike: upstream


=== TEST 3: configured full inspection sends the complete file-backed body
--- http_config
    lua_package_path "lib/?.lua;;";
--- config
    location /t {
        client_max_body_size 2m;
        client_body_buffer_size 1k;
        content_by_lua_block {
            local t1k = require "resty.t1k"
            local ok, err, result = t1k.do_access({ mode = "block",
                host = "127.0.0.1", port = 18000, req_body_size = 0.0625,
                overflow_action = "inspect" }, false)
            if not ok then
                ngx.log(ngx.ERR, err)
                return ngx.exit(500)
            end
            ngx.say(result.action)
        }
    }
--- tcp_listen: 18000
--- tcp_reply eval
"\xc1\x01\x00\x00\x00."
--- tcp_query eval
qr/a{2048}/
--- tcp_query_auto_timeout: 0.5
--- request eval
"POST /t\n" . ("a" x 2048)
--- response_body
.


=== TEST 4: chunked file-backed body above inspection limit is rejected
--- http_config
    lua_package_path "lib/?.lua;;";
--- config
    location /t {
        client_max_body_size 2m;
        client_body_buffer_size 1k;
        access_by_lua_block {
            local t1k = require "resty.t1k"
            t1k.do_access({ mode = "block", host = "127.0.0.1", port = 18000,
                req_body_size = 0.0625 }, true)
        }
        content_by_lua_block { ngx.say("upstream") }
    }
--- raw_request eval
"POST /t HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\nTransfer-Encoding: chunked\r\n\r\n" .
    sprintf("%x\r\n", 2048) . ("a" x 2048) . "\r\n0\r\n\r\n"
--- error_code: 413
--- response_body_unlike: upstream
