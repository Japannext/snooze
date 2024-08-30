package lang



/*
import (
	"github.com/japannext/snooze/pkg/common/utils"
	api "github.com/japannext/snooze/pkg/common/api/v2"
)

// The fast mapper is an alternative to if/conditions.
// Rules that need to match a certain rule can register
// a simple condition (a switch), that is a key=value
// exact match. Then a map will be used to retrieve
// all the objects that registered to match this.

type Switch struct {
	Key string `yaml:"key"`
	Value string `yaml:"value"`
}

type switchable interface {
	GetSwitch() Switch
	comparable
}

type FastMapper[T switchable] struct {
	fields []Field
	mapper map[Switch][]T
}

func NewFastMapper[T switchable](items []T) *FastMapper[T] {
    mapper := make(map[Switch][]T)
	keys := utils.OrderedSetFrom[string](items)
    for _, item := range items {
		sw := item.GetSwitch()
        if _, found := mapper[sw]; !found {
            mapper[sw] = []T{}
        }
        mapper[sw] = append(mapper[sw], item)
    }
	fields := NewFields(keys.Items())
    return &FastMapper[T]{fields, mapper}
}

func (fast *FastMapper[T]) MatchLog(item api.Log) []T {
    for _, key := range fast.keys {
		sw := key.GetSwitch()
        value, ok := lang.ExtractFromLog(item, sw.Key)
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
*/
