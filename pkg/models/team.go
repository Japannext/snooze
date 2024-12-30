package models

type Team struct {
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	Members              []string              `json:"members"`
	NotificationProfiles []NotificationProfile `json:"notificationProfiles"`
}

type NotificationProfile struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
