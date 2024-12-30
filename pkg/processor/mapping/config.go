package mapping

import (
	"context"
)

type Mapping struct {
	Name string            `json:"name" yaml:"name"`
	Map  map[string]string `json:"map"  yaml:"map"`
}

func WithMappings(ctx context.Context) context.Context {
	return context.WithValue(ctx, "mappings", mappings)
}
