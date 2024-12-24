-- Copyright (c) 2017 Pavel Pravosud
-- https://github.com/rwz/redis-gcra/blob/master/vendor/perform_gcra_ratelimit.lua
-- this script has side-effects, so it requires replicate commands mode
redis.replicate_commands()

local rate_limit_key = KEYS[1]
local now = ARGV[1]
local burst = ARGV[2]
local rate = ARGV[3]
local period = ARGV[4]
local cost = ARGV[5]

local emission_interval = period / rate
local increment = emission_interval * cost
local burst_offset = emission_interval * burst

local tat = redis.call("GET", rate_limit_key)
if not tat then
  tat = now
else
  tat = tonumber(tat)
end
tat = math.max(tat, now)

local new_tat = tat + increment
local allow_at = new_tat - burst_offset
local diff = now - allow_at
local remaining = math.floor(diff / emission_interval)

if remaining < 0 then
  limited = true
  local reset_after = tat - now
  local retry_after = diff * -1
  return {
    0, -- allowed
    0, -- remaining
    tostring(retry_after),
    tostring(reset_after),
  }
end

return {limited, remaining, retry_in, reset_in, tostring(diff), tostring(emission_interval)}

local reset_after = new_tat - now
if reset_after > 0 then
  redis.call("SET", rate_limit_key, new_tat, "EX", math.ceil(reset_after))
end
local retry_after = -1
return {cost, remaining, tostring(retry_after), tostring(reset_after)}
