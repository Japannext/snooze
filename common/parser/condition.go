package parser

import (
	"fmt"
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"

	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
)

type binopBuilder = func(*field.AlertField, string) *condition.Condition
type unopBuilder = func(*field.AlertField) *condition.Condition
type expression = p.Combinator[rune, RecursivePosition, *condition.Condition]

func binop_op() p.Combinator[rune, RecursivePosition, binopBuilder] {
	return Choice(
		Try(MatchAndReturn(String("=~"), condition.NewMatch)),
		Try(MatchAndReturn(String("!="), condition.NewNotEqual)),
		Try(MatchAndReturn(String("!~"), condition.NewNotMatch)),
		Try(MatchAndReturn(String("="), condition.NewEqual)),
	)
}

func quoted(quote rune) p.Combinator[rune, RecursivePosition, string] {
	return TextInBetween(Eq(quote), Any(), Eq(quote))
}

func str() p.Combinator[rune, RecursivePosition, string] {
	return Choice(
		Try(quoted('"')),
		Try(quoted('\'')),
		Try(quoted('`')),
	)
}

func SpacePadded[T any](x p.Combinator[rune, RecursivePosition, T]) p.Combinator[rune, RecursivePosition, T] {
	return Padded(Eq(' '), x)
}

func binop() expression {
	return func(buffer p.Buffer[rune, RecursivePosition]) (*condition.Condition, error) {
		fmt.Printf("[binop:%d] ...\n", buffer.RecursivePosition().Column())
		f, err := Field()(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		builder, err := SpacePadded(binop_op())(buffer)
		if err != nil {
			return new(condition.Condition), err
		}

		v, err := SpacePadded(str())(buffer)
		if err != nil {
			return new(condition.Condition), err
		}

		return builder(f, v), nil
	}
}

func unop_op() p.Combinator[rune, RecursivePosition, unopBuilder] {
	return Choice(
		Try(MatchAndReturn(String("has"), condition.NewHas)),
	)
}

func unop() expression {
	return func(buffer p.Buffer[rune, RecursivePosition]) (*condition.Condition, error) {
		fmt.Printf("[unop:%d] ...\n", buffer.RecursivePosition().Column())
		builder, err := unop_op()(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		f, err := Skip(Eq(' '), Field())(buffer)
		if err != nil {
			return new(condition.Condition), err
		}

		return builder(f), nil
	}
}

func operation() expression {
	return Choice(
		Try(unop()),
		Try(binop()),
	)
}

/*
func not(e expression) expression {
	return func(buffer p.Buffer[rune, RecursivePosition]) (*condition.Condition, error) {
		fmt.Printf("[not:%d] ---\n", buffer.RecursivePosition().Column())
		_, err := Eq('!')(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		x, err := SpacePadded(e)(buffer)
		if err != nil {
			return new(condition.Condition), err
		}
		return condition.NewNot(x), nil
	}
}
*/

type andorBuilder = func(*condition.Condition, *condition.Condition) *condition.Condition

func andorParser() p.Combinator[rune, RecursivePosition, andorBuilder] {
	return MapStrings(map[string]andorBuilder{
		"and": condition.NewAnd,
		"or":  condition.NewOr,
	})
}

/*
func andor(e expression) expression {
	return Recursing(
		Chainl1(SpacePadded(e), SpacePadded(andorParser())),
	)
}
*/

func te() expression {
	return Choice(
		Try(operation()),
		Try(not(expr())),
		Try(expr()),
	)
}

func term() expression {
	return Choice(
		Try(Parens(operation())),
		Try(operation()),
	)
}

func andor() expression {
	return Try(Sequence(andorParser(), term()))
}

func not() expression {
	return Case(
		Sequence(Eq('!'), SpacePadded(term())),
		func(c *condition.Condition) (*condition.Condition, error) {
			return condition.NewNot(c)
		},
	)
}

func expr() expression {
	return Choice(
		Try(Sequence(operator(), andor())),
	)
}

func expr() expression {
	return Chainl1(
		SpacePadded(te()),
		SpacePadded(andorParser()),
	)
}


func ConditionExpr() p.Combinator[rune, RecursivePosition, *condition.Condition] {
	var e func() expression
	// var te expression

	not := func() expression {
		return func(buffer p.Buffer[rune, RecursivePosition]) (*condition.Condition, error) {
			fmt.Printf("[not:%d] ...\n", buffer.RecursivePosition().Column())
			_, err := Eq('!')(buffer)
			if err != nil {
				return new(condition.Condition), err
			}
			x, err := SpacePadded(e())(buffer)
			if err != nil {
				return new(condition.Condition), err
			}
			return condition.NewNot(x), nil
		}
	}
	/*
		andor := func(buffer p.Buffer[rune, RecursivePosition]) (*condition.Condition, error) {
			fmt.Printf("[andor:%d] ...\n", buffer.RecursivePosition().Column())
			c1, err := te(buffer)
			if err != nil {
				return new(condition.Condition), err
			}
			builder, err := SpacePadded(andorParser())(buffer)
			if err != nil {
				return new(condition.Condition), err
			}
			c2, err := SpacePadded(e)(buffer)
			if err != nil {
				return new(condition.Condition), err
			}
			return builder(c1, c2), nil
		}
	*/

	andor := func() expression {
		return Chainl1[*condition.Condition](e(), andorParser())
	}

	e = func() expression {
		return Choice(
			Try(not()),
			Try(andor()),
		)
	}
	// Terminal expression.
	// Solve loop issues by trying the terminal expression (oeprator) first
	/*
	te = Choice(
		Try(operation()),
		Try(e),
	)
	*/

	return e()
}

func ParseCondition(text string) (*condition.Condition, error) {
	res, err := StrictParseString(text, expr())
	if err != nil {
		return new(condition.Condition), err
	}
	return res, nil
}
