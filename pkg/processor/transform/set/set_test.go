package set_test

import (
	"context"
	"testing"

	"github.com/japannext/snooze/pkg/processor/transform/set"
	"github.com/japannext/snooze/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	t.Parallel()

	tests := []struct{
		Name string
		Config *set.Config
		Log *models.Log
		ExpectedContext map[string]string
		ExpectedLabels  map[string]string
		ExpectedIdentity map[string]string
	}{
		{
			Name: "set 1 label",
			Config: &set.Config{Labels: map[string]string{"a": "x"}},
			Log: &models.Log{Labels: map[string]string{"v": "1"}},
			ExpectedLabels: map[string]string{"v": "1", "a": "x"},
		},
		{
			Name: "override 1 label",
			Config: &set.Config{Labels: map[string]string{"a": "x", "b": "y"}},
			Log: &models.Log{Labels: map[string]string{"a": "1"}},
			ExpectedLabels: map[string]string{"a": "x", "b": "y"},
		},
		{
			Name: "append 1 identity",
			Config: &set.Config{Identity: map[string]string{"process": "myapp"}},
			Log: &models.Log{Identity: map[string]string{"host": "myapp01"}},
			ExpectedIdentity: map[string]string{"host": "myapp01", "process": "myapp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			action, err := set.New(tt.Config)
			require.NoError(t, err)

			ctx := context.TODO()

			ctx, err = action.Process(ctx, tt.Log)

			if tt.ExpectedContext != nil {
				for key, value := range tt.ExpectedContext {
					assert.Equal(t, value, ctx.Value(key))
				}
			}

			if tt.ExpectedLabels != nil {
				assert.Equal(t, tt.ExpectedLabels, tt.Log.Labels)
			}

			if tt.ExpectedIdentity != nil {
				assert.Equal(t, tt.ExpectedIdentity, tt.Log.Identity)
			}
		})
	}

}
