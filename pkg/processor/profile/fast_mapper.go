package profile

import (
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/models"
)

type Kv struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type FastMapper struct {
	keys   []string
	mapper map[Kv][]*Profile
	fields map[string]*lang.Field
}

func NewFastMapper(prfs []*Profile) *FastMapper {
	m := map[Kv][]*Profile{}
	keys := utils.NewOrderedSet[string]()
	fields := map[string]*lang.Field{}
	var err error
	for _, prf := range prfs {
		prf.Load()
		m[prf.Switch] = append(m[prf.Switch], prf)
		fields[prf.Switch.Key], err = lang.NewField(prf.Switch.Key)
		if err != nil {
			log.Fatalf("invalid field `%s`", prf.Switch.Key)
		}
		keys.Append(prf.Switch.Key)
	}
	return &FastMapper{keys.Items(), m, fields}
}

func (m *FastMapper) GetMatches(item *models.Log) []*Profile {
	for _, key := range m.keys {
		field, found := m.fields[key]
		if !found {
			log.Warnf("unexpected field `%s`", key)
			continue
		}
		value, err := lang.ExtractField(item, field)
		if err != nil {
			continue
		}
		prfs, found := m.mapper[Kv{key, value}]
		if !found {
			continue
		}
		return prfs
	}
	return []*Profile{}
}
