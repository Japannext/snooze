package kv

import (
	"github.com/japannext/snooze/internal/orderedset"
	"github.com/japannext/snooze/internal/pipeline"
)

type HasLabels interface {
	GetLabels() map[string]string
}

type HasSelectors interface {
	GetSelectors() []pipeline.Selector
}

type Kv[K HasLabels, V HasSelectors] struct {
	keys   []string
	mapper map[pipeline.Selector][]*V
}

func New[K HasLabels, V HasSelectors](results []V) (*Kv[K, V], error) {
	m := map[pipeline.Selector][]*V{}
	keys := orderedset.New[string]()

	for _, result := range results {
		selectors := result.GetSelectors()
		for _, selector := range selectors {
			if _, exist := m[selector]; !exist {
				m[selector] = []*V{}
			}
			m[selector] = append(m[selector], &result)
			keys.Append(selector.Key)
		}
	}

	return &Kv[K, V]{
		keys:   keys.Items(),
		mapper: m,
	}, nil
}

func (kv *Kv[K, V]) FirstMatch(item K) (*V, bool) {
	labels := item.GetLabels()
	for _, key := range kv.keys {
		value, exist := labels[key]
		if !exist {
			continue
		}

		selector := pipeline.Selector{Key: key, Value: value}
		results, found := kv.mapper[selector]
		if !found || len(results) == 0 {
			continue
		}

		return results[0], true
	}

	return nil, false
}

func (kv *Kv[K, V]) AllMatches(item K) []*V {
	results := []*V{}
	labels := item.GetLabels()
	for _, key := range kv.keys {
		value, exist := labels[key]
		if !exist {
			continue
		}

		selector := pipeline.Selector{Key: key, Value: value}
		rr, found := kv.mapper[selector]
		if !found {
			continue
		}
		results = append(results, rr...)
	}

	return results
}
