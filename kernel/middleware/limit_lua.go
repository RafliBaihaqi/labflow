package middleware

const slidingWindowLuaScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

local clearBefore = now - window

-- Remove timestamps older than the window
redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)

-- Count remaining timestamps
local count = redis.call('ZCARD', key)

if count < limit then
    redis.call('ZADD', key, now, member)
    redis.call('EXPIRE', key, math.ceil(window / 1000) + 1)
    return {1, limit - count - 1} -- [allowed, remaining]
else
    return {0, 0} -- [allowed, remaining]
end
`
