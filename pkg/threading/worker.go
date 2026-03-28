package threading

import (
	"context"
	"log/slog"
	"sync"
)

// Worker is used to handle work of the [WorkerPool].
type Worker struct {
	id            int
	active        bool
	activeMu      *sync.RWMutex
	isDoingWork   bool
	isDoingWorkMu *sync.RWMutex
	pool          *WorkerPool
}

// NewWorker creates a new [Worker].
func NewWorker(id int, pool *WorkerPool) Worker {
	worker := Worker{
		id:            id,
		active:        false,
		activeMu:      &sync.RWMutex{},
		isDoingWork:   false,
		isDoingWorkMu: &sync.RWMutex{},
		pool:          pool,
	}
	return worker
}

// Active fetches the current state of the worker.
func (worker *Worker) Active() bool {
	worker.activeMu.RLock()
	defer worker.activeMu.RUnlock()

	return worker.active
}

// IsDoingWork fetches the current state of the worker.
func (worker *Worker) IsDoingWork() bool {
	worker.isDoingWorkMu.RLock()
	defer worker.isDoingWorkMu.RUnlock()

	return worker.isDoingWork
}

// Stop stops the worker.
func (worker *Worker) Stop() {
	worker.activeMu.Lock()
	defer worker.activeMu.Unlock()

	worker.active = false
}

// Run makes [Worker] start doing work.
func (worker *Worker) Run(ctx context.Context, logger *slog.Logger) error {
	worker.activeMu.Lock()
	worker.active = true
	worker.activeMu.Unlock()

	for worker.Active() {
		doWork := <-worker.pool.queue

		worker.isDoingWorkMu.Lock()
		worker.isDoingWork = true
		worker.isDoingWorkMu.Unlock()

		err := doWork(ctx, logger)
		if err != nil {
			logger.ErrorContext(ctx, err.Error())
		}

		worker.isDoingWorkMu.Lock()
		worker.isDoingWork = false
		worker.isDoingWorkMu.Unlock()
	}

	return nil
}
