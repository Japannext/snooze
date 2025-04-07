package writer

import (
	"bytes"
	"context"
	// "encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Run() error {
	for {
		log.Debugf("Ready to fetch items...")
		msgs, err := storeQ.Fetch(config.BatchSize)
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
		realCause := cause.Cause
		if realCause != nil && realCause.Reason != nil {
			fmt.Fprintf(&builder, ": [%s] %s", realCause.Type, *realCause.Reason)
		}
	}
	return fmt.Errorf("%s", builder.String())
}

/*
func bulkHeader(action, index string) []byte {
	m := map[string]map[string]string{
		action: {
			"_index": index,
		},
	}
	b, _ := json.Marshal(m)
	return b
}
*/

func bulkWrite(ctx context.Context, msgs []mq.MsgWithContext) {
	ctx, bulkSpan := tracer.Start(ctx, "bulkWrite")
	defer bulkSpan.End()
	var buf bytes.Buffer
	messages := map[int]jetstream.Msg{}
	for i, m := range msgs {
		msg, msgCtx := m.Extract()
		messages[i] = msg

		// Tracing
		msgCtx, span := tracer.Start(msgCtx, "bulkWrite", trace.WithLinks(trace.LinkFromContext(ctx)))
		defer span.End()

		buf.Write(msg.Data())
		buf.WriteString("\n")
	}
	params := opensearchapi.BulkParams{
		Timeout: 10 * time.Second,
	}
	req := opensearchapi.BulkReq{
		Body:   bytes.NewReader(buf.Bytes()),
		Params: params,
	}
	log.Debugf("Inserting bulk into opensearch...")
	resp, err := opensearch.Bulk(ctx, req)
	if err != nil {
		log.Errorf("failed to send bulk message: %s", err)
		log.Debugf("Query: %s", buf.Bytes())
		log.Debugf("Result: %+v", resp.Items)
		for _, m := range msgs {
			if err := m.Msg.NakWithDelay(1 * time.Minute); err != nil {
				log.Warn("fail to reschedule log")

				continue
			}
		}
	}
	if resp.Errors {
		log.Debugf("Query: %s", buf.Bytes())
		log.Debugf("Result: %+v", resp.Items)
	}
	for i, result := range resp.Items {
		if len(result) == 0 {
			continue
		}
		_, res := getActionAndResult(result)
		if res.Error != nil {
			log.Warnf("failed to write item to '%s': %s", res.Index, extractBulkError(res))
			if msg, ok := messages[i]; ok {
				if err := msg.Term(); err != nil {
					log.Warnf("failed to cancel log")

					continue
				}
			} else {
				log.Warnf("failed to term!")
			}
			errorItems.WithLabelValues(res.Index).Inc()
		} else {
			writeItems.WithLabelValues(res.Index).Inc()
			if msg, ok := messages[i]; ok {
				if err := msg.Ack(); err != nil {
					log.Warnf("failed to acknowledge log")

					continue
				}
			} else {
				log.Warnf("failed to ack!")
			}
		}
	}
}

func getActionAndResult(results map[string]opensearchapi.BulkRespItem) (string, opensearchapi.BulkRespItem) {
	for key, value := range results {
		return key, value
	}

	return "", opensearchapi.BulkRespItem{}
}

func (w *Consumer) Stop() {
}
