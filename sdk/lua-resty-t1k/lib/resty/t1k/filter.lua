local _M = {
    _VERSION = '1.0.0'
}

local find = string.find
local sub = string.sub

local ngx = ngx

local function parse_extra_header(extra_header)
    local t = {}
    local idx = 1

    while (idx <= #extra_header) do
        local key, val
        local _, to = find(extra_header, ":", idx)
        if to == 0 then
            break
        else
            key = sub(extra_header, idx, to - 1)
        end

        idx = to + 1
        _, to = find(extra_header, "\n", idx)
        if to == 0 then
            break
        else
            val = sub(extra_header, idx, to - 1)
        end

        t[key] = val
        idx = to + 1
    end

    return t
end

function _M.do_header_filter()
    local extra_header = ngx.ctx.t1k_extra_header
    if extra_header ~= nil then
        local header_table = parse_extra_header(extra_header)
        for k, v in pairs(header_table) do
            if k ~= nil and v ~= nil then
                ngx.header[k] = v
            end
        end
    end
end

return _M
