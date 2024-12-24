package redis

import (
	"context"
	"time"
	_ "embed"

	redisv9 "github.com/redis/go-redis/v9"
)

// Generic Cell Rate Algorithm
// https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
type GCRA struct {
	Burst int
	Rate int
	Period time.Duration
}

type GCRAStatus struct {
}

// go:embed perform_gcra.lua
var PERFORM_GCRA_LUA string
// go:embed inspect_gcra.lua
var INSPECT_GCRA_LUA string

var PERFORM_GCRA_SCRIPT = redisv9.NewScript(PERFORM_GCRA_LUA)
var INSPECT_GCRA_SCRITP = redisv9.NewScript(INSPECT_GCRA_LUA)

func (gcra *GCRA) Perform(ctx context.Context, key string) (*GCRAStatus, error) {
	/*
	v, err := PERFORM_GCRA_SCRIPT.Run(ctx, Client, []string{key},
		[]interface{}{key, gcra.Burst, gcra.Rate, gcra.Period}).Result()
	if err != nil {
		return nil, err
	}
	values := v.([]interface{})
	*/

	// TODO

	return &GCRAStatus{}, nil
}

func (gcra *GCRA) Inspect(ctx context.Context) (*GCRAStatus, error) {

	// TODO

	return &GCRAStatus{}, nil
}
