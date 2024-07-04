use Test::Nginx::Socket;

our $HttpConfig = <<'_EOC_';
    lua_package_path "lib/?.lua;/usr/local/share/lua/5.1/?.lua;;";
_EOC_

repeat_each(3);

plan tests => repeat_each() * (blocks() * 3);

run_tests();

__DATA__

=== TEST 1: int_to_char_length
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("255 to char length: ", utils.int_to_char_length(255))
            ngx.print("16777216 to char length: ", utils.int_to_char_length(16777216))
        }
    }
--- request
GET /t
--- response_body eval
"255 to char length: \x{ff}\x{00}\x{00}\x{00}
16777216 to char length: \x{00}\x{00}\x{00}\x{01}"
--- no_error_log
[error]



=== TEST 2: int_to_char_length
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("0xff 0x00 0x00 0x00 to int length: ", utils.char_to_int_length("\xff\x00\x00\x00"))
            ngx.say("0x00 0x00 0x00 0x01 to int length: ", utils.char_to_int_length("\x00\x00\x00\x01"))
        }
    }
--- request
GET /t
--- response_body
0xff 0x00 0x00 0x00 to int length: 255
0x00 0x00 0x00 0x01 to int length: 16777216
--- no_error_log
[error]



=== TEST 3: is_mask_first
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("0x30 is MASK FIRST: ", utils.is_mask_first(0x30))
            ngx.say("0x41 is MASK FIRST: ", utils.is_mask_first(0x41))
            ngx.say("0xc0 is MASK FIRST: ", utils.is_mask_first(0xc0))
            ngx.say("0xc1 is MASK FIRST: ", utils.is_mask_first(0xc1))
        }
    }
--- request
GET /t
--- response_body
0x30 is MASK FIRST: false
0x41 is MASK FIRST: true
0xc0 is MASK FIRST: true
0xc1 is MASK FIRST: true
--- no_error_log
[error]



=== TEST 4: is_mask_last
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("0x41 is MASK LAST: ", utils.is_mask_last(65))
            ngx.say("0x80 is MASK LAST: ", utils.is_mask_last(128))
        }
    }
--- request
GET /t
--- response_body
0x41 is MASK LAST: false
0x80 is MASK LAST: true
--- no_error_log
[error]



=== TEST 5: packet_parser unfinished
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            local finished, tag, length = utils.packet_parser("\x41\x59\x00\x00\x00")
            ngx.say("finished: ", finished)
            ngx.say("tag == TAG_HEAD: ", tag == 1)
            ngx.say("length: ", 89)
        }
    }
--- request
GET /t
--- response_body
finished: false
tag == TAG_HEAD: true
length: 89
--- no_error_log
[error]



=== TEST 6: packet_parser finished
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            local finished, tag, length = utils.packet_parser("\xa0\x08\x00\x00\x00")
            ngx.say("finished: ", finished)
            ngx.say("tag == TAG_VERSION: ", tag == 32)
            ngx.say("length: ", 8)
        }
    }
--- request
GET /t
--- response_body
finished: true
tag == TAG_VERSION: true
length: 8
--- no_error_log
[error]



=== TEST 7: starts_with
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("http://www.baidu.com starts with \"http\": ", utils.starts_with("http://www.baidu.com", "http"))
            ngx.say("http://www.baidu.com starts with \"https\": ", utils.starts_with("http://www.baidu.com", "https"))
        }
    }
--- request
GET /t
--- response_body
http://www.baidu.com starts with "http": true
http://www.baidu.com starts with "https": false
--- no_error_log
[error]



=== TEST 8: to_var_idx
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            local fmt = string.format
            ngx.say(fmt("http_x_real_ip is %s and %s", utils.to_var_idx("http_x_real_ip")))
            ngx.say(fmt("X_REAL_IP is %s and %s", utils.to_var_idx("X_REAL_IP")))
            ngx.say(fmt("X-Forwarded-For: 1 is %s and %d", utils.to_var_idx("X-Forwarded-For: 1")))
            ngx.say(fmt("X-Forwarded-For: -1 is %s and %d", utils.to_var_idx("X-Forwarded-For: -1")))
            ngx.say(fmt("X-FORWARDED-FOR:100 is %s and %d", utils.to_var_idx("X-FORWARDED-FOR:-100")))
        }
    }
--- request
GET /t
--- response_body
http_x_real_ip is http_x_real_ip and nil
X_REAL_IP is http_x_real_ip and nil
X-Forwarded-For: 1 is http_x_forwarded_for and 1
X-Forwarded-For: -1 is http_x_forwarded_for and -1
X-FORWARDED-FOR:100 is http_x_forwarded_for and -100
--- no_error_log
[error]



=== TEST 9: get_indexed_element
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say("X-Forwarded-For: 1 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 1))
            ngx.say("X-Forwarded-For: 2 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 2))
            ngx.say("X-Forwarded-For: 3 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 3))
            ngx.say("X-Forwarded-For: 4 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 4))
            ngx.say("X-Forwarded-For: 5 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 5))
            ngx.say("X-Forwarded-For: -1 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -1))
            ngx.say("X-Forwarded-For: -2 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -2))
            ngx.say("X-Forwarded-For: -3 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -3))
            ngx.say("X-Forwarded-For: -4 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -4))
            ngx.say("X-Forwarded-For: -5 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, -5))
            ngx.say("X-Forwarded-For: 0 is ", utils.get_indexed_element(ngx.var.http_x_forwarded_for, 0))
            ngx.say("X-Non-Existent-Header is ", utils.get_indexed_element(ngx.var.http_non_existent_header))
            ngx.say("X-Real-IP is ", utils.get_indexed_element(ngx.var.http_x_real_ip))
            ngx.say("X-Real-IP: 1 is ", utils.get_indexed_element(ngx.var.http_x_real_ip, 1))
            ngx.say("X-Real-IP: 2 is ", utils.get_indexed_element(ngx.var.http_x_real_ip, 2))
            ngx.say("X-Real-IP: -1 is ", utils.get_indexed_element(ngx.var.http_x_real_ip, -1))
            ngx.say("X-Real-IP: -2 is ", utils.get_indexed_element(ngx.var.http_x_real_ip, -2))
        }
    }
--- request
GET /t
--- more_headers
X-Forwarded-For: 1.1.1.1, 2.2.2.2
X-Forwarded-For: 3.3.3.3, 4.4.4.4
X-Real-IP: 10.10.10.10
--- response_body
X-Forwarded-For: 1 is 1.1.1.1
X-Forwarded-For: 2 is 2.2.2.2
X-Forwarded-For: 3 is 3.3.3.3
X-Forwarded-For: 4 is 4.4.4.4
X-Forwarded-For: 5 is nil
X-Forwarded-For: -1 is 4.4.4.4
X-Forwarded-For: -2 is 3.3.3.3
X-Forwarded-For: -3 is 2.2.2.2
X-Forwarded-For: -4 is 1.1.1.1
X-Forwarded-For: -5 is nil
X-Forwarded-For: 0 is 1.1.1.1, 2.2.2.2, 3.3.3.3, 4.4.4.4
X-Non-Existent-Header is nil
X-Real-IP is 10.10.10.10
X-Real-IP: 1 is 10.10.10.10
X-Real-IP: 2 is nil
X-Real-IP: -1 is 10.10.10.10
X-Real-IP: -2 is nil
--- no_error_log
[error]



=== TEST 10: get_event_id
--- http_config eval: $::HttpConfig
--- config
    location /t {
        content_by_lua_block {
            local utils = require "resty.t1k.utils"
            ngx.say(utils.get_event_id("<!-- event_id: 0988987de04844c7a3ce6d27865c9513 -->"))
            ngx.say(utils.get_event_id("<!-- event_id: 8bae6adf33864c7f8bf715a9b7a65b2c TYPE: A -->"))
            ngx.say(utils.get_event_id("<!-- TYPE B -->"))
        }
    }
--- request
GET /t
--- response_body
0988987de04844c7a3ce6d27865c9513
8bae6adf33864c7f8bf715a9b7a65b2c
nil
--- no_error_log
[error]
