package ref

import (
  "testing"
  "github.com/stretchr/testify/assert"

  api "github.com/japannext/snooze/common/api/v2"
)

func TestReferenceGet(t *testing.T) {
  a := &api.Alert{
    Labels: map[string]string{
      "k8s.namespace": "ns1",
      "k8s.deployment.name": "myapp",
    },
  }
  ref := Reference{Kv: KvKind("labels"), Key: "k8s.namespace"}
  v, ok := ref.Get(a)
  assert.Equal(t, true, ok)
  assert.Equal(t, "ns1", v)
}

func TestReferenceSetOverride(t *testing.T) {
  a := &api.Alert{
    Labels: map[string]string{
      "k8s.namespace": "ns1",
      "k8s.deployment.name": "myapp",
    },
  }
  ref := Reference{Kv: KvKind("labels"), Key: "k8s.namespace"}
  ref.Set(a, "ns2")
  m := map[string]string{
    "k8s.namespace": "ns2",
    "k8s.deployment.name": "myapp",
  }
  assert.Equal(t, m, a.Labels)
}

func TestReferenceSetNew(t *testing.T) {
  a := &api.Alert{
    Labels: map[string]string{
      "k8s.namespace": "ns1",
      "k8s.deployment.name": "myapp",
    },
  }
  ref := Reference{Kv: KvKind("labels"), Key: "k8s.cluster.name"}
  ref.Set(a, "dev")
  m := map[string]string{
      "k8s.namespace": "ns1",
      "k8s.deployment.name": "myapp",
      "k8s.cluster.name": "dev",
  }
  assert.Equal(t, m, a.Labels)
}
