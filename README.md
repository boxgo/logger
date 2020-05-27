# Log Format

[time] [level] [uid1|uid2|...] [traceId] [spanId] [bizId] [logId] [message]

## HTTP Log

2020-04-22T10:00:01.494+0800 INFO [openId|portalId|clientId] [PPFmbx3ZR] [1] [GET /user/info] {"method": "", "path": "", "ip": "", "user-agent": "", "query": "", "body": ""}
2020-04-22T10:00:99.000+0800 INFO [openId|portalId|clientId] [PPFmbx3ZR] [1] [GET /user/info] {"method": "", "path": "", "ip": "", "user-agent": "", "query": "", "body": "", "resp":""}

## Request Log

2020-04-22T10:00:01.494+0800 INFO [openId|portalId|clientId] [PPFmbx3ZR] [2] [GET /user/info proxy_start] {"method": "", "path": "", "ip": "", "user-agent": "", "query": "", "body": ""}
2020-04-22T10:00:01.494+0800 INFO [openId|portalId|clientId] [PPFmbx3ZR] [2] [GET /user/info proxy_end] {"method": "", "path": "", "ip": "", "user-agent": "", "query": "", "body": ""}

## Normal log

2020-04-22T10:00:01.494+0800 INFO [openId|portalId|clientId] [PPFmbx3ZR] [2] [bizId] {"method": "", "path": "", "ip": "", "user-agent": "", "query": "", "body": ""}