package parser

import (
	"fmt"
	. "github.com/prataprc/goparsec"
)

func as_str(n ParsecNode) string {
	return n.(*Terminal).Value
}
func not_none(n ParsecNode) bool {
	_, ok := n.(MaybeNone)
	return !ok
}
func f_one(ns []ParsecNode) ParsecNode {
	fmt.Printf("[f_one] %#v\n", ns)
	for _, n := range ns {
		if v, ok := n.(*Terminal); ok {
			fmt.Printf("        - %s\n", v.GetName())
		}
	}
	if len(ns) == 0 {
		return nil
	}
	return ns[0]
}
