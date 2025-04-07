package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
	_ "embed"

	redisv9 "github.com/redis/go-redis/v9"
)

// Generic Cell Rate Algorithm
// https://en.wikipedia.org/wiki/Generic_cell_rate_algorithm
type GCRA struct {
	Burst  int
	Rate   int
	Period time.Duration
}

type GCRAStatus struct{
	Limited bool
	Remaining int
	RetryAfter time.Duration
	ResetAfter time.Duration
}

//go:embed perform_gcra.lua
var PerformGCRALua string

//go:embed inspect_gcra.lua
var InspectGCRALua string

//nolint:gochecknoglobals
var (
	PerformGCRAScript = redisv9.NewScript(PerformGCRALua)
	InspectGCRAScript = redisv9.NewScript(InspectGCRALua)
)

// Should never happen, except when there is a bug in the Lua script,
// or in the code that evaluate it.
func ScriptError(msg string, args ...interface{}) error {
	return fmt.Errorf("SCRIPT ERROR: " + msg, args)
}

func (gcra *GCRA) Perform(ctx context.Context, key string) (*GCRAStatus, error) {
	keys := []string{key}
	args := []interface{}{
		gcra.Burst,             // burst
		gcra.Rate,              // rate
		gcra.Period.Seconds(),  // period
		1,                      // cost
	}

	v, err := PerformGCRAScript.Run(ctx, Client, keys, args).Result()
	if IsError(err) {
		return &GCRAStatus{}, fmt.Errorf("error running perform_gcra script: %w", err)
	}

	status, err := scanStatus(v)
	if err != nil {
		return  &GCRAStatus{}, fmt.Errorf("error in response: %w", err)
	}

	return status, nil
}

func toDuration(f float64) time.Duration {
	if f == -1 {
		return -1
	}

	return time.Duration(f * float64(time.Second))
}

const gcraStatusArgsNb = 4

// Scan the output of perform_gcra and inspect_gcra scripts.
func scanStatus(v interface{}) (*GCRAStatus, error) {
	values, ok := v.([]interface{})
	if !ok {
		return &GCRAStatus{}, fmt.Errorf("invalid return value")
	}

	if len(values) != gcraStatusArgsNb {
		return &GCRAStatus{}, fmt.Errorf("wrong number of return args: got %d, expected %d", len(values), gcraStatusArgsNb)
	}

	var limited bool

	limitedString, ok := values[0].(string)

	switch {
	case !ok:
		return &GCRAStatus{}, fmt.Errorf("invalid value for 'limited' `%s` (not string)", values[0])
	case limitedString == "true":
		limited = true
	case limitedString == "false":
		limited = false
	default:
		return &GCRAStatus{}, fmt.Errorf("invalid value for 'limited' `%s` (should be 'true' or 'false')", values[0])
	}

	if limitedString == "true" {
		limited = true
	}

	remaining, ok := values[1].(int64)
	if !ok {
		return &GCRAStatus{}, fmt.Errorf("invalid value `%s` (not int64)", values[1])
	}

	retryAfter, err := strconv.ParseFloat(values[2].(string), 64)
	if err != nil {
		return &GCRAStatus{}, fmt.Errorf("invalid value `%s`: %w", values[2], err)
	}

	resetAfter, err := strconv.ParseFloat(values[3].(string), 64)
	if err != nil {
		return &GCRAStatus{}, fmt.Errorf("invalid value `%s`: %w", values[3], err)
	}

	return &GCRAStatus{
		Limited: limited,
		Remaining: int(remaining),
		RetryAfter: toDuration(retryAfter),
		ResetAfter: toDuration(resetAfter),
	}, nil
}

func (gcra *GCRA) Inspect(ctx context.Context, key string) (*GCRAStatus, error) {
	keys := []string{key}
	args := []interface{}{
		gcra.Burst,             // burst
		gcra.Rate,              // rate
		gcra.Period.Seconds(),  // period
		1,                      // cost
	}

	v, err := InspectGCRAScript.Run(ctx, Client, keys, args).Result()
	if IsError(err) {
		return &GCRAStatus{}, fmt.Errorf("error running perform_gcra script: %w", err)
	}

	status, err := scanStatus(v)
	if err != nil {
		return  &GCRAStatus{}, fmt.Errorf("error in response: %w", err)
	}

	return status, nil
}
