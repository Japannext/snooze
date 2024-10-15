package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/opensearch"
)

type Consumer struct {}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Run() error {
	for {
		batch, err := storeQ.Fetch(config.BatchSize, jetstream.FetchMaxWait(5*time.Second))
		if err != nil {
			log.Warnf("failed to fetch items: %s", err)
			continue
		}
		var msgs = []jetstream.Msg{}
		for msg := range batch.Messages() {
			msgs = append(msgs, msg)
		}
		if len(msgs) == 0 {
			continue
		}
		ctx := context.TODO()
		bulkWrite(ctx, msgs)
	}
}

func bulkWrite(ctx context.Context, msgs []jetstream.Msg) {
	//var buf bytes.Buffer
	reader, writer := io.Pipe()
	var inserting map[int]jetstream.Msg
	for i, msg := range msgs {
		msg.InProgress()
		splits := strings.Split(msg.Subject(), ".")
		if len(splits) < 1 {
			// TODO
			errorItems.WithLabelValues("unknown").Inc()
			continue
		}
		indexString, err := json.Marshal(splits[1])
		if err != nil {
			// TODO
			errorItems.WithLabelValues("unknown").Inc()
			continue
		}
		fmt.Fprintf(writer, `{"index": {"_index": %s}}`+"\n", indexString)
		writer.Write(msg.Data())
		writer.Write([]byte("\n"))
		inserting[i] = msg
	}
	req := opensearchapi.BulkReq{
		Body: reader,
	}
	resp, err := opensearch.Bulk(ctx, req)
	if err != nil {
		log.Errorf("failed to send bulk message: %s", err)
		for _, msg := range msgs {
			msg.NakWithDelay(1*time.Minute)
		}
	}
	if resp.Errors {
		for i, actionMap := range resp.Items {
			item, ok := actionMap["index"]
			if !ok {
				// TODO
				continue
			}
			if item.Error != nil {
				msg, ok := inserting[i]
				if !ok {
					// TODO
					continue
				}
				msg.Term()
				errorItems.WithLabelValues().Inc()
			} else {
				//writeItems.Inc()
			}
		}
	} else {
		// No error: ACK everything
		for _, msg := range msgs {
			msg.Ack()
		}
		// writeItems.Add(float64(len(msgs)))
	}
}

func (w *Consumer) Stop() {
}
