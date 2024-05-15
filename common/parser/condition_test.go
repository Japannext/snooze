package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

var _f = field.NewAlertField

type parsingTest struct {
	Text     string
	Expected interface{}
	Success  bool
}

func strictParse[T any](t *testing.T, parser p.Combinator[rune, Position, T], tests []parsingTest) {
	for _, data := range tests {
		t.Run(data.Text, func(t *testing.T) {
			res, err := StrictParseString(data.Text, parser)
			if data.Success {
				assert.NoError(t, err)
				assert.Equal(t, data.Expected, res)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestStr(t *testing.T) {
	// parser := Trace(t, "str", str())
	parser := str()
	tests := []parsingTest{
		{`"abc"`, "abc", true},
		{`'abc123'`, "abc123", true},
		{"`abc123`", "abc123", true},
		{`"abc`, "", false},
		{`abc123`, "", false},
	}
	strictParse(t, parser, tests)
}

func TestBinOp(t *testing.T) {
	parser := binop()
	tests := []parsingTest{
		{`labels[host.name] = "host01"`, condition.NewEqual(_f("labels", "host.name"), "host01"), true},
		{`labels[host.name] != "host01"`, condition.NewNotEqual(_f("labels", "host.name"), "host01"), true},
		{`severity =~ 'error'`, condition.NewMatch(_f("severity", ""), "error"), true},
		{`body[message] !~ "fatal"`, condition.NewNotMatch(_f("body", "message"), "fatal"), true},
		{`labels[host.name] =`, nil, false},
		{`labels[host.name]`, nil, false},
		{`severity`, nil, false},
		{`severity = 'test`, nil, false},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}

func TestUnOp(t *testing.T) {
	parser := unop()
	tests := []parsingTest{
		{`has labels[host.name]`, condition.NewHas(_f("labels", "host.name")), true},
		{`has severity_text`, condition.NewHas(_f("severity_text", "")), true},
		{`hastimestamp`, nil, false},
		{`has `, nil, false},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}

func TestOp(t *testing.T) {
	parser := operation()
	tests := []parsingTest{
		{`labels[host.name] = "host01"`, condition.NewEqual(_f("labels", "host.name"), "host01"), true},
		{`labels[host.name] != "host01"`, condition.NewNotEqual(_f("labels", "host.name"), "host01"), true},
		{`severity =~ 'error'`, condition.NewMatch(_f("severity", ""), "error"), true},
		{`body[message] !~ "fatal"`, condition.NewNotMatch(_f("body", "message"), "fatal"), true},
		{`labels[host.name] =`, nil, false},
		{`labels[host.name]`, nil, false},
		{`severity`, nil, false},
		{`severity = 'test`, nil, false},
		{`has labels[host.name]`, condition.NewHas(_f("labels", "host.name")), true},
		{`has severity_text`, condition.NewHas(_f("severity_text", "")), true},
		{`hastimestamp`, nil, false},
		{`has `, nil, false},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}

func TestNot(t *testing.T) {
	parser := Trace(t, "not", not(operation()))
	tests := []parsingTest{
		{`! labels[host.name] = "host01"`, condition.NewNot(condition.NewEqual(_f("labels", "host.name"), "host01")), true},
		{`! has severity_text`, condition.NewNot(condition.NewHas(_f("severity_text", ""))), true},
		{`! `, nil, false},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}

func TestAndOr(t *testing.T) {
	parser := Trace(t, "andor", andor(operation()))
	tests := []parsingTest{
		{`severity_text = "error" and severity_number = "12"`, condition.NewAnd(
			condition.NewEqual(_f("severity_text", ""), "error"),
			condition.NewEqual(_f("severity_number", ""), "12"),
		), true},
		{`labels[host.name] = "host01" or labels[host.name] = "host02"`, condition.NewOr(
			condition.NewEqual(_f("labels", "host.name"), "host01"),
			condition.NewEqual(_f("labels", "host.name"), "host02"),
		), true},
		{`labels[host.name]`, nil, false},
		{`labels[host.name] = "host01"`, condition.NewEqual(_f("labels", "host.name"), "host01"), true},
		{`labels[host.name] = "host01" and labels[host.name]`, nil, false},
		{`labels[host.name] = "host01" and labels[host.name] =`, nil, false},
		{`! `, nil, false},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}

func _operation(t *testing.T) expression {
    return Trace(t, "operation", Choice(
        Try(unop()),
        Try(binop()),
    ))
}

func _te(t *testing.T) expression {
    return Trace(t, "te", Choice(
        Try(_operation(t)),
		Try(_pexpr(t)),
		Try(not(_expr(t))),
    ))
}

func _pexpr(t *testing.T) expression {
	return Trace(t, "pexpr", Parens(
		_te(t),
	))
}

func _expr(t *testing.T) expression {
    return Trace(t, "expr", Chainl1(
        SpacePadded(_operation(t)),
        SpacePadded(andorParser()),
    ))
}

type fop func(string, string) string

func op() p.Combinator[rune, Position, func(string, string) string] {
	return MapStrings(map[string]func(string, string) string{
		"+": func(x, y string) string { return fmt.Sprintf("(%s+%s)", x, y) },
		"-": func(x, y string) string { return fmt.Sprintf("(%s-%s)", x, y) },
	})
}

func element() p.Combinator[rune, Position, string] {
    return Choice(
        Try(String("x")),
        Try(String("y")),
        Try(String("z")),
    )
}

func e() p.Combinator[rune, Position, string] {
	return Choice(
		Try(element()),
		Try(Parens(__expr())),
	)
}

func __expr() p.Combinator[rune, Position, string] {
	return Chainl1(e(), op())
}

func TestChain(t *testing.T) {
	parser := __expr()
	tests := []parsingTest{
		{`x+y+z-x`, `(((x+y)+z)-x)`, true},
		{`(x+y)+z-x`, `(((x+y)+z)-x)`, true},
	}
	strictParse(t, parser, tests)
}

func TestExpr(t *testing.T) {
	parser := Trace(t, "condition", _expr(t))

	tests := []parsingTest{
		{`labels[host.name] = "host01"`, condition.NewEqual(_f("labels", "host.name"), "host01"), true},
		{`labels[host.name] != "host01"`, condition.NewNotEqual(_f("labels", "host.name"), "host01"), true},
		{`severity =~ 'error'`, condition.NewMatch(_f("severity", ""), "error"), true},
		{`!labels[host.name] = "host01"`, condition.NewNot(condition.NewEqual(_f("labels", "host.name"), "host01")), true},
		{`! has severity_text`, condition.NewNot(condition.NewHas(_f("severity_text", ""))), true},
		{`! `, nil, false},
		{`severity_text = "error" and severity_number = "12"`, condition.NewAnd(
			condition.NewEqual(_f("severity_text", ""), "error"),
			condition.NewEqual(_f("severity_number", ""), "12"),
		), true},
		{`labels[host.name] = "host01" or labels[host.name] = "host02"`, condition.NewOr(
			condition.NewEqual(_f("labels", "host.name"), "host01"),
			condition.NewEqual(_f("labels", "host.name"), "host02"),
		), true},
		{`a = "1" and (b = "x" or b = "y")`, condition.NewAnd(
			condition.NewEqual(_f("a", ""), "1"),
			condition.NewOr(
				condition.NewEqual(_f("b", ""), "x"),
				condition.NewEqual(_f("b", ""), "y"),
			),
		), true},
		{``, nil, false},
	}
	strictParse(t, parser, tests)
}
