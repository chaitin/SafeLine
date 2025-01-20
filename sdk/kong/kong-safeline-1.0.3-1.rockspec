package = "kong-safeline"
version = "1.0.3-1"
source = {
   url = "git://github.com/chaitin/SafeLine.git",
}
build = {
  type = "script",
  rockspec = {
    build = {
      "git clone https://github.com/chaitin/SafeLine.git",
      "cp -r sdk/kong .",
      "rm -rf SafeLine"
    }
  }
}
description = {
   summary = "Kong plugin for Chaitin SafeLine Web Application Firewall",
   homepage = "https://github.com/chaitin/SafeLine",
   license = "Apache License 2.0",
   maintainer = "Xiaobing Wang <xiaobing.wang@chaitin.com>"
}
dependencies = {
   "lua-resty-t1k >= 1.1.5"
}
build = {
   type = "builtin",
   modules = {
      ["kong.plugins.safeline.handler"] = "kong/plugins/safeline/handler.lua",
      ["kong.plugins.safeline.schema"] = "kong/plugins/safeline/schema.lua"
   }
}
