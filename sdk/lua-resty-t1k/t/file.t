use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3);

run_tests();

__DATA__

=== TEST 1: read full file
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local file = require "resty.t1k.file"
            local path = ngx.var.document_root .. "/foo.bar"
            local ok, err, content = file.read(path, 2 ^ 15)
            if not ok then
                ngx.say(err)
            end
            ngx.print(table.concat(content))
        }
    }
--- user_files eval
[
    ["foo.bar" => "a" x (2 ** 14) ],
]
--- request
GET /t
--- response_body eval
"a" x (2 ** 14)
--- no_error_log
[error]



=== TEST 2: read partial file
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local file = require "resty.t1k.file"
            local path = ngx.var.document_root .. "/foo.bar"
            local ok, err, content = file.read(path, 2 ^ 13)
            if not ok then
                ngx.say(err)
            end
            ngx.print(table.concat(content))
        }
    }
--- user_files eval
[
    ["foo.bar" => "a" x (2 ** 14) ],
]
--- request
GET /t
--- response_body eval
"a" x (2 ** 13)
--- no_error_log
[error]



=== TEST 3: read negative bytes
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local file = require "resty.t1k.file"
            local path = ngx.var.document_root .. "/foo.bar"
            local ok, err, content = file.read(path, -1)
            if not ok then
                ngx.say(err)
            end
            ngx.print(table.concat(content))
        }
    }
--- user_files
>>> foo.bar
--- request
GET /t
--- response_body eval
""
--- no_error_log
[error]



=== TEST 4: read empty file
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local file = require "resty.t1k.file"
            local path = ngx.var.document_root .. "/foo.bar"
            local ok, err, content = file.read(path, 1)
            if not ok then
                ngx.say(err)
            end
            ngx.print(table.concat(content))
        }
    }
--- user_files
>>> foo.bar
--- request
GET /t
--- response_body eval
""
--- no_error_log
[error]



=== TEST 5: read non-existent file
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local file = require "resty.t1k.file"
            local ok, err, buffer = file.read("/opt/non_existent_file", 0)
            if not ok then
                ngx.say(err)
            end
        }
    }
--- request
GET /t
--- response_body
/opt/non_existent_file: No such file or directory
--- no_error_log
[error]
