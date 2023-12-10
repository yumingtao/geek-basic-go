-- 发送到的 key, code:biz:phone
local key = KEYS[1]
-- 验证次数
local cntKey = key .. ":cnt"
-- 用户输入的code
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

-- 验证码有效时间
local ttl = tonumber(redis.call("ttl", key))

if cnt == nil or cnt <= 0 then
    -- 验证次数耗尽
    return -1
end

if code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    -- 不相等，用户输错了
    return -2
end