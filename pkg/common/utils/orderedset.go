package utils

type OrderedSet[T comparable] struct {
	a []T
	m map[T]bool
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		a: []T{},
		m: map[T]bool{},
	}
}

func OrderedSetFrom[T comparable](items []T) *OrderedSet[T] {
	set := NewOrderedSet[T]()
	set.AppendMany(items)
	return set
}

func (set *OrderedSet[T]) Append(item T) {
	if _, exists := set.m[item]; exists {
		return
	}
	set.m[item] = true
	set.a = append(set.a, item)
}

func (set *OrderedSet[T]) AppendMany(items []T) {
	for _, item := range items {
		set.Append(item)
	}
}

func (set *OrderedSet[T]) Items() []T {
	return set.a
}
