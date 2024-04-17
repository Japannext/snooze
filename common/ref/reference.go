package ref

import (
  "fmt"

  api "github.com/japannext/snooze/common/api/v2"
)

type KvKind string

const (
  LABELS KvKind = "labels"
  ATTRIBUTES    = "attributes"
  GROUP         = "group_kv"
)

type Reference struct {
  Kv KvKind
  Key string
}

func (r *Reference) Get(a *api.Alert) (string, bool) {
  var (
    v string
    found bool
  )
  switch r.Kv {
    case LABELS:
      v, found = a.Labels[r.Key]
    case ATTRIBUTES:
      v, found = a.Attributes[r.Key]
    case GROUP:
      v, found = a.GroupKv[r.Key]
    default:
      v, found = "", false
  }
  return v, found
}

func (r *Reference) Set(a *api.Alert, v string) {
  switch r.Kv {
    case LABELS:
      a.Labels[r.Key] = v
    case ATTRIBUTES:
      a.Attributes[r.Key] = v
    case GROUP:
      a.GroupKv[r.Key] = v
  }
}

func (r *Reference) Validate() error {
  switch r.Kv {
    case LABELS, ATTRIBUTES, GROUP:
      return nil
    default:
      return fmt.Errorf("Unsupported map '%s'", r.Kv)
  }
}

func (r *Reference) String() string {
  return fmt.Sprintf("%s[%s]", r.Kv, r.Key)
}
