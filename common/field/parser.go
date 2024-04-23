package field

import (
  p "github.com/vektah/goparsify"
)

var parser p.Parser

var fname = p.Chars("a-z_", 1, 255)
var fkey = p.Chars("a-zA-Z0-9.,/_-", 1, 255)
var FieldParser = p.Seq(fname, "[", p.Cut(), fkey, "]").Map(func(r *p.Result) {
  r.Result = AlertField{
    r.Child[0].Token,
    r.Child[3].Token,
  }
})

func Parse(data string) (*AlertField, error) {
  f, err := p.Run(FieldParser, data)
  if err != nil {
    return nil, err
  }
  fi := f.(AlertField)
  return &fi, nil
}
