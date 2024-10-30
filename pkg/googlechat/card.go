package googlechat

import (
	"fmt"

	chat "google.golang.org/api/chat/v1"

	"github.com/japannext/snooze/pkg/models"
)

func GetCard(item *models.Notification) *chat.Message {

	decoratedText := getDecoratedText(item)

	footer := &chat.GoogleAppsCardV1CardFixedFooter{
		PrimaryButton: &chat.GoogleAppsCardV1Button{
			Text: "Acknowledge",
			AltText: "Acknowledge this notification",
			Color: &chat.Color{Green: 0.8},
		},
		SecondaryButton: &chat.GoogleAppsCardV1Button{
			Text: "Snooze",
			AltText: "Open the snooze dialog",
			OnClick: &chat.GoogleAppsCardV1OnClick{
				Action: &chat.GoogleAppsCardV1Action{
					Interaction: "OPEN_DIALOG",
					Function: "",
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
					},
				},
				{
					Header: "Message",
					Widgets: []*chat.GoogleAppsCardV1Widget{
						{
							TextParagraph: &chat.GoogleAppsCardV1TextParagraph{Text: item.Message},
						},
					},
				},
			},
			FixedFooter: footer,
		},
	}

	return &chat.Message{
		CardsV2: []*chat.CardWithId{cardv2},
	}

}

// Use the "identity" of the notification to determine the decorated text labels
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
		Icon: &chat.GoogleAppsCardV1Icon{
			MaterialIcon: &chat.GoogleAppsCardV1MaterialIcon{Name: icon},
		},
		TopLabel: top,
		Text: mid,
		BottomLabel: bottom,
	}
}
