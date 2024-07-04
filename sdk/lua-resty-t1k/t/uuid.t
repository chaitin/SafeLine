use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3);

run_tests();

__DATA__

=== TEST 1: generate_v4
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local uuid = require "resty.t1k.uuid"
            ngx.say(uuid.generate_v4())
        }
    }
--- request
GET /t
--- response_body_like
^[0-9a-f]{8}[0-9a-f]{4}4[0-9a-f]{3}[89ab][0-9a-f]{3}[0-9a-f]{12}$
--- no_error_log
[error]
