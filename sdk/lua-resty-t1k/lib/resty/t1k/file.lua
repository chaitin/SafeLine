local buffer = require "resty.t1k.buffer"

local _M = {
    _VERSION = '1.0.0'
}

local buffer_size = 2 ^ 13

function _M.read(p, size)
    size = (not size or size < 0) and 0 or size

    local f, err = io.open(p, "rb")
    if not f or err then
        return nil, err, nil
    end

    local left = size
    local buf = buffer:new()

    while left ~= 0 do
        local block_size = math.min(left, buffer_size)
        local block = f:read(block_size)
        if not block then
            break
        end
        buf:add(block)
        left = math.max(left - block_size, 0)
    end
    f:close()

    return true, nil, buf
end

return _M
