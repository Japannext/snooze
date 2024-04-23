package redis

import (
	redisv9 "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var (
	Client *RedisClient
)

var Nil = redisv9.Nil

type RedisClient struct {
	*redisv9.Client
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
