package threading

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/xdoubleu/essentia/pkg/sentry"
)

// DoWork describes the interface for work executed by the workers.
type DoWork = func(ctx context.Context, logger *slog.Logger) error

// WorkerPool is used to divide [Subscriber]s between [Worker]s.
// This prevents one [Worker] of being very busy.
type WorkerPool struct {
	logger  *slog.Logger
	workers []Worker
	queue   chan DoWork
}

// NewWorkerPool creates a new [WorkerPool].
func NewWorkerPool(
	logger *slog.Logger,
	amountWorkers int,
	queueSize int,
) *WorkerPool {
	pool := &WorkerPool{
		logger:  logger,
		workers: make([]Worker, amountWorkers),
		queue:   make(chan DoWork, queueSize),
	}

	pool.createWorkers(amountWorkers)
	pool.Start()

	return pool
}

// Active checks if the [WorkerPool] is active
// by checking if any [Worker] is active.
func (pool *WorkerPool) Active() bool {
	for i := range pool.workers {
		if pool.workers[i].Active() {
			return true
		}
	}
	return false
}

// IsDoingWork checks if the [WorkerPool] is still processing work.
func (pool *WorkerPool) IsDoingWork() bool {
	for i := range pool.workers {
		if pool.workers[i].IsDoingWork() {
			return true
		}
	}
	return false
}

// Start starts [Worker]s of a [WorkerPool] if they weren't active yet.
func (pool *WorkerPool) Start() {
	for i := range pool.workers {
		go sentry.GoRoutineWrapper(
			context.Background(),
			pool.logger,
			fmt.Sprintf("Worker %d", i),
			pool.workers[i].Run,
		)
	}
}

// EnqueueWork puts an work on the queue.
func (pool *WorkerPool) EnqueueWork(doWork DoWork) {
	pool.queue <- doWork
}

// IsWorkRemaining checks if there is still work on the queue.
func (pool *WorkerPool) IsWorkRemaining() bool {
	return len(pool.queue) > 0 || pool.IsDoingWork()
}

// WaitUntilDone blocks until the queue is empty.
func (pool *WorkerPool) WaitUntilDone() {
	for pool.IsWorkRemaining() {
		//nolint:mnd //no magic number
		time.Sleep(100 * time.Millisecond)
	}
}

// Stop stops all workers.
func (pool *WorkerPool) Stop() {
	for i := range pool.workers {
		pool.workers[i].Stop()
	}
}

func (pool *WorkerPool) createWorkers(amountWorkers int) {
	for i := 0; i < amountWorkers; i++ {
		pool.workers[i] = NewWorker(i, pool)
	}
}
