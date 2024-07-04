local typedefs = require "kong.db.schema.typedefs"

return {
    name = "kong-safeline",
    fields = {{
        consumer = typedefs.no_consumer
    }, {
        protocols = typedefs.protocols_http
    }, {
        config = {
            type = "record",
            fields = {{
                host = {
                    type = "string",
                    required = true
                }
            }, {
                port = {
                    type = "number",
                    required = false
                }
            }, {
                mode = {
                    type = "string",
                    required = false,
                    default = "block",
                    one_of = {"monitor", "block", "off"}
                }
            }, {
                connect_timeout = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                send_timeout = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                read_timeout = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                req_body_size = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                keepalive_size = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                keepalive_timeout = {
                    type = "number",
                    required = false,
                    default = 1000
                }
            }, {
                remote_addr = {
                    type = "string",
                    required = false
                }
            }}
        }
    }}
}
