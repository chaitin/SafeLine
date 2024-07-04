use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3 + 1);

run_tests();

__DATA__

=== TEST 1: do_access nil option
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local ok, err, _ = t1k.do_access(nil)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
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
[debug]
--- log_level: debug



=== TEST 2: do_access disabled
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "off",
            }

            local ok, err, _ = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
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
--- log_level: debug



=== TEST 3: do_access invalid mode
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "invalid",
            }

            local ok, err, _ = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
            end
        }

        content_by_lua_block {
            ngx.say("passed")
        }
    }
--- request
GET /t/shell.php
--- response_body
passed
--- error_log
lua-resty-t1k: invalid t1k mode: invalid



=== TEST 4: do_access invalid host
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block"
            }

            local ok, err, _ = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
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
lua-resty-t1k: invalid t1k host: nil
--- log_level: debug



=== TEST 5: do_access invalid port
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            local t1k = require "resty.t1k"

            local t = {
                mode = "block",
                host = "127.0.0.1"
            }

            local ok, err, _ = t1k.do_access(t)
            if not ok then
                ngx.log(ngx.ERR, err)
                return
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
lua-resty-t1k: invalid t1k port: nil
--- log_level: debug
