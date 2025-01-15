package redis

import (
	"context"
)

type TestKeys struct {
	kvs map[string]string
}

func NewKeys(kvs map[string]string) *TestKeys {
	tks := &TestKeys{kvs: kvs}

	args := []string{}
	for key, value := range kvs {
		log.Debugf("setting %s=%s", key, value)
		args = append(args, key, value)
	}

	Client.MSet(context.Background(), args).Err()

	return tks
}

func (tks *TestKeys) Cleanup() {
	ctx := context.Background()

	keys := []string{}
	for key, _ := range tks.kvs {
		log.Debugf("deleting %s", key)
		keys = append(keys, key)
	}

	Client.Del(ctx, keys...).Err()
}
