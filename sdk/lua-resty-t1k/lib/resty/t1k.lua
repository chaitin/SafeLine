local consts = require "resty.t1k.constants"
local filter = require "resty.t1k.filter"
local handler = require "resty.t1k.handler"
local log = require "resty.t1k.log"
local request = require "resty.t1k.request"
local utils = require "resty.t1k.utils"

local lower = string.lower

local ngx = ngx
local nlog = ngx.log

local log_fmt = log.fmt
local debug_fmt = log.debug_fmt

local _M = {
    _VERSION = '1.0.0'
}

local DEFAULT_T1K_CONNECT_TIMEOUT = 1000 -- 1s
local DEFAULT_T1K_SEND_TIMEOUT = 1000 -- 1s
local DEFAULT_T1K_READ_TIMEOUT = 1000 -- 1s
local DEFAULT_T1K_REQ_BODY_SIZE = 1024 -- 1024 KB
local DEFAULT_T1K_KEEPALIVE_SIZE = 256
local DEFAULT_T1K_KEEPALIVE_TIMEOUT = 60 * 1000 -- 60s

function _M.do_access(t, handle)
    local ok, err, result
    local opts = {}
    t = t or {}

    if not t.mode then
        return true, nil, nil
    end

    opts.mode = lower(t.mode)
    if opts.mode == consts.MODE_OFF then
        nlog(debug_fmt("t1k is not enabled"))
        return true, nil, nil
    end

    if opts.mode ~= consts.MODE_OFF and opts.mode ~= consts.MODE_BLOCK and opts.mode ~= consts.MODE_MONITOR then
        err = log_fmt("invalid t1k mode: %s", t.mode)
        return nil, err, nil
    end

    if not t.host then
        err = log_fmt("invalid t1k host: %s", t.host)
        return nil, err, nil
    end
    opts.host = t.host

    if utils.starts_with(opts.host, consts.UNIX_SOCK_PREFIX) then
        opts.uds = true
    else
        if not tonumber(t.port) then
            err = log_fmt("invalid t1k port: %s", t.port)
            return nil, err, nil
        end
        opts.port = tonumber(t.port)
    end

    opts.connect_timeout = t.connect_timeout or DEFAULT_T1K_CONNECT_TIMEOUT
    opts.send_timeout = t.send_timeout or DEFAULT_T1K_SEND_TIMEOUT
    opts.read_timeout = t.read_timeout or DEFAULT_T1K_READ_TIMEOUT
    opts.req_body_size = t.req_body_size or DEFAULT_T1K_REQ_BODY_SIZE
    opts.keepalive_size = t.keepalive_size or DEFAULT_T1K_KEEPALIVE_SIZE
    opts.keepalive_timeout = t.keepalive_timeout or DEFAULT_T1K_KEEPALIVE_TIMEOUT

    if t.remote_addr then
        local var, idx = utils.to_var_idx(t.remote_addr)
        opts.remote_addr_var = var
        opts.remote_addr_idx = idx
    end

    ok, err, result = request.do_request(opts)
    if not ok then
        return ok, err, result
    end

    if handle and opts.mode == consts.MODE_BLOCK then
        ok, err = _M.do_handle(result)
    end

    return ok, err, result
end

function _M.do_handle(t)
    local ok, err = handler.handle(t)
    return ok, err
end

function _M.do_header_filter()
    filter.do_header_filter()
end

return _M
