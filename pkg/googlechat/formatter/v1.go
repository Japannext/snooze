package formatter

import (
	"fmt"
	"strings"

	"github.com/japannext/snooze/pkg/models"
	chat "google.golang.org/api/chat/v1"
)

type V1 struct{}

func NewV1() *V1 {
	return &V1{}
}

func FormatWithoutFail(item *models.Notification) *chat.Message {
	var builder strings.Builder

	fmt.Fprintf(&builder, "*Date:* %s\n", item.NotificationTime)
	if value, ok := item.Identity["host"]; ok {
		fmt.Fprintf(&builder, "*Host:* %s\n", value)
	}
	if value, ok := item.Identity["process"]; ok {
		fmt.Fprintf(&builder, "*Process:* %s\n", value)
	}
	fmt.Fprintf(&builder, "*Source:* %s/%s\n", item.Source.Kind, item.Source.Name)

	// fmt.Fprintf(&builder, "*Severity:* %s\n", item.SeverityText)

	fmt.Fprintf(&builder, "*Message:* %s", item.Message)

	return &chat.Message{
		Text: builder.String(),
	}
}

func (v1 *V1) Format(item *models.Notification) (*chat.Message, error) {
	return FormatWithoutFail(item), nil
}
