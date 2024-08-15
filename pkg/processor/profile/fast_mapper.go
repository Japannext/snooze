package profile

import (
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type FastMapper struct {
	keys []string
	mapper map[Kv][]Rule
}

// Return a set of keys, ordered by order of appearance
// in rules
func orderedSetOfKeys(rules []Rule) (keys []string) {
	inserted := make(map[string]bool)
	for _, rule := range rules {
		key := rule.Switch.Key
		if !inserted[key] {
			keys = append(keys, key)
		}
	}
	return
}

func NewFastMapper(rules []Rule) *FastMapper {
	m := make(map[Kv][]Rule)
	keys := orderedSetOfKeys(rules)
	for _, rule := range rules {
		if _, ok := m[rule.Switch]; !ok {
			m[rule.Switch] = []Rule{}
		}
		m[rule.Switch] = append(m[rule.Switch], rule)
	}
	return &FastMapper{keys, m}
}

func (m *FastMapper) GetRules(item *api.Log) (rules []Rule) {
	for _, key := range m.keys {
		value, ok := FindValue(key, item)
		if !ok {
			continue
		}
		rules, ok = m.mapper[Kv{key, value}]
		if !ok {
			continue
		}
		return
	}
	return
}
