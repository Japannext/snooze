package parser

import (
	"fmt"

	. "github.com/prataprc/goparsec"

	"github.com/japannext/snooze/common/field"
)

var fieldParser Parser

func initFieldParser() {
	var (
		sbrk         = AtomExact("[", "START_BRACKET")
		ebrk         = AtomExact("]", "END_BRACKET")
		label        = Token("[a-zA-Z0-9._/-]+", "LABEL")
		strict_label = TokenExact("[a-zA-Z0-9._/-]+", "STRICT_LABEL")
		subkey       = And(f_subkey, sbrk, strict_label, ebrk)
	)

	fieldParser = And(f_field, label, Maybe(f_one, subkey))
}

func f_field(ns []ParsecNode) ParsecNode {
	fmt.Printf("[f_field] %#v\n", ns)
	name := as_str(ns[0])
	var subkey string
	if not_none(ns[1]) {
		subkey, _ = ns[1].(string)
	}
	return &field.AlertField{Name: name, SubKey: subkey}
}

func f_subkey(ns []ParsecNode) ParsecNode {
	fmt.Printf("[f_subkey] %#v\n", ns)
	return ns[1].(*Terminal).GetValue()
}

func ParseField(text string) (*field.AlertField, error) {
	s := NewScanner([]byte(text))
	res, next := fieldParser(s)
	if !next.Endof() {
		return nil, fmt.Errorf("failed to parse '%s' char number %d", text, next.GetCursor())
	}
	f, ok := res.(*field.AlertField)
	if !ok {
		return nil, fmt.Errorf("no field matched (`%s`)", text)
	}
	fmt.Printf("field: %#v\n", f)
	return f, nil
}
