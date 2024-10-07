package models

type Mute struct {
    // Skip the notification step. Usually on.
    SkipNotification bool `json:"skipNotification"`
    // Skip storing into the database (opensearch). Usually used for testing
    SkipStorage bool `json:"skipStorage"`
	// Reason why it was muted
	Reason string `json:"reason"`
}

func (m *Mute) Silence(reason string) {
	m.SkipNotification = true
	m.SkipStorage = false
	m.Reason = reason
}

func (m *Mute) Drop(reason string) {
	m.SkipNotification = true
	m.SkipStorage = true
	m.Reason = reason
}

func (m *Mute) Enabled() bool {
	return m.SkipNotification || m.SkipStorage
}
