package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
	"github.com/japannext/snooze/common/condition"
	"github.com/japannext/snooze/common/field"
)

var _f = field.NewAlertField

type parsingTest struct {
	Text string
	Expected interface{}
	Success bool
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
	}
	strictParse(t, parser, tests)
}

func TestNot(t *testing.T) {
	e := Trace(t, "operation", operation())
	parser := Trace(t, "not", not(e))
	tests := []parsingTest{
		{`! labels[host.name] = "host01"`, condition.NewNot(condition.NewEqual(_f("labels", "host.name"), "host01")), true},
		{`! has severity_text`, condition.NewNot(condition.NewHas(_f("severity_text", ""))), true},
		{`! `, nil, false},
	}
	strictParse(t, parser, tests)
}

func TestAndOr(t *testing.T) {
	e := Trace(t, "operation", operation())
	parser := Trace(t, "andor", andor(e))
	tests := []parsingTest{
		{`severity_text = "error" and severity_number = "12"`, condition.NewAnd(
			condition.NewEqual(_f("severity_text", ""), "error"),
			condition.NewEqual(_f("severity_number", ""), "12"),
		), true},
		{`labels[host.name] = "host01" or labels[host.name] = "host02"`, condition.NewOr(
			condition.NewEqual(_f("labels", "host.name"), "host01"),
			condition.NewEqual(_f("labels", "host.name"), "host02"),
		), true},
	}
	strictParse(t, parser, tests)
}

