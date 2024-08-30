package utils

type Set[T comparable] struct {
	m map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: map[T]bool{},
	}
}

func SetFrom[T comparable](items []T) *Set[T] {
	set := NewSet[T]()
	set.AppendMany(items)
	return set
}

func (set *Set[T]) Append(item T) {
	if _, exists := set.m[item]; exists {
		return
	}
	set.m[item] = true
}

func (set *Set[T]) Contain(item T) bool {
	_, found := set.m[item]
	return found
}

func (set *Set[T]) AppendMany(items []T) {
	for _, item := range items {
		set.Append(item)
	}
}

func (set *Set[T]) Items() []T {
	keys := []T{}
	for key := range set.m {
		keys = append(keys, key)
	}
	return keys
}
