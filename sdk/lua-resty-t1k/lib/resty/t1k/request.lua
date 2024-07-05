local bit = require "bit"

local buffer = require "resty.t1k.buffer"
local consts = require "resty.t1k.constants"
local file = require "resty.t1k.file"
local log = require "resty.t1k.log"
local utils = require "resty.t1k.utils"
local uuid = require "resty.t1k.uuid"

local _M = {
    _VERSION = '1.0.0',
}

local bor = bit.bor
local byte = string.byte
local char = string.char
local fmt = string.format
local sub = string.sub
local concat = table.concat

local ngx = ngx
local nlog = ngx.log
local ngx_req = ngx.req
local ngx_socket = ngx.socket
local ngx_var = ngx.var

local warn_fmt = log.warn_fmt
local debug_fmt = log.debug_fmt

local KEY_EXTRA_UUID = "UUID"
local KEY_EXTRA_LOCAL_ADDR = "LocalAddr"
local KEY_EXTRA_LOCAL_PORT = "LocalPort"
local KEY_EXTRA_REMOTE_ADDR = "RemoteAddr"
local KEY_EXTRA_REMOTE_PORT = "RemotePort"
local KEY_EXTRA_SCHEME = "Scheme"
local KEY_EXTRA_SERVER_NAME = "ServerName"
local KEY_EXTRA_PROXY_NAME = "ProxyName"
local KEY_EXTRA_REQ_BEGIN_TIME = "ReqBeginTime"
local KEY_EXTRA_HAS_RSP_IF_OK = "HasRspIfOK"
local KEY_EXTRA_HAS_RSP_IF_BLOCK = "HasRspIfBlock"

local TAG_HEAD_WITH_MASK_FIRST = bor(consts.TAG_HEAD, consts.MASK_FIRST)
local TAG_EXTRA_WITH_MASK_LAST = bor(consts.TAG_EXTRA, consts.MASK_LAST)

local T1K_PROTO = "Proto:2\n"
local T1K_PROTO_DATA = fmt("%s%s%s", char(consts.TAG_VERSION), utils.int_to_char_length(#T1K_PROTO), T1K_PROTO)

local function read_request_body(opt_req_body_size)
    local ok, err
    local req_body, req_body_size

    local content_length = tonumber(ngx_var.http_content_length) or 0
    local transfer_encoding = ngx_var.http_transfer_encoding
    if content_length == 0 and not transfer_encoding then
        return true, nil, nil
    end

    ngx_req.read_body()
    req_body = ngx_req.get_body_data()
    if req_body then
        req_body_size = #req_body
        if req_body_size > opt_req_body_size then
            nlog(debug_fmt("request body is too long: %d bytes, cut to %d bytes", req_body_size, opt_req_body_size))
            req_body = sub(req_body, 1, opt_req_body_size)
        end

        return true, nil, req_body
    end

    local path = ngx_req.get_body_file()
    if not path then
        return true, nil, nil
    end

    ok, err, req_body = file.read(path, opt_req_body_size)
    if not ok then
        err = fmt("failed to read temporary file %s: %s", path, err)
        return ok, err, nil
    end

    return true, nil, req_body
end

local function get_remote_addr(remote_addr_var, remote_addr_idx)
    local addr
    if remote_addr_var then
        addr = utils.get_indexed_element(ngx_var[remote_addr_var], remote_addr_idx)
    end
    return addr or ngx_var.remote_addr
end

local function parse_v(v)
    if type(v) == "table" then
        return concat(v, ", ")
    end
    return tostring(v)
end

local function build_header()
    local http_version = ngx_req.http_version()
    if http_version < 2.0 then
        return true, nil, ngx_req.raw_header()
    end

    local headers, err = ngx_req.get_headers(0, true)
    if err then
        err = fmt("failed to call ngx_req.get_headers: %s", err)
        return nil, err, nil
    end

    local buf = buffer:new()
    buf:add(fmt("%s %s HTTP/%.1f\r\n", ngx_req.get_method(), ngx_var.request_uri, http_version))

    for k, v in pairs(headers) do
        buf:add(fmt("%s: %s\r\n", k, parse_v(v)))
    end
    buf:add("\r\n")

    return true, nil, buf
end

local function build_body(opts)
    local ok, err
    local body

    local req_body_size = opts.req_body_size * 1024
    ok, err, body = read_request_body(req_body_size)
    if not ok then
        return ok, err, nil
    end

    return true, nil, body
end

local function build_extra(opts)
    local err

    local src_ip = get_remote_addr(opts.remote_addr_var, opts.remote_addr_idx)
    if not src_ip then
        err = fmt("failed to get remote_addr, var: %s, idx %d", opts.remote_addr_var, opts.remote_addr_idx)
        return nil, err
    end

    local src_port = ngx_var.remote_port
    if not src_port then
        err = "failed to get ngx_var.remote_port"
        return nil, err, nil
    end

    local local_ip = ngx_var.server_addr
    if not local_ip then
        err = "failed to get ngx_var.server_addr"
        return nil, err, nil
    end

    local local_port = ngx_var.server_port
    if not local_port then
        err = "failed to get ngx_var.server_port"
        return nil, err, nil
    end

    local extra = buffer:new({
        KEY_EXTRA_UUID, ":", uuid.generate_v4(), "\n",
        KEY_EXTRA_REMOTE_ADDR, ":", src_ip, "\n",
        KEY_EXTRA_REMOTE_PORT, ":", src_port, "\n",
        KEY_EXTRA_LOCAL_ADDR, ":", local_ip, "\n",
        KEY_EXTRA_LOCAL_PORT, ":", local_port, "\n",
        KEY_EXTRA_SCHEME, ":", ngx_var.scheme, "\n",
        KEY_EXTRA_SERVER_NAME, ":", ngx_var.server_name, "\n",
        KEY_EXTRA_PROXY_NAME, ":", ngx_var.hostname, "\n",
        KEY_EXTRA_REQ_BEGIN_TIME, ":", fmt("%.0f", ngx_req.start_time() * 1000000), "\n",
        KEY_EXTRA_HAS_RSP_IF_OK, ":n\n",
        KEY_EXTRA_HAS_RSP_IF_BLOCK, ":n\n"
    })

    return true, nil, extra
end

local function do_send(sock, data)
    local ok, err = sock:send(data)
    if not ok then
        return ok, err
    end
    return true, nil
end

local function receive_data(s, srv)
    local t = {}
    local ft = true
    local finished

    repeat
        local err
        local tag, length, packet, rsp_body

        packet, err = s:receive(consts.T1K_HEADER_SIZE)
        if err then
            err = fmt("failed to receive info packet from t1k server %s: %s", srv, err)
            return nil, err, nil
        end
        if not packet then
            err = fmt("empty packet from t1k server %s", srv)
            return nil, err, nil
        end

        if ft then
            if not utils.is_mask_first(byte(packet, 1, 1)) then
                err = fmt("first packet is not MASK_FIRST from t1k server %s", srv)
                return nil, err, nil
            end
            ft = false
        end

        finished, tag, length = utils.packet_parser(packet)
        if length > 0 then
            rsp_body, err = s:receive(length)
            if not rsp_body or #rsp_body ~= length then
                err = fmt("failed to receive payload from t1k server %s: %s", srv, err)
                return nil, err, nil
            end
            t[tag] = rsp_body
        end

    until (finished)

    return true, nil, t
end

local function get_socket(opts)
    local ok, err
    local count, sock, server

    sock, err = ngx_socket.tcp()
    if not sock then
        err = fmt("failed to create socket: %s", err)
        return nil, err, nil, nil
    end

    sock:settimeouts(opts.connect_timeout, opts.send_timeout, opts.read_timeout)

    if opts.uds then
        server = opts.host
        ok, err = sock:connect(opts.host)
    else
        server = fmt("%s:%d", opts.host, opts.port)
        ok, err = sock:connect(opts.host, opts.port)
    end
    if not ok then
        sock:close()
        err = fmt("failed to connect to t1k server %s: %s", server, err)
        return ok, err, nil, nil
    end
    nlog(debug_fmt("successfully connected to t1k server %s", server))

    count, err = sock:getreusedtimes()
    if not count then
        nlog(warn_fmt("failed to get reused times from t1k server %s: %s", server, err))
    end

    if count == 0 then
        ok, err = sock:setoption("keepalive", true)
        if not ok then
            nlog(warn_fmt("failed to set keepalive for t1k server %s: %s", server, err))
        end
        ok, err = sock:setoption("reuseaddr", true)
        if not ok then
            nlog(warn_fmt("failed to set reuseaddr for t1k server %s: %s", server, err))
        end
        ok, err = sock:setoption("tcp-nodelay", true)
        if not ok then
            nlog(warn_fmt("failed to set tcp-nodelay for t1k server %s: %s", server, err))
        end
    end

    return true, nil, sock, server
end

local function do_socket(opts, header, body, extra)
    local ok, err
    local t, sock, server

    ok, err, sock, server = get_socket(opts)
    if not ok then
        err = fmt("failed to get socket: %s", err)
        return ok, err, nil
    end

    ok, err = do_send(sock, { char(TAG_HEAD_WITH_MASK_FIRST), utils.int_to_char_length(header:len()), header })
    if not ok then
        sock:close()
        err = fmt("failed to send header data to t1k server %s: %s", server, err)
        return ok, err, nil
    end

    if body ~= nil then
        ok, err = do_send(sock, { char(consts.TAG_BODY), utils.int_to_char_length(body:len()), body })
        if not ok then
            sock:close()
            err = fmt("failed to send body data to t1k server %s: %s", server, err)
            return ok, err, nil
        end
    end

    ok, err = do_send(sock, { T1K_PROTO_DATA, char(TAG_EXTRA_WITH_MASK_LAST), utils.int_to_char_length(extra:len()), extra })
    if not ok then
        sock:close()
        err = fmt("failed to send extra data to t1k server %s: %s", server, err)
        return ok, err, nil
    end

    ok, err, t = receive_data(sock, server)
    if not ok then
        return ok, err, nil
    end

    ok, err = sock:setkeepalive(opts.keepalive_timeout, opts.keepalive_size)
    if not ok then
        nlog(warn_fmt("failed to set keepalive: %s", err))
        sock:close()
    end

    return true, nil, t
end

function _M.do_request(opts)
    local ok, err
    local header, body, extra, t

    ok, err, header = build_header(opts)
    if not ok then
        return ok, err, nil
    end

    ok, err, body = build_body(opts)
    if not ok then
        return ok, err, nil
    end

    ok, err, extra = build_extra(opts)
    if not ok then
        return ok, err, nil
    end

    ok, err, t = do_socket(opts, header, body, extra)
    if not ok then
        return ok, err, nil
    end

    if opts.mode == consts.MODE_BLOCK then
        local extra_header = t[consts.TAG_EXTRA_HEADER]
        if extra_header then
            ngx.ctx.t1k_extra_header = extra_header
        end
    end

    local result = {
        action = t[consts.TAG_HEAD],
        status = t[consts.TAG_BODY],
        event_id = utils.get_event_id(t[consts.TAG_EXTRA_BODY]),
    }

    return true, nil, result
end

return _M
