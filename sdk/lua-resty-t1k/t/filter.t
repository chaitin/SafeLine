use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 5);

run_tests();

__DATA__

=== TEST 1: do_header_filter
--- http_config eval: $::HttpConfig
--- config
    location /t {
        access_by_lua_block {
            ngx.ctx.t1k_extra_header = "k1:v1\nk2:v2\nk3:v3\n"
        }

        header_filter_by_lua_block {
            local filter = require "resty.t1k.filter"
            filter.do_header_filter()
        }

        content_by_lua_block {
            ngx.say("hi")
        }
    }
--- request
GET /t
--- response_headers
k1: v1
k2: v2
k3: v3
--- no_error_log
[error]
