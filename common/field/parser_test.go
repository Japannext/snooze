package field

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLabelAlertField(t *testing.T) {
	r, err := Parse("labels[host.name]")
	if assert.NoError(t, err) {
		assert.Equal(t, &AlertField{"labels", "host.name"}, r)
	}
}

func TestParserAttributeAlertField(t *testing.T) {
	r, err := Parse("attributes[host.arch]")
	if assert.NoError(t, err) {
		assert.Equal(t, &AlertField{"attributes", "host.arch"}, r)
	}
}

func TestParserGroupAlertField(t *testing.T) {
	r, err := Parse("group[k8s.namespace]")
	if assert.NoError(t, err) {
		assert.Equal(t, &AlertField{"group", "k8s.namespace"}, r)
	}
}

var notParsable = []string{
	"labels[blah",
	"attribtues[valid.name]",
	"labels[a][b]",
}

func TestAlertFieldNotParsable(t *testing.T) {
	var err error
	for _, data := range notParsable {
		_, err = Parse(data)
		assert.Error(t, err)
	}
}
