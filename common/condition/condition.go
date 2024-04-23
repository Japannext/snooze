package condition

import (
	"fmt"
	"regexp"
	"strings"

	api "github.com/japannext/snooze/common/api/v2"
	"github.com/japannext/snooze/common/field"
)

type ConditionKind string

const (
	AND       ConditionKind = "and"
	OR                      = "or"
	NOT                     = "not"
	EQUAL                   = "equal"
	NOT_EQUAL               = "not_equal"
	MATCH                   = "match"
	NOT_MATCH               = "not_match"
	HAS                     = "has"
)

type Condition struct {
	Kind     ConditionKind      `json:"kind"`
	And      *AndCondition      `json:"and,omitempty"`
	Or       *OrCondition       `json:"or,omitempty"`
	Not      *NotCondition      `json:"not,omitempty"`
	Equal    *EqualCondition    `json:"equal,omitempty"`
	NotEqual *NotEqualCondition `json:"notEqual,omitempty"`
	Match    *MatchCondition    `json:"match,omitempty"`
	NotMatch *NotMatchCondition `json:"not_match,omitempty"`
	Has      *HasCondition      `json:"has,omitempty"`
}

func (c *Condition) Resolve() Interface {
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
	case HAS:
		return c.Has
	}
	// Should never be reached
	return nil
}

func (c *Condition) Validate() error {
	switch c.Kind {
	case AND, OR, EQUAL, NOT, NOT_EQUAL:
		return c.Resolve().Validate()
	default:
		return fmt.Errorf("unsupported condition kind '%s'", c.Kind)
	}
}

type Interface interface {
	Test(*api.Alert) bool
	Validate() error
	String() string
}

type AndCondition struct {
	Conditions []*Condition `json:"conditions"`
}

func (c *AndCondition) Validate() error {
	for _, cc := range c.Conditions {
		err := cc.Resolve().Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *AndCondition) Test(a *api.Alert) bool {
	for _, cc := range c.Conditions {
		if !cc.Resolve().Test(a) {
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
		b.WriteString(cc.Resolve().String())
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
		err := cc.Resolve().Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *OrCondition) Test(a *api.Alert) bool {
	for _, cc := range c.Conditions {
		if cc.Resolve().Test(a) {
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
		b.WriteString(cc.Resolve().String())
		b.WriteString(")")
	}
	b.WriteString("")
	return b.String()
}

type NotCondition struct {
	Condition *Condition `json:"condition"`
}

func (c *NotCondition) Validate() error {
	return c.Condition.Resolve().Validate()
}

func (c *NotCondition) Test(a *api.Alert) bool {
	return c.Condition.Resolve().Test(a)
}

func (c *NotCondition) String() string {
	return fmt.Sprintf("!(%s)", c.Condition)
}

type EqualCondition struct {
	Field field.AlertField `json:"field"`
	Value string           `json:"value"`
}

func (c *EqualCondition) Validate() error {
	return c.Field.Validate()
}

func (c *EqualCondition) Test(a *api.Alert) bool {
	v, found := c.Field.Get(a)
	if found && v == c.Value {
		return true
	}
	return false
}

func (c *EqualCondition) String() string {
	return fmt.Sprintf("%s == '%s'", c.Field, c.Value)
}

type NotEqualCondition struct {
	Field field.AlertField `json:"field"`
	Value string           `json:"value"`
}

func (c *NotEqualCondition) Validate() error {
	return c.Field.Validate()
}

func (c *NotEqualCondition) Test(a *api.Alert) bool {
	v, found := c.Field.Get(a)
	if found && v == c.Value {
		return false
	}
	return true
}

func (c *NotEqualCondition) String() string {
	return fmt.Sprintf("%s != '%s'", c.Field.String(), c.Value)
}

type MatchCondition struct {
	Field  field.AlertField `json:"field"`
	Value  string           `json:"value"`
	regexp *regexp.Regexp
}

func (c *MatchCondition) Validate() error {
	if err := c.Field.Validate(); err != nil {
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
	v, found := c.Field.Get(a)
	if found && c.regexp.Match([]byte(v)) {
		return true
	}
	return false
}

func (c *MatchCondition) String() string {
	return fmt.Sprintf("%s =~ /%s/", c.Field, c.Value)
}

type NotMatchCondition struct {
	Field  field.AlertField `json:"field"`
	Value  string           `json:"value"`
	regexp *regexp.Regexp
}

func (c *NotMatchCondition) Validate() error {
	if err := c.Field.Validate(); err != nil {
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
	v, found := c.Field.Get(a)
	if found && c.regexp.Match([]byte(v)) {
		return false
	}
	return true
}

func (c *NotMatchCondition) String() string {
	return fmt.Sprintf("%s !~ /%s/", c.Field, c.Value)
}

type HasCondition struct {
	Field field.AlertField `json:"field"`
}

func (c *HasCondition) Validate() error {
	if err := c.Field.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *HasCondition) Test(a *api.Alert) bool {
	_, found := c.Field.Get(a)
	return found
}

func (c *HasCondition) String() string {
	return fmt.Sprintf("has %s", c.Field)
}
