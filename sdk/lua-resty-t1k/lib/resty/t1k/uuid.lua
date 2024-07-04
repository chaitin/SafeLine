local bit = require "bit"
local ffi = require "ffi"

local log = require "resty.t1k.log"

local _M = {
    _VERSION = '1.0.0'
}

local C = ffi.C
local N_BYTES = 32
local random = math.random

local nlog = ngx.log

local warn_fmt = log.warn_fmt

ffi.cdef [[
    int RAND_bytes(unsigned char *buf, int num);
]]

local function _rand_bytes(buf, len)
    return C.RAND_bytes(buf, len)
end

local function rand_bytes(len)
    local buf = ffi.new("char[?]", len)
    local ok, ret = pcall(_rand_bytes, buf, len)

    if not ok or ret ~= 1 then
        nlog(warn_fmt("call RAND_bytes failed: %s", ret))
        return nil
    end

    return ffi.string(buf, len)
end

do
    local band = bit.band
    local bor = bit.bor
    local tohex = bit.tohex

    local fmt = string.format
    local byte = string.byte

    function _M.generate_v4()
        local bytes = rand_bytes(N_BYTES)

        -- fallback to math.random based method
        if not bytes then
            return (fmt('%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s',
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),

                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),

                    tohex(bor(band(random(0, 255), 0x0F), 0x40), 2),
                    tohex(random(0, 255), 2),

                    tohex(bor(band(random(0, 255), 0x3F), 0x80), 2),
                    tohex(random(0, 255), 2),

                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2),
                    tohex(random(0, 255), 2)))
        end

        return fmt('%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s',
                tohex(byte(bytes, 1, 2), 2),
                tohex(byte(bytes, 3, 4), 2),
                tohex(byte(bytes, 5, 6), 2),
                tohex(byte(bytes, 7, 8), 2),

                tohex(byte(bytes, 9, 10), 2),
                tohex(byte(bytes, 11, 12), 2),

                tohex(bor(band(byte(bytes, 13, 14), 0x0F), 0x40), 2),
                tohex(byte(bytes, 15, 16), 2),

                tohex(bor(band(byte(bytes, 17, 18), 0x3F), 0x80), 2),
                tohex(byte(bytes, 19, 20), 2),

                tohex(byte(bytes, 21, 22), 2),
                tohex(byte(bytes, 23, 24), 2),
                tohex(byte(bytes, 25, 26), 2),
                tohex(byte(bytes, 27, 28), 2),
                tohex(byte(bytes, 29, 30), 2),
                tohex(byte(bytes, 31, 32), 2))
    end
end

return setmetatable(_M, {
    __call = _M.generate_v4
})
