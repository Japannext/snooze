package lang

import (
	"context"
	"fmt"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

type Field struct {
	raw  string
	gval gval.Evaluable
}

func NewField(raw string) (*Field, error) {
	computed, err := jsonpath.New(raw)
	if err != nil {
		return &Field{}, err
	}
	return &Field{raw, computed}, nil
}

func NewFields(raws []string) ([]Field, error) {
	var fields []Field
	for _, raw := range raws {
		f, err := NewField(raw)
		if err != nil {
			return []Field{}, fmt.Errorf("with field `%s`: %w", raw, err)
		}
		fields = append(fields, *f)
	}
	return fields, nil
}

func (f *Field) String() string {
	return f.raw
}

func (f *Field) Get(ctx context.Context, item *api.Log) (string, error) {
	v, err := f.gval.EvalString(ctx, map[string]interface{}{"log": item})
	if err != nil {
		return "", err
	}
	return v, nil
}
