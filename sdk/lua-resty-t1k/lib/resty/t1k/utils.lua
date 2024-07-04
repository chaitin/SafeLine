local consts = require "resty.t1k.constants"

local _M = {
    _VERSION = '1.0.0'
}

local band = bit.band
local bnot = bit.bnot
local lshift = bit.lshift
local rshift = bit.rshift

local abs = math.abs
local char = string.char

local ngx = ngx
local ngx_re = ngx.re
local re_match = ngx_re.match
local re_gmatch = ngx_re.gmatch
local re_gsub = ngx_re.gsub

local NOT_MASK_FIRST = bnot(consts.MASK_FIRST)
local NOT_MASK_LAST = bnot(consts.MASK_LAST)

function _M.int_to_char_length(x)
    return char(band(x, 0xff)) .. char(band(rshift(x, 8), 0xff)) ..
            char(band(rshift(x, 16), 0xff)) .. char(band(rshift(x, 24), 0xff))
end

function _M.char_to_int_length(l)
    return l:byte(1, 1) + lshift(l:byte(2, 2), 8) + lshift(l:byte(3, 3), 16) + lshift(l:byte(4, 4), 24)
end

function _M.is_mask_first(b)
    return band(b, consts.MASK_FIRST) == consts.MASK_FIRST
end

function _M.is_mask_last(b)
    return band(b, consts.MASK_LAST) == consts.MASK_LAST
end

local function tag_parser(tag)
    return band(band(tag, NOT_MASK_FIRST), NOT_MASK_LAST)
end

function _M.packet_parser(header)
    if #header ~= consts.T1K_HEADER_SIZE then
        return true, nil, 0
    end

    local fb = header:byte(1, 1)
    local finished = _M.is_mask_last(fb)
    local tag = tag_parser(fb)
    local length = _M.char_to_int_length(header:sub(2, 5))

    return finished, tag, length
end

function _M.starts_with(str, start)
    return str:sub(1, #start) == start
end

function _M.to_var_idx(o)
    local var = o
    local idx

    local _, p = o:find(":")
    if p then
        var = o:sub(1, p - 1)
        idx = tonumber(o:sub(p + 1))
    end

    var = re_gsub(var:lower(), "-", "_")
    if not _M.starts_with(var, consts.NGX_HTTP_HEADER_PREFIX) then
        var = consts.NGX_HTTP_HEADER_PREFIX .. var
    end

    return var, idx
end

function _M.get_indexed_element(str, idx)
    if not str or not idx or idx == 0 then
        return str
    end

    local it, err = re_gmatch(str, [[([^,\s]+)]], "jo")
    if err then
        return nil
    end

    local t = {}
    for m, e in it do
        if e then
            return nil
        end
        table.insert(t, m[1])
    end

    local len = #t
    if len < abs(idx) then
        return nil
    end

    return t[idx > 0 and idx or (len + idx + 1)]
end

function _M.get_event_id(str)
    if not str then
        return nil
    end

    local m, err = re_match(str, [[<!-- event_id: ([A-Za-z0-9]+)(?: TYPE: [a-zA-Z])? -->]], "jo")
    if err then
        return nil
    end

    if m then
        return m[1]
    end

    return nil
end

return _M
