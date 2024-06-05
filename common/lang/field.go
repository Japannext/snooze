package lang

import (
	"context"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"

	api "github.com/japannext/snooze/common/api/v2"
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

func (f *Field) String() string {
	return f.raw
}

func (f *Field) Get(ctx context.Context, alert *api.Alert) (string, error) {
	v, err := f.gval.EvalString(ctx, map[string]interface{}{"alert": alert})
	if err != nil {
		return "", err
	}
	return v, nil
}
