package parser

import (
	"fmt"
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

func StrictParseString[T any](str string, parse p.Combinator[rune, Position, T]) (T, error) {
	buf := Buffer([]rune(str))
	res, err := p.Parse[rune, Position, T](buf, parse)
	if err != nil {
		return *new(T), err
	}
	if !buf.IsEOF() {
		leftover := str[buf.Position().Column():len(str)-1]
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
