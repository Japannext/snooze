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

func NewFields(raws []string) ([]*Field, error) {
	var fields = []*Field{}
	for _, raw := range raws {
		f, err := NewField(raw)
		if err != nil {
			return fields, fmt.Errorf("with field `%s`: %w", raw, err)
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func (f *Field) String() string {
	return f.raw
}

func ExtractField(item api.HasContext, field *Field) (string, error) {
	ctx := context.Background()
	return field.gval.EvalString(ctx, item.Context())
}

func ExtractFields(item api.HasContext, fields []*Field) (map[string]string, error) {
	mapper := map[string]string{}
	for _, field := range fields {
		value, err := ExtractField(item, field)
		if err != nil {
			return mapper, err
		}
		mapper[field.String()] = value
	}
	return mapper, nil
}
