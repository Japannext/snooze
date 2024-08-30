package silence

import (
	"context"
	"fmt"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func Process(item *api.Log) error {

	ctx := context.Background()

	for _, s := range silences {
		v, err := s.internal.condition.MatchLog(ctx, item)
		if err != nil {
			return err
		}
		if v {
			item.Mute.Silence(fmt.Sprintf("Silenced by rule %s", s))
			break
		}
	}

	return nil
}
