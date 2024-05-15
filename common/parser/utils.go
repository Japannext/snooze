package parser

import (
	"fmt"
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

type index struct {
	i int
}

type RecursivePosition struct {
	line   uint
	column uint
	index  int
	recursion int
}

// Line - line number.
func (p RecursivePosition) Line() uint {
	return p.line
}

// Column - column number.
func (p RecursivePosition) Column() uint {
	return p.column
}

// String - return string representation of opsition.
func (p RecursivePosition) String() string {
	return fmt.Sprintf("line=%d column=%d index=%d", p.line, p.column, p.index)
}

type buffer struct {
	data []rune
	position RecursivePosition
	newLineRunes map[rune]struct{}
	// Allow for combinator to add recusion info
	recursion int
}

func (b *buffer) ResetRecursion() {
	b.recursion = 0
}

func (b *buffer) Recurse() bool {
	b.position.recursion++
	// Went through this function twice at the same position
	if b.position.recursion > 1 {
		return true
	}
	return false
}

type RecursiveBuffer[T any, P any] interface {
	p.Buffer[T, P]
	ResetRecursion()
	Recurse() bool
}

type Combinator[T any, P any, S any] func(RecursiveBuffer[T, P]) (S, error)

func Recursing[T any](c p.Combinator[rune, Position, T]) Combinator[rune, Position, T] {
	return func(b RecursiveBuffer[rune, Position]) (T, error) {
		if b.Recurse() {
			b.ResetRecursion()
			return *new(T), fmt.Errorf("buffer is recursing")
		}
		return c(b)
	}
}

func (b *buffer) Read(greedy bool) (rune, error) {
	if b.IsEOF() {
		return 0, p.EndOfFile
	}

	x := b.data[b.position.index]

	if greedy {
		b.position.index++
		b.recursion = 0

		if _, isNewLine := b.newLineRunes[x]; isNewLine {
			b.position.column = 0
			b.position.line++
		} else {
			b.position.column++
		}
	}

	return x, nil
}

// Seek - change buffer position
func (b *buffer) Seek(x RecursivePosition) {
	b.position = x
}

// Position - return current buffer position
func (b *buffer) Position() RecursivePosition {
	return b.position
}

// IsEOF - true if buffer ended.
func (b *buffer) IsEOF() bool {
	return b.position.index >= len(b.data)
}

var defaultNewLineRunes = map[rune]struct{}{'\n': {}}

// Buffer - make buffer which can read text on input and use
// struct for positions.
func NewRecursiveBuffer(data []rune, newLineRunes ...rune) *buffer {
	b := new(buffer)
	b.data = data
	b.position = RecursivePosition{0, 0, 0, 0}

	if len(newLineRunes) == 0 {
		b.newLineRunes = defaultNewLineRunes
	} else {
		b.newLineRunes = make(map[rune]struct{})

		for _, x := range newLineRunes {
			b.newLineRunes[x] = struct{}{}
		}
	}

	return b
}

func StrictParseString[T any](str string, parse p.Combinator[rune, RecursivePosition, T]) (T, error) {
	buf := NewRecursiveBuffer([]rune(str))
	res, err := p.Parse[rune, RecursivePosition, T](buf, parse)
	if err != nil {
		return *new(T), err
	}
	if !buf.IsEOF() {
		leftover := str[buf.Position().Column() : len(str)-1]
		return *new(T), fmt.Errorf("failed to parse: %s", leftover)
	}
	return res, nil
}

// Match one combinator, and return another thing
func MatchAndReturn[T any](c p.Combinator[rune, Position, string], f T) p.Combinator[rune, Position, T] {
	return func(buffer p.Buffer[rune, Position]) (T, error) {
		_, err := c(buffer)
		if err != nil {
			return *new(T), err
		}
		return f, nil
	}
}

func SuperChain[T any, P any, S any](
	c p.Combinator[T, P, S],
	op p.Combinator[T, P, func(S, S) S],
	prefix p.Combinator[T, P, func(S, S) S],
) p.Combinator[T, P, S] {
	return func(buffer p.Buffer[T, P]) (S, error) {
		x, err := c(buffer)
		if err != nil {
			return *new(S), err
		}

		rest := x

		for !buffer.IsEOF() {
			f, err := op(buffer)
			if err != nil {
				break
			}

			y, err := c(buffer)
			if err != nil {
				break
			}

			rest = f(rest, y)
		}

		return rest, nil
	}
}


func ManyUntil[T any, P any, S any, B any](cap int, c p.Combinator[T, P, S], end p.Combinator[T, P, B]) p.Combinator[T, P, []S] {
	z := p.Try(end)
	return func(buffer p.Buffer[T, P]) ([]S, error) {
		result := make([]S, 0, cap)

		for {
			_, err := z(buffer)
			if err == nil {
				break
			}

			x, err := c(buffer)
			if err != nil {
				return result, err
			}

			result = append(result, x)
		}

		return result, nil
	}
}

func TextInBetween[S any, E any](pre p.Combinator[rune, Position, S], c p.Combinator[rune, Position, rune], suf p.Combinator[rune, Position, E]) p.Combinator[rune, Position, string] {
	return func(buffer p.Buffer[rune, Position]) (string, error) {
		_, err := pre(buffer)
		if err != nil {
			return "", err
		}
		res, err := ManyUntil(0, Try(c), suf)(buffer)
		if err != nil {
			return "", err
		}
		return string(res), nil
	}
}
