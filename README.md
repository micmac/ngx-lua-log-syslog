ngx-lua-log-syslog
==================

Playground for generating access log from nginx to a syslog server and forwarding syslog to scribe

Lua script for Nginx's log_by_lua_file and a syslog module in Lua to send access logs to a remote UDP port.

Log forwarder in Go to receive syslog protocol and send it to scribe.

Todo:

  - routing (this port/facility/level to that scribe category)
  - actual writing to scribe
  - tests :(
