      
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

syslog = require 'syslog'

syslog:log(logfacility, loglevel, logtag, 0, logline)

