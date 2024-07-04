use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 2 + 2);

run_tests();

__DATA__

=== TEST 1: err_fmt
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local log = require "resty.t1k.log"
            ngx.log(log.err_fmt("%s - %04d - %.4f", "test", 1, 1))
        }
    }
--- request
GET /t
[error]
--- error_log
lua-resty-t1k: test - 0001 - 1.0000
--- log_level: error



=== TEST 2: warn_fmt
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local log = require "resty.t1k.log"
            ngx.log(log.warn_fmt("%s - %04d - %.4f", "test", 1, 1))
        }
    }
--- request
GET /t
--- no_error_log
[error]
--- error_log
lua-resty-t1k: test - 0001 - 1.0000
--- log_level: warn



=== TEST 3: debug_fmt
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local log = require "resty.t1k.log"
            ngx.log(log.debug_fmt("%s - %04d - %.4f", "test", 1, 1))
        }
    }
--- request
GET /t
--- no_error_log
[error]
--- error_log
lua-resty-t1k: test - 0001 - 1.0000
--- log_level: debug
