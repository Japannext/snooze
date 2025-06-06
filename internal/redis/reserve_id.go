package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	redisv9 "github.com/redis/go-redis/v9"
)

// Generate a UUID, and attempt to reserve a given key.
// If the key is already taken, return the ID that has
// been reserved by another process.
func ReserveID(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {

	u, err := uuid.NewRandom()
	if err != nil {
		return "", false, fmt.Errorf("failed to generate UID")
	}
	newID := u.String()

	id, err := Client.SetArgs(ctx, key, newID, redisv9.SetArgs{
		Mode: "NX",
		Get:  true,
		TTL:  ttl,
	}).Result()

	// The ID is ours.
	if IsNil(err) || id == newID {
		return newID, false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to reserve ID: %w", err)
	}

	return id, true, nil
}
