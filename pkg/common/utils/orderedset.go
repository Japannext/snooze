package utils

type OrderedStringSet struct {
	a []string
	m map[string]bool
}

func NewOrderedStringSet() *OrderedStringSet {
	return &OrderedStringSet{
		a: []string{},
		m: map[string]bool{},
	}
}

func (set *OrderedStringSet) Append(item string) {
	if _, exists := set.m[item]; exists {
		return
	}
	set.m[item] = true
	set.a = append(set.a, item)
}

func (set *OrderedStringSet) AppendMany(items []string) {
	for _, item := range items {
		set.Append(item)
	}
}

func (set *OrderedStringSet) Items() []string {
	return set.a
}
