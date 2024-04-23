package condition

import (
  p "github.com/vektah/goparsify"

  "github.com/japannext/snooze/common/field"
)

var parser p.Parser

var str = p.Any(p.StringLit(`"`), p.StringLit(`'`)).Map(func (r *p.Result){ r.Result = r.Token })

var equal = p.Seq(field.FieldParser, p.Exact("="), str).Map(func (r *p.Result) {
  r.Result = &Condition{Kind: "equal", Equal: &EqualCondition{
    Field: r.Child[0].Result.(field.AlertField),
    Value: r.Child[2].Result.(string),
  }}
})
var notEqual = p.Seq(field.FieldParser, p.Exact("!="), str).Map(func (r *p.Result) {
  r.Result = &Condition{Kind: "not_equal", NotEqual: &NotEqualCondition{
      Field: r.Child[0].Result.(field.AlertField),
      Value: r.Child[2].Result.(string),
    }}
})

var match = p.Seq(field.FieldParser, p.Exact("=~"), str).Map(func (r *p.Result) {
  r.Result = &Condition{Kind: "match", Match: &MatchCondition{
      Field: r.Child[0].Result.(field.AlertField),
      Value: r.Child[2].Result.(string),
    }}
})

var notMatch = p.Seq(field.FieldParser, p.Exact("!="), str).Map(func (r *p.Result) {
  r.Result = &Condition{Kind: "not_match", NotMatch: &NotMatchCondition{
    Field: r.Child[0].Result.(field.AlertField),
    Value: r.Child[2].Result.(string),
  }}
})

var has = p.Seq("has", field.FieldParser).Map(func (r *p.Result) {
  r.Result = HasCondition{r.Child[1].Result.(field.AlertField)}
})

var rawOp = p.Any(has, equal, notEqual, match, notMatch)
var wrappedOp = p.Seq("(", rawOp, ")").Map(func (r *p.Result) { r.Result = r.Child[1].Result })
var op = p.Any(wrappedOp, rawOp)

var and = p.Seq(op, "and", op).Map(func (r *p.Result) {
  r.Result = Condition{Kind: "and", And: &AndCondition{
    Conditions: []*Condition{
      r.Child[0].Result.(*Condition),
      r.Child[2].Result.(*Condition),
    },
  }}
})
var or = p.Seq(op, "or", op).Map(func (r *p.Result) {
  r.Result = Condition{Kind: "or", Or: &OrCondition{
    Conditions: []*Condition{
      r.Child[0].Result.(*Condition),
      r.Child[2].Result.(*Condition),
    },
  }}
})

var not = p.Seq("!", op)

func Parse(data string) (*Condition, error) {
  res, err := p.Run(parser, data)
  if err != nil {
    return nil, err
  }
  c := res.(Condition)
  return &c, nil
}
