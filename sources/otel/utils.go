package otel

import (
	"strconv"
	"strings"
	"time"

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

type AnyValue struct {
	*commonv1.AnyValue
}

// Utility to return a single key-value
func m(key, value string) map[string]string {
	return map[string]string{key: value}
}

func (x AnyValue) ToMap() map[string]string {

	switch x.GetValue().(type) {
	case *commonv1.AnyValue_StringValue:
		return m("string", x.GetStringValue())
	case *commonv1.AnyValue_BytesValue:
		return m("bytes", string(x.GetBytesValue()))
	case *commonv1.AnyValue_BoolValue:
		return m("bool", strconv.FormatBool(x.GetBoolValue()))
	case *commonv1.AnyValue_IntValue:
		return m("int", strconv.Itoa(int(x.GetIntValue())))
	case *commonv1.AnyValue_DoubleValue:
		return m("float", strconv.FormatFloat(x.GetDoubleValue(), 'f', -1, 64))
	case *commonv1.AnyValue_ArrayValue:
		var b strings.Builder
		for i, y := range x.GetArrayValue().GetValues() {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(AnyValue{y}.ToString())
		}
		return m("array", b.String())
	case *commonv1.AnyValue_KvlistValue:
		var b strings.Builder
		for i, kv := range x.GetKvlistValue().GetValues() {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(kv.Key)
			b.WriteString("=")
			b.WriteString(AnyValue{kv.Value}.ToString())
		}
		return m("map", b.String())
	}
	return m("none", "")
}

func (x AnyValue) ToString() string {
	switch x.GetValue().(type) {
	case *commonv1.AnyValue_StringValue:
		return x.GetStringValue()
	case *commonv1.AnyValue_BytesValue:
		return string(x.GetBytesValue())
	case *commonv1.AnyValue_BoolValue:
		return strconv.FormatBool(x.GetBoolValue())
	case *commonv1.AnyValue_IntValue:
		return strconv.Itoa(int(x.GetIntValue()))
	case *commonv1.AnyValue_DoubleValue:
		return strconv.FormatFloat(x.GetDoubleValue(), 'f', -1, 64)
	case *commonv1.AnyValue_ArrayValue:
		var b strings.Builder
		for i, y := range x.GetArrayValue().GetValues() {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(AnyValue{y}.ToString())
		}
		return b.String()
	case *commonv1.AnyValue_KvlistValue:
		var b strings.Builder
		for i, kv := range x.GetKvlistValue().GetValues() {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(kv.Key)
			b.WriteString("=")
			b.WriteString(AnyValue{kv.Value}.ToString())
		}
		return b.String()
	}
	return ""
}

// Return a time in the correct format
func timeNow() uint64 {
	return uint64(time.Now().UnixNano())
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
