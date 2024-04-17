package v1

import (
  "fmt"
  "strings"
  "regexp"

  api "github.com/japannext/snooze/common/api/v2"
  "github.com/japannext/snooze/common/ref"
)

type ConditionKind string

const (
  AND ConditionKind = "and"
  OR                = "or"
  NOT               = "not"
  EQUAL             = "equal"
  NOT_EQUAL         = "not_equal"
  MATCH             = "match"
  NOT_MATCH         = "not_match"
)

type Condition struct {
  Kind ConditionKind `json:"kind"`
  And *AndCondition `json:"and,omitempty"`
  Or *OrCondition `json:"or,omitempty"`
  Not *NotCondition `json:"not,omitempty"`
  Equal *EqualCondition `json:"equal,omitempty"`
  NotEqual *NotEqualCondition `json:"notEqual,omitempty"`
  Match *MatchCondition `json:"match,omitempty"`
  NotMatch *NotMatchCondition `json:"not_match,omitempty"`
}

func (c *Condition) Get() ConditionInterface {
  switch c.Kind {
    case AND:
      return c.And
    case OR:
      return c.Or
    case NOT:
      return c.Not
    case EQUAL:
      return c.Equal
    case NOT_EQUAL:
      return c.NotEqual
  }
  // Should never be reached
  return nil
}

func (c *Condition) Validate() error {
  switch c.Kind {
    case AND, OR, EQUAL, NOT, NOT_EQUAL:
      return c.Get().Validate()
    default:
      return fmt.Errorf("unsupported condition kind '%s'", c.Kind)
  }
}

type ConditionInterface interface {
  Test(*api.Alert) bool
  Validate() error
  String() string
}

type AndCondition struct {
  Conditions []*Condition `json:"conditions"`
}

func (c *AndCondition) Validate() error {
  for _, cc := range c.Conditions {
    err := cc.Get().Validate()
    if err != nil {
      return err
    }
  }
  return nil
}

func (c *AndCondition) Test(a *api.Alert) bool {
  for _, cc := range c.Conditions {
    if !cc.Get().Test(a) {
      return false
    }
  }
  return true
}

func (c *AndCondition) String() string {
  var b strings.Builder
  for i, cc := range c.Conditions {
    if i != 0 {
      b.WriteString(" and ")
    }
    b.WriteString("(")
    b.WriteString(cc.Get().String())
    b.WriteString(")")
  }
  b.WriteString("")
  return b.String()
}

type OrCondition struct {
  Conditions []*Condition `json:"conditions"`
}

func (c *OrCondition) Validate() error {
  for _, cc := range c.Conditions {
    err := cc.Get().Validate()
    if err != nil {
      return err
    }
  }
  return nil
}

func (c *OrCondition) Test(a *api.Alert) bool {
  for _, cc := range c.Conditions {
    if cc.Get().Test(a) {
      return true
    }
  }
  return false
}

func (c *OrCondition) String() string {
  var b strings.Builder
  for i, cc := range c.Conditions {
    if i != 0 {
      b.WriteString(" or ")
    }
    b.WriteString("(")
    b.WriteString(cc.Get().String())
    b.WriteString(")")
  }
  b.WriteString("")
  return b.String()
}

type NotCondition struct {
  Condition *Condition `json:"condition"`
}

func (c *NotCondition) Validate() error {
  return c.Condition.Get().Validate()
}

func (c *NotCondition) Test(a *api.Alert) bool {
  return c.Condition.Get().Test(a)
}

func (c *NotCondition) String() string {
  return fmt.Sprintf("!(%s)", c.Condition)
}

type EqualCondition struct {
  Ref ref.Reference `json:"reference"`
  Value string `json:"value"`
}

func (c *EqualCondition) Validate() error {
  return c.Ref.Validate()
}

func (c *EqualCondition) Test(a *api.Alert) bool {
  v, found := c.Ref.Fetch(a)
  if found && v == c.Value {
    return true
  }
  return false
}

func (c *EqualCondition) String() string {
  return fmt.Sprintf("%s == '%s'", c.Ref.String(), c.Value)
}

type NotEqualCondition struct {
  Ref ref.Reference `json:"reference"`
  Value string `json:"value"`
}

func (c *NotEqualCondition) Validate() error {
  return c.Ref.Validate()
}

func (c *NotEqualCondition) Test(a *api.Alert) bool {
  v, found := c.Ref.Fetch(a)
  if found && v == c.Value {
    return false
  }
  return true
}

func (c *NotEqualCondition) String() string {
  return fmt.Sprintf("%s != '%s'", c.Ref.String(), c.Value)
}

type MatchCondition struct {
  Ref ref.Reference
  Value string
  regexp *regexp.Regexp
}

func (c *MatchCondition) Validate() error {
  if err := c.Ref.Validate(); err != nil {
    return err
  }
  re, err := regexp.Compile(c.Value)
  if err != nil {
    return err
  }
  c.regexp = re
  return nil
}

func (c *MatchCondition) Test(a *api.Alert) bool {
  v, found := c.Ref.Fetch(a)
  if found && c.regexp.Match([]byte(v)) {
    return true
  }
  return false
}

func (c *MatchCondition) String() string {
  return fmt.Sprintf("%s =~ /%s/", c.Ref.String(), c.Value)
}

type NotMatchCondition struct {
  Ref ref.Reference
  Value string
  regexp *regexp.Regexp
}

func (c *NotMatchCondition) Validate() error {
  if err := c.Ref.Validate(); err != nil {
    return err
  }
  re, err := regexp.Compile(c.Value)
  if err != nil {
    return err
  }
  c.regexp = re
  return nil
}

func (c *NotMatchCondition) Test(a *api.Alert) bool {
  v, found := c.Ref.Fetch(a)
  if found && c.regexp.Match([]byte(v)) {
    return false
  }
  return true
}

func (c *NotMatchCondition) String() string {
  return fmt.Sprintf("%s !~ /%s/", c.Ref.String(), c.Value)
}
