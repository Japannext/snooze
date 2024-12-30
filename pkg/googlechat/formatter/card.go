package formatter

import (
	"fmt"

	"github.com/japannext/snooze/pkg/models"
	chat "google.golang.org/api/chat/v1"
)

// Format things within a googlechat CardV2.
type Card struct{}

func NewCard() *Card {
	return &Card{}
}

func (p *Card) Format(item *models.Notification) (*chat.Message, error) {
	decoratedText := getDecoratedText(item)

	footer := &chat.GoogleAppsCardV1CardFixedFooter{
		PrimaryButton: &chat.GoogleAppsCardV1Button{
			Text:    "Acknowledge",
			AltText: "Acknowledge this notification",
			Color:   &chat.Color{Green: 0.8},
			Icon:    materialIcon("check"),
		},
		SecondaryButton: &chat.GoogleAppsCardV1Button{
			Text:    "Snooze",
			AltText: "Open the snooze dialog",
			OnClick: &chat.GoogleAppsCardV1OnClick{
				Action: &chat.GoogleAppsCardV1Action{
					Interaction: "OPEN_DIALOG",
					Function:    "",
				},
			},
		},
	}

	cardv2 := &chat.CardWithId{
		Card: &chat.GoogleAppsCardV1Card{
			Sections: []*chat.GoogleAppsCardV1Section{
				{
					Header: "Identity",
					Widgets: []*chat.GoogleAppsCardV1Widget{
						{
							DecoratedText: decoratedText,
						},
						{
							TextParagraph: &chat.GoogleAppsCardV1TextParagraph{Text: fmt.Sprintf("<b>Message:</b> %s", item.Message)},
						},
					},
				},
			},
			FixedFooter: footer,
		},
	}

	return &chat.Message{
		CardsV2: []*chat.CardWithId{cardv2},
	}, nil
}

// Use the "identity" of the notification to determine the decorated text labels.
func getIdentityLabels(identity map[string]string) (icon, top string, mid string, bottom string) {
	icon = "help"

	if host, ok := identity["host"]; ok {
		mid = host
		icon = "host"
		if process, ok := identity["process"]; ok {
			bottom = process
		}
		return
	}
	if deploy, ok := identity["k8s.deployment"]; ok {
		icon = "deployed_code"
		mid = fmt.Sprintf("deploy/%s", deploy)
		bottom = identity["k8s.namespace"]
		top = identity["k8s.cluster"]
		return
	}
	return
}

func getDecoratedText(item *models.Notification) *chat.GoogleAppsCardV1DecoratedText {
	icon, top, mid, bottom := getIdentityLabels(item.Identity)

	return &chat.GoogleAppsCardV1DecoratedText{
		Icon:        materialIcon(icon),
		TopLabel:    top,
		Text:        mid,
		BottomLabel: bottom,
	}
}

func materialIcon(name string) *chat.GoogleAppsCardV1Icon {
	return &chat.GoogleAppsCardV1Icon{MaterialIcon: &chat.GoogleAppsCardV1MaterialIcon{Name: name}}
}
