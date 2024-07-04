local consts = require "resty.t1k.constants"
local log = require "resty.t1k.log"

local fmt = string.format

local ngx = ngx

local log_fmt = log.fmt

local _M = {
    _VERSION = '1.0.0'
}

function _M.handle(t)
    local t_type = type(t)
    if t_type ~= "table" then
        local err = log_fmt("invalid result type: %s", t_type)
        return nil, err
    end

    local action = t["action"]
    if action == consts.ACTION_PASSED then
        return true, nil
    elseif action == consts.ACTION_BLOCKED then
        ngx.status = t["status"] or ngx.HTTP_FORBIDDEN
        ngx.header.content_type = consts.BLOCK_CONTENT_TYPE
        ngx.say(fmt(consts.BLOCK_CONTENT_FORMAT, ngx.status, t["event_id"]))

        return ngx.exit(ngx.status)
    else
        local err = log_fmt("unknown action from t1k server: %s", action)
        return nil, err
    end
end

return _M
