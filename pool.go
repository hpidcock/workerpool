package workerpool

import (
	"context"
	"runtime"
)

// WorkerPool provides various parallelism utilities via on a pool of goroutines.
type WorkerPool struct {
	acquireWorker chan interface{}
	jobChan       chan func()
	ctx           context.Context
	count         int
}

// New returns a WorkerPool with maxWorker goroutines. Cancelation facilitated via the context.
func New(ctx context.Context, maxWorkers int) *WorkerPool {
	count := runtime.NumCPU() * 10
	if maxWorkers > 0 {
		count = maxWorkers
	}
	pw := &WorkerPool{
		ctx:           ctx,
		jobChan:       make(chan func(), 0),
		acquireWorker: make(chan interface{}, count),
		count:         count,
	}
	for i := 0; i < count; i++ {
		pw.acquireWorker <- nil
		go pw.run()
	}
	go pw.close()
	return pw
}

// Push enqueues a function to be run on a worker goroutine.
func (wp *WorkerPool) Push(fn func()) error {
	_, ok := <-wp.acquireWorker
	if ok == false {
		return context.Canceled
	}
	wp.jobChan <- fn
	return nil
}

func (wp *WorkerPool) close() {
	<-wp.ctx.Done()
	for i := 0; i < wp.count; i++ {
		<-wp.acquireWorker
	}
	close(wp.acquireWorker)
	close(wp.jobChan)
}

func (wp *WorkerPool) run() {
	closeChan := wp.ctx.Done()
	for {
		select {
		case <-closeChan:
			return
		case job, ok := <-wp.jobChan:
			if ok == false {
				return
			}
			job()
			wp.acquireWorker <- nil
		}
	}
}
