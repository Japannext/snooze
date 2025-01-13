
redis.call("zrangebyscore", key, "-inf", "+inf", "LIMIT", "0", "1")
