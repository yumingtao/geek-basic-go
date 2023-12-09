-- 发送到的 key, code:biz:phone
local key = KEYS[1]
-- 验证次数
local cntKey = key .. ":cnt"
-- 准备存储的验证码
local val = ARGV[1]
-- 验证码有效时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- key 存在但是没有过期时间
    return -2 -- Go里边得到的返回值是-2
elseif ttl == -2 or ttl < 540 then
    -- 可以发送验证码
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end
