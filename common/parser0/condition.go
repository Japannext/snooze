package parser

import (
	"fmt"

	. "github.com/prataprc/goparsec"

	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
)

var conditionParser Parser

func initConditionParser() {
	// Literals
	ff := fieldParser
	str := String()

	// Binary op
	eq := Atom(`=`, "EQUAL")
	neq := Atom(`!=`, "NOT_EQUAL")
	match := Atom("=~", "MATCH")
	nmatch := Atom(`!=`, "NOT_MATCH")
	has := Atom(`has`, "HAS")

	bin_tk := OrdChoice(f_one, match, nmatch, eq, neq)
	bin_op := And(f_bin, ff, bin_tk, str)
	op := OrdChoice(f_one,
		And(f_has, has, ff),
		bin_op,
	)

	// Circular reference
	var expr Parser

	not := Atom("not", "NOT")
	not_op := And(f_not, not, &expr)

	and := Atom("and", "AND")
	or := Atom("or", "OR")
	andor := OrdChoice(f_one, and, or)
	andor_op := And(f_andor, &expr, andor, &expr)

	spth := Atom(`(`, "START_PARENS")
	epth := Atom(`)`, "END_PARENS")
	pexpr := And(f_pth, spth, &expr, epth)

	expr = OrdChoice(f_one,
		op,
		not_op,
		andor_op,
		pexpr,
	)

	conditionParser = OrdChoice(f_one, expr)
}

func f_bin(ns []ParsecNode) ParsecNode {
	fi := ns[0].(*field.AlertField)
	op := ns[1].(*Terminal).GetName()
	v := ns[2].(string)
	fmt.Printf("[bin] fi: %#v, op: %#v, v: %#v\n", fi, op, v)
	switch op {
	case "EQUAL":
		return condition.NewEqual(fi, v)
	case "NOT_EQUAL":
		return condition.NewNotEqual(fi, v)
	case "MATCH":
		return condition.NewMatch(fi, v)
	case "NOT_MATCH":
		return condition.NewNotMatch(fi, v)
	}
	panic(fmt.Sprintf("op: %s", op))
}

func f_has(ns []ParsecNode) ParsecNode {
	fi := ns[1].(*field.AlertField)
	fmt.Printf("[has] fi: %#v\n", fi)
	return condition.NewHas(fi)
}

func f_pth(ns []ParsecNode) ParsecNode {
	fmt.Printf("[f_pth] %#v\n", ns)
	if len(ns) == 0 {
		return nil
	}
	return ns[1]
}

func f_andor(ns []ParsecNode) ParsecNode {
	c1 := ns[0].(*condition.Condition)
	op := ns[1].(*Terminal).GetName()
	c2 := ns[2].(*condition.Condition)
	fmt.Printf("[andor] c1: %#v, op: %#v, c2: %#v\n", c1, op, c2)
	switch op {
	case "AND":
		return condition.NewAnd(c1, c2)
	case "OR":
		return condition.NewOr(c1, c2)
	}
	panic(fmt.Sprintf("op: %s", op))
}

func f_not(ns []ParsecNode) ParsecNode {
	c := ns[1].(*condition.Condition)
	fmt.Printf("[f_not] %#v\n", ns)
	return condition.NewNot(c)
}

func ParseCondition(text string) (*condition.Condition, error) {
	s := NewScanner([]byte(text))
	res, next := conditionParser(s)
	if !next.Endof() {
		return nil, fmt.Errorf("failed to parse '%s' char number %d", text, next.GetCursor())
	}
	c, ok := res.(*condition.Condition)
	if !ok {
		return nil, fmt.Errorf("no condition matched (`%s`)", text)
	}
	return c, nil
}
