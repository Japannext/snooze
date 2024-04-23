package ratelimit

import (
  "time"

  "github.com/japannext/snooze/common/redis"
)

type Rule struct {
  Burst int64 `yaml:"burst"`
  Period time.Duration `yaml:"period"`
  Cooldown time.Duration `yaml:"cooldown"`
}

var rdb *redis.RedisClient
var burst int64
var period int64
var rateLimit *Rule

func Init(rule *Rule) {
  rateLimit = rule
  rdb = redis.Client
  burst = rule.Burst
  period = int64(rule.Period.Seconds())
}
