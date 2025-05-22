package daemon

type Ping struct{}

func (ping *Ping) Pass() error {
	return nil
}

func (ping *Ping) Name() string {
	return "ping"
}

func NewPing() *Ping {
	return &Ping{}
}
