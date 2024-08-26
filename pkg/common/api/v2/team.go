package v2

type Team struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Members []string `json:"members"`
	NotificationProfiles []NotificationProfile
}

// Example: kind=mail,name=unix
type NotificationProfile struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
