package parser

import (
	"fmt"
	"unicode"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"

	"github.com/japannext/snooze/common/field"
)

func allowedChar(r rune) bool {
	return unicode.IsLetter(r) ||
		unicode.IsNumber(r) ||
		r == '.' ||
		r == '/' ||
		r == '-' ||
		r == '_'
}

func label() p.Combinator[rune, Position, string] {
	return Cast(
		Many(255, Try(p.Satisfy[rune, Position](true, allowedChar))),
		func(x []rune) (string, error) {
			if len(x) == 0 {
				return "", fmt.Errorf("field is empty")
			}
			return string(x), nil
		},
	)
}

func lonefield() p.Combinator[rune, Position, *field.AlertField] {
	return Cast(
		label(),
		func(x string) (*field.AlertField, error) {
			return field.NewAlertField(x, ""), nil
		},
	)
}

func multifield() p.Combinator[rune, Position, *field.AlertField] {
	return Cast(
		Sequence(2, label(), Between(Eq('['), label(), Eq(']'))),
		func (x []string) (*field.AlertField, error) {
			return field.NewAlertField(x[0], x[1]), nil
		},
	)
}

func Field() p.Combinator[rune, Position, *field.AlertField] {
	return Choice(
			Try(multifield()),
			Try(lonefield()),
	)
}

func ParseField(text string) (*field.AlertField, error) {
	res, err := StrictParseString(text, Field())
	if err != nil {
		return nil, err
	}
	return res, nil
}
