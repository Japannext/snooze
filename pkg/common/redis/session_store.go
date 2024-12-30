package redis

import (
	"context"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v9"
)

func GetSessionStore() sessions.Store {
	ctx := context.Background()
	store, err := redisstore.NewRedisStore(ctx, Client)
	if err != nil {
		log.Fatalf("failed to initialize redis store")
	}
	return store
}
