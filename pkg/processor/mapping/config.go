package mapping

import (
	"context"
)

type Mapping struct {
	Name string            `yaml:"name" json:"name"`
	Map  map[string]string `yaml:"map" json:"map"`
}

func WithMappings(ctx context.Context) context.Context {
	return context.WithValue(ctx, "mappings", mappings)
}
