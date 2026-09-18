local t1k = require "resty.t1k"
local t1k_constants = require "resty.t1k.constants"

local ngx = ngx
local fmt = string.format
local lower = string.lower

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

local function is_block_mode(t)
    -- resty.t1k lower cases the mode before comparing it, so do the same here:
    -- a mode of "BLOCK" has to fail closed exactly like "block".
    return type(t.mode) == "string" and lower(t.mode) == t1k_constants.MODE_BLOCK
end

-- fail_closed answers a request that could not be checked.
--
-- The detector was not asked (or the plugin is not configured at all), so
-- nothing is known about this request. Answering it while the plugin runs in
-- block mode would turn a detector outage into a complete bypass of the WAF,
-- which is the one outcome an inline protection must not have. monitor mode
-- keeps passing traffic, and off mode never reaches this point.
local function fail_closed(t)
    if not is_block_mode(t) then
        return
    end

    local status = ngx.HTTP_INTERNAL_SERVER_ERROR
    local msg = fmt(blocked_message, status, "")
    ngx.status = status
    ngx.header.content_type = t1k_constants.BLOCK_CONTENT_TYPE
    ngx.say(msg)
    return ngx.exit(status)
end

function _M.rewrite()
    local t = get_conf()
    if not t.host then
        ngx.log(ngx.ERR, "safeline host is required")
        -- resty.t1k refuses to run without a host, so this is the same
        -- "cannot be checked" situation as a detector that does not answer.
        return fail_closed(t)
    end
    local ok, err, result = t1k.do_access(t, false)
    if not ok then
        ngx.log(ngx.ERR, "failed to detector req: ", err)
        return fail_closed(t)
    end
    if result then
        if result.action == t1k_constants.ACTION_BLOCKED then
            local status = tonumber(result.status, 10) or ngx.HTTP_FORBIDDEN
            local msg = fmt(blocked_message, status, result.event_id)
            ngx.log(ngx.ERR, "blocked by safeline waf: ", msg)
            ngx.status = status
            ngx.header.content_type = t1k_constants.BLOCK_CONTENT_TYPE
            ngx.say(msg)
            return ngx.exit(status)
        end
    end
end

return _M
