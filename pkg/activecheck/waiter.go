package activecheck

import (
	"fmt"
	"sync"
	"time"

	"github.com/japannext/snooze/pkg/models"
)

// An global object to wait for multiple check to callback
// through a webhook

type Waiter struct {
	mu sync.Mutex
	m  map[string](chan models.SourceActiveCheck)
}

func NewWaiter() *Waiter {
	return &Waiter{
		m: make(map[string](chan models.SourceActiveCheck)),
	}
}

// Get a channel or initialize it.
func (waiter *Waiter) getChannel(key string) chan models.SourceActiveCheck {
	waiter.mu.Lock()
	ch, found := waiter.m[key]
	if !found {
		ch = make(chan models.SourceActiveCheck)
		waiter.m[key] = ch
	}
	waiter.mu.Unlock()
	return ch
}

// Cleanup the key to avoid growing the map infinitely.
func (waiter *Waiter) cleanup(key string) {
	waiter.mu.Lock()
	delete(waiter.m, key)
	waiter.mu.Unlock()
}

func (waiter *Waiter) Wait(key string, timeout time.Duration) (models.SourceActiveCheck, error) {
	ch := waiter.getChannel(key)
	defer waiter.cleanup(key)
	select {
	case callback := <-ch:
		return callback, nil
	case <-time.After(timeout):
		return models.SourceActiveCheck{}, fmt.Errorf("timed out while waiting for callback")
	}
}

func (waiter *Waiter) Insert(key string, callback models.SourceActiveCheck) {
	ch := waiter.getChannel(key)
	ch <- callback
}

var waiter *Waiter

func init() {
	waiter = NewWaiter()
}
