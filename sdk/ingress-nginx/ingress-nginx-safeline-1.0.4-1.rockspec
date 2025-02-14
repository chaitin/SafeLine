package = "ingress-nginx-safeline"
version = "1.0.4-1"
source = {
   url = "git://github.com/chaitin/ingress-nginx-safeline.git"
}
description = {
   summary = "Ingress-Nginx plugin for Chaitin SafeLine Web Application Firewall",
   homepage = "https://github.com/chaitin/ingress-nginx-safeline",
   license = "Apache License 2.0",
   maintainer = "Xiaobing Wang <xiaobing.wang@chaitin.com>"
}
dependencies = {
   "lua-resty-t1k >= 1.1.5"
}
build = {
   type = "builtin",
   modules = {
       ["safeline.main"] = "lib/safeline/main.lua"
   }
}
