local t1k = require "resty.t1k"
local t1k_constants = require "resty.t1k.constants"

local ngx = ngx
local fmt = string.format

local blocked_message = [[{"code": %s, "success":false, ]] ..
        [["message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "%s"}]]

local _M = {}

local mode = os.getenv("SAFELINE_MODE")
local host = os.getenv("SAFELINE_HOST")
local port = os.getenv("SAFELINE_PORT")
local connect_timeout = os.getenv("SAFELINE_CONNECT_TIMEOUT")
local send_timeout = os.getenv("SAFELINE_SEND_TIMEOUT")
local read_timeout = os.getenv("SAFELINE_READ_TIMEOUT")
local req_body_size = os.getenv("SAFELINE_REQ_BODY_SIZE")
local keepalive_size = os.getenv("SAFELINE_KEEPALIVE_SIZE")
local keepalive_timeout = os.getenv("SAFELINE_KEEPALIVE_TIMEOUT")
local remote_addr = os.getenv("SAFELINE_REMOTE_ADDR")

local function get_conf()
    local t = {
        mode = mode or "block",
        host = host,
        port = port,
        connect_timeout = connect_timeout or 1000,
        send_timeout = send_timeout or 1000,
        read_timeout = read_timeout or 1000,
        req_body_size = req_body_size or 1024,
        keepalive_size = keepalive_size or 256,
        keepalive_timeout = keepalive_timeout or 60000,
        remote_addr = remote_addr
    }
    return t
end

function _M.rewrite()
    local t = get_conf()
    if not t.host then
        ngx.log(ngx.ERR, "safeline host is required")
        return
    end
    local ok, err, result = t1k.do_access(t, false)
    if not ok then
        ngx.log(ngx.ERR, "failed to detector req: ", err)
        return
    end
    if result then
        if result.action == t1k_constants.ACTION_BLOCKED then
            local msg = fmt(blocked_message, result.status, result.event_id)
            ngx.log(ngx.ERR, "blocked by safeline waf: ",msg)
            ngx.status = tonumber(result.status,10)
            ngx.say(msg)
            return ngx.exit(ngx.HTTP_OK)
        end
    end
end

return _M