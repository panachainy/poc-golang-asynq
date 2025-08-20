package background_job

import (
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"
)

type BackgroundJob interface {
	Submit(task func()) error
	Wait()
}

type antsBackgroundJob struct {
	pool *ants.Pool
}

func NewAntsBackgroundJob(size int, opts ...ants.Option) (BackgroundJob, error) {
	p, err := ants.NewPool(size, opts...)
	if err != nil {
		return nil, err
	}

	instance := &antsBackgroundJob{pool: p}
	instance.runMonitor()

	return instance, nil
}

func (a *antsBackgroundJob) Submit(task func()) error {
	return a.pool.Submit(task)
}

func (a *antsBackgroundJob) Wait() {
	for a.pool.Running() > 0 {
		time.Sleep(1000 * time.Millisecond)
	}
}

func (a *antsBackgroundJob) runMonitor() {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("Pool stats - Running: %d, Free: %d, Cap: %d\n",
					a.pool.Running(), a.pool.Free(), a.pool.Cap())
			}
		}
	}()
}
