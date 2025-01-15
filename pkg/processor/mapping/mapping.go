package mapping

import (
	"context"
	"fmt"

	"github.com/japannext/snooze/pkg/common/utils"
)

type ContextProcessor struct {
	mappings map[string]map[string]string
}

func New(cfg Config) (*ContextProcessor, error) {
	mappings := make(map[string]map[string]string)
	duplicate := utils.NewDuplicateChecker()

	for _, m := range cfg.Mappings {
		if err := duplicate.Check(m.Name); err != nil {
			return &ContextProcessor{}, fmt.Errorf("duplicate mapping '%s': %w", m.Name, err)
		}

		mappings[m.Name] = m.Map
	}

	return &ContextProcessor{mappings: mappings}, nil
}

type Config struct {
	Mappings []Mapping `json:"mappings" yaml:"mappings"`
}

type Mapping struct {
	Name string            `json:"name" yaml:"name"`
	Map  map[string]string `json:"map"  yaml:"map"`
}

func (p *ContextProcessor) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "mappings", p.mappings)
}
