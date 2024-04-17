package redis

import (
  log "github.com/sirupsen/logrus"
  redisv9 "github.com/redis/go-redis/v9"
)

var (
  Client *RedisClient
)

type RedisClient struct {
  Client *redisv9.Client
}

func Init() {

  options, err := initConfig()
  if err != nil {
    log.Fatal(err)
  }

  Client = &RedisClient{
    redisv9.NewClient(options),
  }
}
