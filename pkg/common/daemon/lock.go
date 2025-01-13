package daemon

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/common/redis"
)

type LockJob = func() error

const (
	retryInterval = 5 * time.Second
)

type LockDaemon struct {
	name     string
	lock     *redis.Mutex
	interval time.Duration
	fn       LockJob
}

func NewLockDaemon(name string, interval time.Duration, fn LockJob) *LockDaemon {
	return &LockDaemon{
		name:     name,
		lock:     redis.NewMutex(name),
		interval: interval,
		fn:       fn,
	}
}

func (d *LockDaemon) Run() error {
	ctx := context.TODO()

	for {
		// Acquire the lock (or wait for it forever)
		if err := d.lock.Acquire(ctx); err != nil {
			d.lock.Release()

			// Context cancel. This is to exit the for loop.
			return nil
		}

		log.Debugf("Executing job %s...", d.name)

		// Execute the job
		if err := d.fn(); err != nil {
			log.Errorf("failed to run job %s: %s", d.name, err)
			time.Sleep(retryInterval)

			continue
		}

		log.Debugf("Executed job %s successfully", d.name)

		time.Sleep(d.interval)
	}
}

func (d *LockDaemon) Stop() {
	d.lock.Release()
}
