package kv

import (
	"fmt"

	"github.com/japannext/snooze/internal/hash"
	"github.com/japannext/snooze/internal/models"
)

type Kv[T any] struct {
	keys   []string
	mapper map[string][]*T
}

func New[T any]() *Kv[T] {
	return &Kv[T]{
		keys:   []string{},
		mapper: map[string][]*T{},
	}
}

func (kv *Kv[T]) Add(selector models.Selector, value *T) error {
	h := hash.HashSelector(selector)
	if h == "" {
		return fmt.Errorf("empty selector")
	}

	// Initialize the array if no exist
	if _, exist := kv.mapper[h]; !exist {
		kv.mapper[h] = []*T{}
	}

	kv.mapper[h] = append(kv.mapper[h], value)

	return nil
}

func (kv *Kv[T]) FirstMatch(key string) (*T, bool) {
	matches := kv.mapper[key]
	if len(matches) == 0 {
		return new(T), false
	}
	return matches[0], true
}

func (kv *Kv[T]) AllMatches(key string) []*T {
	matches, ok := kv.mapper[key]
	if !ok {
		return []*T{}
	}
	return matches
}

// func (kv *Kv[K, V]) FirstMatch(item K) (*V, bool) {
// 	labels := item.GetLabels()
// 	log.Debugf("labels: %+v", labels)
// 	log.Debugf("kv.keys: %+v", kv.keys)
// 	for _, key := range kv.keys {
// 		value, exist := labels[key]
// 		if !exist {
// 			continue
// 		}

// 		log.Debugf("found key %s", key)
// 		selector := pipeline.Selector{Key: key, Value: value}
// 		results, found := kv.mapper[selector]
// 		if !found || len(results) == 0 {
// 			continue
// 		}

// 		log.Debugf("found selector for key %s: %+v", key, selector)
// 		return results[0], true
// 	}

// 	return nil, false
// }

// func (kv *Kv[K, V]) AllMatches(item K) []*V {
// 	results := []*V{}
// 	labels := item.GetLabels()
// 	log.Debugf("labels: %+v", labels)
// 	log.Debugf("kv.keys: %+v", kv.keys)
// 	for _, key := range kv.keys {
// 		value, exist := labels[key]
// 		if !exist {
// 			continue
// 		}

// 		log.Debugf("found key %s", key)
// 		selector := pipeline.Selector{Key: key, Value: value}
// 		rr, found := kv.mapper[selector]
// 		if !found {
// 			continue
// 		}

// 		log.Debugf("found selector for key %s: %+v", key, selector)
// 		results = append(results, rr...)
// 	}

// 	return results
// }
