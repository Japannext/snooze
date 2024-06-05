package redis

import (
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	Client *RedisClient
)

var Nil = redisv9.Nil

type RedisClient struct {
	*redisv9.Client
}

var log *logrus.Entry

func Init() {

	log = logrus.WithFields(logrus.Fields{"module": "redis"})

	options, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	Client = &RedisClient{
		redisv9.NewClient(options),
	}
}
