use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3);

run_tests();

__DATA__

=== TEST 1: handle passed action
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local handler = require "resty.t1k.handler"

            local t = {
                action = ".",
            }

            local ok, err = handler.handle(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t
--- response_body
passed
--- no_error_log
[error]



=== TEST 2: handle blocked action
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local handler = require "resty.t1k.handler"

            local t = {
                action = "?",
                status = 405,
                event_id = "c0c039a7c348486eaffd9e2f9846b66b",
            }

            local ok, err = handler.handle(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        header_filter_by_lua_block {
            local filter = require "resty.t1k.filter"
            filter.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t
--- response_body
{"code": 405, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "c0c039a7c348486eaffd9e2f9846b66b"}
--- error_code eval
"405"
--- no_error_log
[error]
--- no_error_log
[error]



=== TEST 3: handle unknown action
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local handler = require "resty.t1k.handler"

            local t = {
                action = "~"
            }

            local ok, err = handler.handle(t)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t
--- response_body
passed
--- error_log
lua-resty-t1k: unknown action from t1k server: ~



=== TEST 4: handle nil result
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local handler = require "resty.t1k.handler"

            local ok, err = handler.handle(nil)
            if not ok then
                ngx.log(ngx.ERR, err)
            end
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t
--- response_body
passed
--- error_log
lua-resty-t1k: invalid result type: nil
