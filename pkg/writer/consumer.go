package writer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"go.opentelemetry.io/otel/trace"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/mq"
)

type Consumer struct {}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Run() error {
	for {
		log.Debugf("Ready to fetch items...")
		msgs, err := storeQ.Fetch(config.BatchSize, jetstream.FetchMaxWait(5*time.Second))
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		log.Debugf("Fetched %d items", len(msgs))
		if len(msgs) == 0 {
			continue
		}
		ctx := context.TODO()
		bulkWrite(ctx, msgs)
	}
}

func extractBulkError(resp opensearchapi.BulkRespItem) error {
	if resp.Error == nil {
		return nil
	}
	err := resp.Error
	var builder strings.Builder
	fmt.Fprintf(&builder, "[%s] %s", err.Type, err.Reason)
	if err.Cause.Type != "" {
		cause := err.Cause
		fmt.Fprintf(&builder, ": [%s] %s", cause.Type, cause.Reason)
		if cause.Cause.Reason != nil {
			cause2 := cause.Cause
			fmt.Fprintf(&builder, ": [%s] %s", cause2.Type, *cause2.Reason)
		}
	}
	return fmt.Errorf("%s", builder.String())
}

func bulkHeader(action, index string) []byte {
	m := map[string]map[string]string{
		action: map[string]string{
			"_index": index,
		},
	}
	b, _ := json.Marshal(m)
	return b
}

func bulkWrite(ctx context.Context, msgs []mq.MsgWithContext) {
	ctx, bulkSpan := tracer.Start(ctx, "bulkWrite")
	defer bulkSpan.End()
	var buf bytes.Buffer
	var inserting = map[int]jetstream.Msg{}
	for i, m := range msgs {
		msg, msgCtx := m.Extract()

		// Tracing
		msgCtx, span := tracer.Start(msgCtx, "bulkWrite", trace.WithLinks(trace.LinkFromContext(ctx)))
		defer span.End()

		// Extract index
		index := msg.Headers().Get(mq.X_SNOOZE_STORE_INDEX)
		if index == "" {
			log.Warnf("no index specified for writing for item: `%s`", msg.Data())
			continue
		}

		buf.Write(bulkHeader("create", index))
		buf.WriteString("\n")
		buf.Write(msg.Data())
		buf.WriteString("\n")
		inserting[i] = msg
	}
	params := opensearchapi.BulkParams{
		Timeout: 10 * time.Second,
	}
	req := opensearchapi.BulkReq{
		Body: bytes.NewReader(buf.Bytes()),
		Params: params,
	}
	log.Debugf("Inserting bulk into opensearch...")
	ctx, osSpan := osTracer.Start(ctx, "Bulk")
	resp, err := opensearch.Bulk(ctx, req)
	osSpan.End()
	if err != nil {
		log.Errorf("failed to send bulk message: %s", err)
		for _, m := range msgs {
			m.Msg.NakWithDelay(1*time.Minute)
		}
	}
	if resp.Errors {
		log.Debugf("Query: %s", buf.Bytes())
	}
	log.Debugf("opensearch returned result")
	for i, result := range resp.Items {
		res, ok := result["create"]
		if !ok {
			log.Warnf("cannot find action `index` in error response: %+v", result)
			continue
		}
		if res.Error != nil {
			log.Warnf("failed to write item to '%s': %s", res.Index, extractBulkError(res))
			if msg, ok := inserting[i]; ok {
				msg.Term()
			} else {
				log.Warnf("failed to term!")
			}
			errorItems.WithLabelValues(res.Index).Inc()
		} else {
			writeItems.WithLabelValues(res.Index).Inc()
			if msg, ok := inserting[i]; ok {
				msg.Ack()
			} else {
				log.Warnf("failed to ack!")
			}
		}
	}
}

func (w *Consumer) Stop() {
}
