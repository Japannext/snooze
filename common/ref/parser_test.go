package ref

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestParseLabelReference(t *testing.T) {
  r, err := Parse("labels[host.name]")
  if assert.NoError(t, err) {
    assert.Equal(t, &Reference{Kv: KvKind("labels"), Key: "host.name"}, r)
  }
}

func TestParserAttributeReference(t *testing.T) {
  r, err := Parse("attributes[host.arch]")
  if assert.NoError(t, err) {
    assert.Equal(t, &Reference{Kv: KvKind("attributes"), Key: "host.arch"}, r)
  }
}

func TestParserGroupReference(t *testing.T) {
  r, err := Parse("group[k8s.namespace]")
  if assert.NoError(t, err) {
    assert.Equal(t, &Reference{Kv: KvKind("group"), Key: "k8s.namespace"}, r)
  }
}

var notParsable = []string{
  "labels[blah",
  "attribtues[valid.name]",
  "labels[a][b]",
}

func TestReferenceNotParsable(t *testing.T) {
  var err error
  for _, data := range notParsable {
    _, err = Parse(data)
    assert.Error(t, err)
  }
}
