package exporter

import (
	"fmt"
	"time"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

// An global object to wait for multiple check to callback
// through a webhook
type Waiter struct {
	Timeout time.Duration
	mapper map[string](chan *api.ActiveCheckCallback)
}

func NewWaiter(timeout time.Duration) *Waiter {
	return &Waiter{
		Timeout: timeout,
		mapper: map[string](chan *api.ActiveCheckCallback){},
	}
}

func (w *Waiter) Prepare(key string) error {
	if _, ok := w.mapper[key]; ok {
		return fmt.Errorf("waiting more than once for the same key: %s", key)
	}
	w.mapper[key] = make(chan *api.ActiveCheckCallback)
	return nil
}

func (w *Waiter) Wait(key string) (*api.ActiveCheckCallback, error) {
	select {
	case callback := <-w.mapper[key]:
		return callback, nil
	case <-time.After(w.Timeout):
		return nil, fmt.Errorf("timeout while waiting for key %s", key)
	}
}

func (w *Waiter) Insert(key string, callback *api.ActiveCheckCallback) {
	ch, ok := w.mapper[key]
	if !ok {
		return
	}
	ch <- callback
}

func (w *Waiter) Cleanup(key string) {
	delete(w.mapper, key)
}

var waiter = NewWaiter(time.Second * 30)
