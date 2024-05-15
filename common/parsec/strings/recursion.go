package strings

import (
	"fmt"
	p "github.com/okneniz/parsec/common"
)

// Mark one combinator as subject to recursion.
// When staying at the same position, if the parsing loop
// reach twice this function, it will return an error,
// thus breaking the infinite loop.
func Recursion[T any](c p.Combinator[rune, Position, T]) p.Combinator[rune, Position, T] {
	return func(buffer p.Buffer[rune, Position]) (T, error) {
		pos := buffer.Position()
		pos.recursion++
		if pos.recursion > 1 {
			pos.recursion = 0
			return *new(T), fmt.Errorf("buffer is recursing")
		}
		return c(buffer)
	}
}
