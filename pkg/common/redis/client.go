package redis

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/tracing"
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
	tracerProvider := tracing.NewTracerProvider("redis")

	options, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}

	rdb := redisv9.NewClient(options)
	if err := redisotel.InstrumentTracing(rdb, redisotel.WithTracerProvider(tracerProvider)); err != nil {
		log.Fatalf("failed to instrument redis: %s", err)
	}

	Client = &RedisClient{rdb}
}
