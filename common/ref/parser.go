package ref

import (
  p "github.com/vektah/goparsify"
)

var parser p.Parser

func Parse(data string) (*Reference, error) {
  r, err := p.Run(parser, data)
  if err != nil {
    return nil, err
  }
  rr := r.(Reference)
  return &rr, err
}

func init() {
  kv := p.Any(
    p.Exact("labels"),
    p.Exact("attributes"),
    p.Exact("group"),
  ).Map(func (r *p.Result) { r.Result = r.Token })
  key := p.Chars("a-zA-Z0-9.,/_-", 1, 255).Map(func (r *p.Result) { r.Result = r.Token })
  parser = p.Seq(kv, "[", p.Cut(), key, "]").Map(func(res *p.Result) {
    res.Result = Reference{
      Kv: KvKind(res.Child[0].Token),
      Key: res.Child[3].Token,
    }
  })
}


