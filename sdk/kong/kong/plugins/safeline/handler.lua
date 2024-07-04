local t1k = require "resty.t1k"
local t1k_constants = require "resty.t1k.constants"

local fmt = string.format

local SafelineHandler = {
    VERSION = "0.0.1",
    PRIORITY = 1000
}

local blocked_message = [[{"code": %s, "success":false, ]] ..
                            [["message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "%s"}]]

local function get_conf(conf)
    local t = {
        mode = conf.mode,
        host = conf.host,
        port = conf.port,
        connect_timeout = conf.connect_timeout,
        send_timeout = conf.send_timeout,
        read_timeout = conf.read_timeout,
        req_body_size = conf.req_body_size,
        keepalive_size = conf.keepalive_size,
        keepalive_timeout = conf.keepalive_timeout,
        remote_addr = conf.remote_addr
    }
    return t
end

function SafelineHandler:access(conf)
    -- your custom code here
    local t = get_conf(conf)
    local ok, err, result = t1k.do_access(t, false)
    if not ok then
        kong.log.err("failed to detector req: ", err)
    end
    if result then
        if result.action == t1k_constants.ACTION_BLOCKED then
            local msg = fmt(blocked_message, result.status, result.event_id)
            kong.log.debug("blocked by safeline: ",msg)
            return kong.response.exit(result.status, msg)
        end
    end
end


return SafelineHandler
