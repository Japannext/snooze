package mapping

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/utils"
)

type ContextProcessor struct {
	mappings map[string]map[string]string
}

func New(cfgs []pipeline.Mapping) (*ContextProcessor, error) {
	mappings := make(map[string]map[string]string)
	duplicate := utils.NewDuplicateChecker()

	for _, cfg := range cfgs {
		if err := duplicate.Check(cfg.Name); err != nil {
			return &ContextProcessor{}, fmt.Errorf("duplicate mapping '%s': %w", cfg.Name, err)
		}

		mappings[cfg.Name] = cfg.Map
	}

	return &ContextProcessor{mappings: mappings}, nil
}

func (p *ContextProcessor) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "mappings", p.mappings)
}
