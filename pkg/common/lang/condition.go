package lang

import (
	"context"
	"fmt"
	"strconv"
	"reflect"
	"github.com/PaesslerAG/gval"

	"github.com/japannext/snooze/pkg/models"
)

type Condition struct {
	raw  string
	gval gval.Evaluable
}

func variable(path gval.Evaluables) gval.Evaluable {
	return func(c context.Context, v interface{}) (interface{}, error) {
		keys, err := path.EvalStrings(c, v)
		if err != nil {
			return nil, err
		}
		for _, k := range keys {
			switch o := v.(type) {
			case gval.Selector:
				v, err = o.SelectGVal(c, k)
				if err != nil {
					return nil, fmt.Errorf("failed to select '%s' on %T: %w", k, o, err)
				}
				continue
			case map[interface{}]interface{}:
				v = o[k]
				continue
			case map[string]interface{}:
				v = o[k]
				continue
			case []interface{}:
				if i, err := strconv.Atoi(k); err == nil && i >= 0 && len(o) > i {
					v = o[i]
					continue
				}
			default:
				var ok bool
				v, ok = reflectSelect(k, o)
				if !ok {
					return nil, nil
					// return nil, fmt.Errorf("unknown parameter %s", strings.Join(keys[:i+1], "."))
				}
			}
		}
		return v, nil
	}
}

func reflectSelect(key string, value interface{}) (selection interface{}, ok bool) {
	vv := reflect.ValueOf(value)
	vvElem := resolvePotentialPointer(vv)

	switch vvElem.Kind() {
	case reflect.Map:
		mapKey, ok := reflectConvertTo(vv.Type().Key().Kind(), key)
		if !ok {
			return nil, false
		}

		vvElem = vv.MapIndex(reflect.ValueOf(mapKey))
		vvElem = resolvePotentialPointer(vvElem)

		if vvElem.IsValid() {
			return vvElem.Interface(), true
		}

		// key didn't exist. Check if there is a bound method
		method := vv.MethodByName(key)
		if method.IsValid() {
			return method.Interface(), true
		}

	case reflect.Slice:
		if i, err := strconv.Atoi(key); err == nil && i >= 0 && vv.Len() > i {
			vvElem = resolvePotentialPointer(vv.Index(i))
			return vvElem.Interface(), true
		}

		// key not an int. Check if there is a bound method
		method := vv.MethodByName(key)
		if method.IsValid() {
			return method.Interface(), true
		}

	case reflect.Struct:
		field := vvElem.FieldByName(key)
		if field.IsValid() {
			return field.Interface(), true
		}

		method := vv.MethodByName(key)
		if method.IsValid() {
			return method.Interface(), true
		}
	}
	return nil, false
}

func resolvePotentialPointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func reflectConvertTo(k reflect.Kind, value string) (interface{}, bool) {
	switch k {
	case reflect.String:
		return value, true
	case reflect.Int:
		if i, err := strconv.Atoi(value); err == nil {
			return i, true
		}
	}
	return nil, false
}

// Check if a value is not nil/empty
func hasFunction(args ...interface{}) (interface{}, error) {
	for _, arg := range args {
		v := reflect.ValueOf(arg)
		if !v.IsValid() {
			return false, nil
		}
		if v.IsZero() {
			return false, nil
		}
	}
	return true, nil
}

func NewCondition(raw string) (*Condition, error) {
	ll := gval.Full(
		gval.VariableSelector(variable),
		gval.Function("has", hasFunction),
	)
	e, err := ll.NewEvaluable(raw)
	if err != nil {
		return &Condition{}, err
	}
	return &Condition{raw, e}, nil
}

func (c *Condition) String() string {
	return c.raw
}

func (c *Condition) MatchLog(ctx context.Context, item *models.Log) (bool, error) {
	return c.gval.EvalBool(ctx, item.Context())
}
