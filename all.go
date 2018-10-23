package workerpool

import (
	"context"
	"sync"
)

// All runs tasks in the pool, waiting for all to succeed, or cancel the context
// on the first to return a non-nil error.
func (wp *WorkerPool) All(ctx context.Context, fn ...Task) (eout error) {
	ctx, cancel := context.WithCancel(ctx)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	terminate := func(err error) {
		m.Lock()
		defer m.Unlock()
		if eout != nil {
			return
		}
		eout = err
		cancel()
	}
	for _, f := range fn {
		wg.Add(1)
		ff := f
		err := wp.Push(func() {
			err := ff(ctx)
			if err != nil {
				terminate(err)
			}
			wg.Done()
		})
		if err != nil {
			terminate(err)
			wg.Done()
			return
		}
	}
	wg.Wait()
	return
}
