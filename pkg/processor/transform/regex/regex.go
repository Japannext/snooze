package regex

import (
	"context"
	"fmt"
	"regexp"

	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/models"
)

type Config struct {
	Field string `json:"field" yaml:"field"`
	Match string `json:"match" yaml:"match"`
}

type Action struct {
	cfg *Config

	field  *lang.Field
	regexp *regexp.Regexp
}

func New(cfg *Config) (*Action, error) {
	var err error

	action := &Action{cfg: cfg}

	action.field, err = lang.NewField(action.cfg.Field)
	if err != nil {
		return action, fmt.Errorf("invalid field `%s`: %s", action.cfg.Field, err)
	}

	action.regexp, err = regexp.Compile(action.cfg.Match)
	if err != nil {
		return action, fmt.Errorf("invalid regex `%s`: %s", action.cfg.Match, err)
	}

	return action, nil
}

func (action *Action) Process(ctx context.Context, item *models.Log) (context.Context, error) {
	capture, ok := ctx.Value("capture").(map[string]string)
	if !ok {
		capture = make(map[string]string)
	}

	value, err := lang.ExtractField(item, action.field)
	if err != nil {
		return ctx, fmt.Errorf("failed to extract info for field `%s`: %s", action.cfg.Field, err)
	}

	match := action.regexp.MatchString(value)
	if !match {
		return ctx, nil
	}

	keys := action.regexp.SubexpNames()

	if len(keys) > 1 {
		keys = keys[1:]
		values := action.regexp.FindStringSubmatch(value)
		for _, key := range keys {
			i := action.regexp.SubexpIndex(key)
			if i < 0 {
				continue
			}
			capture[key] = values[i]
		}
	}

	ctx = context.WithValue(ctx, "capture", capture)

	return ctx, nil
}
