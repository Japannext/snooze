package redis_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	redis.Init()
	os.Exit(m.Run())
}

func burst(gcra *redis.GCRA, key string, b int) error {
	ctx := context.TODO()

	for i := range b {
		if _, err := gcra.Perform(ctx, key); err != nil {
			return fmt.Errorf("error in burst #%d: %w", i, err)
		}
	}

	return nil
}

func TestGCRA(t *testing.T) {
	tests := []struct{
		Name string
		GCRA *redis.GCRA
		Burst, Remaining int
		Limited bool
	}{
		{
			Name: "under limit",
			GCRA: &redis.GCRA{Burst: 1000, Rate: 10, Period: time.Duration(1)*time.Second},
			Burst: 1,
			Limited: false,
			Remaining: 0,
		},
		{
			Name: "over limit",
			GCRA: &redis.GCRA{Burst: 10, Rate: 1, Period: time.Duration(1)*time.Second},
			Burst: 100,
			Limited: true,
			Remaining: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			key := "ratelimit:" + tt.Name
			err := burst(tt.GCRA, key, tt.Burst)
			require.NoError(t, err)
			status, err := tt.GCRA.Inspect(context.TODO(), key)
			require.NoError(t, err)
			assert.Equal(t, tt.Limited, status.Limited)
			assert.Equal(t, tt.Remaining, status.Remaining)
		})
	}
}
