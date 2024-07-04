local t = {}

t.ACTION_PASSED = "."
t.ACTION_BLOCKED = "?"

t.MODE_OFF = "off"
t.MODE_BLOCK = "block"
t.MODE_MONITOR = "monitor"

t.T1K_HEADER_SIZE = 5

t.TAG_HEAD = 0x01
t.TAG_BODY = 0x02
t.TAG_EXTRA = 0x03
t.TAG_VERSION = 0x20
t.TAG_EXTRA_HEADER = 0x23
t.TAG_EXTRA_BODY = 0x24

t.MASK_FIRST = 0x40
t.MASK_LAST = 0x80

t.NGX_HTTP_HEADER_PREFIX = "http_"

t.BLOCK_CONTENT_TYPE = "application/json"
t.BLOCK_CONTENT_FORMAT = [[
{"code": %s, "success":false, "message": "blocked by Chaitin SafeLine Web Application Firewall", "event_id": "%s"}]]

t.UNIX_SOCK_PREFIX = "unix:"

return t
