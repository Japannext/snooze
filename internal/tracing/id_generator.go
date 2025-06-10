package tracing

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand/v2"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// A type used to store trace information in context
type ContextValue struct{}

type IDGenerator struct{}

func (idgen *IDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	// If something is in the context, attempt to extract the ID.
	traceID, ok := ctx.Value(ContextValue{}).(string)
	if ok {
		tid, err := stringTo16Bytes(traceID)
		if err != nil {
			log.Warnf("invalid traceID (%s): %s", traceID, err)
		} else {
			sid := idgen.NewSpanID(ctx, tid)
			return tid, sid
		}
	}

	// Defaults to random gen.
	return idgen.traceIDGen(ctx)
}

// Copy-pasta of otel trace random id generator
func (idgen *IDGenerator) traceIDGen(ctx context.Context) (trace.TraceID, trace.SpanID) {
	tid := trace.TraceID{}
	sid := trace.SpanID{}

	for {
		binary.NativeEndian.PutUint64(tid[:8], rand.Uint64())
		binary.NativeEndian.PutUint64(tid[8:], rand.Uint64())
		if tid.IsValid() {
			break
		}
	}
	for {
		binary.NativeEndian.PutUint64(sid[:], rand.Uint64())
		if sid.IsValid() {
			break
		}
	}
	return tid, sid
}

// Copy-pasta of otel trace random id generator.
func (idgen *IDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	sid := trace.SpanID{}
	for {
		binary.NativeEndian.PutUint64(sid[:], rand.Uint64())
		if sid.IsValid() {
			break
		}
	}
	return sid
}

func stringTo16Bytes(s string) ([16]byte, error) {
	if len(s) != 16 {
		return [16]byte{}, fmt.Errorf("the provided string is not 16 bytes")
	}
	var res [16]byte
	copy(res[:], []byte(s)[:16])
	return res, nil
}
