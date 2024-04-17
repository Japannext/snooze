package otel

import (
  "time"
  "strconv"

  commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
)

// Check if a map has at least one key prefixed in a certain way
func hasPrefixedKey(m map[string]string, prefix string) bool {
  for k, _ := range m {
    if strings.HasPrefix(k, prefix) {
      return true
    }
  }
  return false
}

// Convert the opentelemetry keyvalue to a map[string]string
// TODO: handle types other than string with a conversion
func kvToMap(kvs []*commonv1.KeyValue) map[string]string {
  var mapping map[string]string = make(map[string]string)

  for _, kv := range kvs {
    if kv != nil {
      mapping[kv.Key] = kv.Value.GetStringValue()
    }
  }
  return mapping
}

// Return a time in the correct format
func timeNow() uint64 {
  return time.Now().Nanoseconds()
}

// Convert a string to the AnyValue type used
// in opentelemetry
func AnyString(s string) *commonv1.AnyValue {
  return &commonv1.AnyValue{
    Value: &commonv1.AnyValue_StringValue{
      StringValue: s,
    },
  }
}

func anyToString(x *commonv1.AnyValue) string {
  if y, ok := x.GetValue().(*commonv1.AnyValue_StringValue); ok {
    return y
  }
  if y, ok := x.GetValue().(*commonv1.AnyValue_BoolValue); ok {
    return strconv.FormatBool(y)
  }
  if y, ok := x.GetValue().(*commonv1.AnyValue_IntValue); ok {
    return strconv.Itoa(y)
  }
  if y, ok := x.GetValue().(*commonv1.AnyValue_DoubleValue); ok {
    return strconv.FormatFloat(y, 'f', -1, 64)
  }
  if y, ok := x.GetValue().(*commonv1.AnyValue_ArrayValue); ok {
    a := []string{}
    for i, z := range y {
      a = append(a, anyToString(z))
    }
    return a
  }
}
