package transform

import (
	"context"
	"fmt"
	"regexp"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/models"
)

type RegexAction struct {
	Field string `yaml:"field" json:"field"`
	Match string `yaml:"match" json:"match"`

	internal struct {
		field  *lang.Field
		regexp *regexp.Regexp
	}
}

func (action *RegexAction) Load() Transformation {
	var err error
	action.internal.field, err = lang.NewField(action.Field)
	if err != nil {
		log.Fatalf("invalid field `%s`: %s", action.Field, err)
	}
	action.internal.regexp, err = regexp.Compile(action.Match)
	if err != nil {
		log.Fatalf("invalid regex `%s`: %s", action.Match, err)
	}
	return action
}

func (action *RegexAction) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	capture, ok := ctx.Value("capture").(map[string]string)
	if !ok {
		capture = make(map[string]string)
	}
	value, err := lang.ExtractField(item, action.internal.field)
	if err != nil {
		return ctx, fmt.Errorf("failed to extract info for field `%s`: %s", action.Field, err)
	}
	match := action.internal.regexp.MatchString(value)
	if !match {
		return ctx, nil
	}
	keys := action.internal.regexp.SubexpNames()
	if len(keys) > 1 {
		keys = keys[1:]
		values := action.internal.regexp.FindStringSubmatch(value)
		for _, key := range keys {
			i := action.internal.regexp.SubexpIndex(key)
			if i < 0 {
				continue
			}
			capture[key] = values[i]
		}
	}
	ctx = context.WithValue(ctx, "capture", capture)
	return ctx, nil
}
