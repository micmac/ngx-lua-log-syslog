
syslog = require 'mod.syslog'

if not ngx then
    ngx = {}
    ngx.var = {}
    ngx.status = 200
    ngx.var.bytes_sent = 123
    ngx.var.http_referer = "http://referer.com/"
    ngx.var.http_user_agent = "User Agent v1"
    ngx.localtime = function() return syslog.iso_timestamp(os.time()) end
    ngx.req = {}
    ngx.req.get_method = function() return "GET" end
    ngx.var.uri = "/foo/bar"
    ngx.var.remote_addr = "1.3.2.4"
end
      
-- log_format combined
-- '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"'

-- you will have to set log format here for now
local fmtstring = '%s - %s [%s] "%s" %s %d "%s" "%s"'
local logline = string.format(fmtstring,
                              ngx.var.remote_addr,
                              "-",
                              ngx.localtime(),
                              ngx.req.get_method() .." ".. ngx.var.uri,
                              ngx.status,
                              ngx.var.bytes_sent,
                              ngx.var.http_referer or "-",
                              ngx.var.http_user_agent
                )

local logfacility = ngx.var.access_loglevel or LOG_LOCAL7
local loglevel = ngx.var.access_logfacility or LOG_INFO
local logtag = ngx.var.access_logtag or "nginx-lua"


syslog:log(logfacility, loglevel, logtag, 0, logline)

