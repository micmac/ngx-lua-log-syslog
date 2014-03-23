#!/usr/bin/lua5.1

local S = {}

-- contants from <sys/syslog.h>
LOG_EMERG    =  0       -- system is unusable */
LOG_ALERT    =  1       -- action must be taken immediately */
LOG_CRIT     =  2       -- critical conditions */
LOG_ERR      =  3       -- error conditions */
LOG_WARNING  =  4       -- warning conditions */
LOG_NOTICE   =  5       -- normal but significant condition */
LOG_INFO     =  6       -- informational */
LOG_DEBUG    =  7       -- debug-level messages */

LOG_KERN     =  (0 *8)  -- kernel messages */
LOG_USER     =  (1 *8)  -- random user-level messages */
LOG_MAIL     =  (2 *8)  -- mail system */
LOG_DAEMON   =  (3 *8)  -- system daemons */
LOG_AUTH     =  (4 *8)  -- security/authorization messages */
LOG_SYSLOG   =  (5 *8)  -- messages generated internally by syslogd */
LOG_LPR      =  (6 *8)  -- line printer subsystem */
LOG_NEWS     =  (7 *8)  -- network news subsystem */
LOG_UUCP     =  (8 *8)  -- UUCP subsystem */
LOG_CRON     =  (9 *8)  -- clock daemon */
LOG_AUTHPRIV =  (10 *8) -- security/authorization messages (private) */
LOG_FTP      =  (11 *8) -- ftp daemon */

-- other codes through 15 reserved for system use */
LOG_LOCAL0   =  (16 *8) -- reserved for local use */
LOG_LOCAL1   =  (17 *8) -- reserved for local use */
LOG_LOCAL2   =  (18 *8) -- reserved for local use */
LOG_LOCAL3   =  (19 *8) -- reserved for local use */
LOG_LOCAL4   =  (20 *8) -- reserved for local use */
LOG_LOCAL5   =  (21 *8) -- reserved for local use */
LOG_LOCAL6   =  (22 *8) -- reserved for local use */
LOG_LOCAL7   =  (23 *8) -- reserved for local use */

S.myhostname = "localhost"

function S.mkprio(fac, sev)
    return fac + sev
end

-- parameter is unix time or nil if now
function S.iso_timestamp(uts)
    tz = os.date("*t", uts).hour - os.date("!*t", uts).hour
    if tz < 0 then tz = tz + 24 end
    tzs = string.format("%.4d", tz * 100)
    ts = os.date("%Y-%m-%dT%H:%M:%S+", uts)
    return ts .. tzs
end

-- timestamp is now added by syslogd (rsyslog)
function S:mklogline(fac, sev, tag, pid, msg)
    prio_field = string.format("<%d>", S.mkprio(fac, sev))
    if pid and (pid > 0) then
        pid_field = string.format("[%d]", pid)
    else
        pid_field = ""
    end
    host_field = self.myhostname
    return prio_field .. host_field .. " " .. tag .. pid_field .. ": " .. msg
end

function S:log(fac, sev, tag, pid, msg)
    socket = require 'socket'
    udpsock = socket.udp()
    udpsock:sendto(self:mklogline(fac, sev, tag, pid, msg), "127.0.0.1", 5140)
end

function S.notice(msg)
   syslog.log(LOG_LOCAL5, LOG_NOTICE, "smartrouter", nil, msg)
end

function S.error(msg)
   syslog.log(LOG_LOCAL5, LOG_ERR, "smartrouter", nil, msg)
end

function S.warning(msg)
   syslog.log(LOG_LOCAL5, LOG_WARNING, "smartrouter", nil, msg)
end

return S
