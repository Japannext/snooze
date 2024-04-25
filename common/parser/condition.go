package parser

import (
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"

	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
)

type binopBuilder = func(*field.AlertField, string) *condition.Condition
type unopBuilder = func(*field.AlertField) *condition.Condition
type expression = p.Combinator[rune, Position, *condition.Condition]

func binop_op() p.Combinator[rune, Position, binopBuilder] {
	return Choice(
		Try(MatchAndReturn(String("=~"), condition.NewMatch)),
		Try(MatchAndReturn(String("!="), condition.NewNotEqual)),
		Try(MatchAndReturn(String("!~"), condition.NewNotMatch)),
		Try(MatchAndReturn(String("="), condition.NewEqual)),
	)
}

func quoted(quote rune) p.Combinator[rune, Position, string] {
	return TextInBetween(Eq(quote), Any(), Eq(quote))
}

func str() p.Combinator[rune, Position, string] {
	return Choice(
		Try(quoted('"')),
		Try(quoted('\'')),
		Try(quoted('`')),
	)
}

func SpacePadded[T any](x p.Combinator[rune, Position, T]) p.Combinator[rune, Position, T] {
	return Padded(Eq(' '), x)
}

func binop() expression {
	return func(buffer p.Buffer[rune, Position]) (*condition.Condition, error) {
		f, err := Field()(buffer)
		if err != nil {
			return nil, err
		}
		builder, err := SpacePadded(binop_op())(buffer)
		if err != nil {
			return nil, err
		}

		v, err := SpacePadded(str())(buffer)
		if err != nil {
			return nil, err
		}

		return builder(f, v), nil
	}
}

func unop_op() p.Combinator[rune, Position, unopBuilder] {
	return Choice(
		Try(MatchAndReturn(String("has"), condition.NewHas)),
	)
}

func unop() expression {
	return func(buffer p.Buffer[rune, Position]) (*condition.Condition, error) {
		builder, err := unop_op()(buffer)
		if err != nil {
			return nil, err
		}
        f, err := Skip(Eq(' '), Field())(buffer)
        if err != nil {
            return nil, err
        }

		return builder(f), nil
	}
}

func operation() expression {
	return Choice(
		Try(binop()),
		Try(unop()),
	)
}

/*
func not(e expression) expression {
	return Cast(Skip(Eq('!'), SpacePadded(e)),
		func(x *condition.Condition) (*condition.Condition, error) {
			return condition.NewNot(x), nil
		})
}
*/
func not(e expression) expression {
	return func(buffer p.Buffer[rune, Position]) (*condition.Condition, error) {
		_, err := Eq('!')(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		c, err := SpacePadded(e)(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		return condition.NewNot(c), nil
	}
}

type andorBuilder = func(...*condition.Condition) *condition.Condition

func andorParser() p.Combinator[rune, Position, andorBuilder] {
	return MapStrings(map[string]andorBuilder{
		"and": condition.NewAnd,
		"or": condition.NewOr,
	})
}

func andor(e expression) p.Combinator[rune, Position, *condition.Condition] {
	return func(buffer p.Buffer[rune, Position]) (*condition.Condition, error) {
		c1, err := e(buffer)
		if err != nil {
			return nil, err
		}
		builder, err := SpacePadded(andorParser())(buffer)
		if err != nil {
			return nil, err
		}
		c2, err := SpacePadded(e)(buffer)
		if err != nil {
			return nil, err
		}
		return builder(c1, c2), nil
	}
}


func pexpr(e expression) p.Combinator[rune, Position, *condition.Condition] {
	return Parens(e)
}

func ConditionExpr() p.Combinator[rune, Position, *condition.Condition] {
	var e expression
	not := not(e)
	andor := andor(e)
	pexpr := pexpr(e)
	e = Choice(
		Try(operation()),
		Try(not),
		Try(andor),
		Try(pexpr),
	)
	return e
}

func ParseCondition(text string) (*condition.Condition, error) {
	res, err := StrictParseString(text, ConditionExpr())
	if err != nil {
		return nil, err
	}
	return res, nil
}
