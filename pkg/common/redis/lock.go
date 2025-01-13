package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	redisv9 "github.com/redis/go-redis/v9"
)

const (
	Nx = "nx"
	lockTimeout = 10 * time.Second
)

type Mutex struct {
	name string
	id   string
	done chan struct{}
}

func NewMutex(name string) *Mutex {
	hostname, _ := os.Hostname()
	id := fmt.Sprintf("%s:%s", hostname, uuid.NewString())

	return &Mutex{
		name: name,
		id: id,
		done: make(chan struct{}),
	}
}

func (m *Mutex) keepLocked(ctx context.Context) {
	for {
		select {
		case <-m.done:
			return
		default:
			id, err := Client.Get(ctx, m.name).Result()
			if err != nil {
				log.Errorf("failed to check lock value: %s", err)

				break
			}

			if id != m.id {
				log.Warnf("lost the lock, now held by %s", id)

				return
			}

			expireAt := time.Now().Add(lockTimeout)

			if err := Client.ExpireAt(ctx, m.name, expireAt).Err(); err != nil {
				log.Errorf("failed to keep lock: %s", err)

				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *Mutex) tryAcquire(ctx context.Context) (bool, error) {
	cmd := Client.SetArgs(ctx, m.name, m.id, redisv9.SetArgs{
		Mode: Nx,
		Get: true,
		TTL: lockTimeout,
	})
	id, err := cmd.Result()

	// Lock successful
	if errors.Is(err, redisv9.Nil) {
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("during `%s`: %s", cmd, err)
	}

	// Return positive result if we were previously the master
	if id == m.id {
		return true, nil
	}

	return false, nil
}

// Return a channel that signals when the lock is acquired.
func (m *Mutex) Acquire(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if ctx.Err() != nil {
				return ctx.Err()
			}

			primary, err := m.tryAcquire(ctx)
			if err != nil {
				log.Errorf("failed to acquire lock: %s", err)

				continue
			}

			if primary {
				go m.keepLocked(ctx)

				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *Mutex) Release() {
	ctx := context.Background()
	// Stop keeping the lock
	m.done <- struct{}{}
	// Check if the lock is still held by self
	id, err := Client.Get(ctx, m.name).Result()
	if err != nil {
		log.Errorf("failed to check lock value: %s", err)

		return
	}

	if id != m.id {
		log.Warnf("lost the lock, now held by %s", id)

		return
	}

	if err := Client.Del(ctx, m.name).Err(); err != nil {
		log.Errorf("failed to delete lock: %s", err)

		return
	}
}
