package threading

import (
	"context"
	"log/slog"
	"sync"
)

// Subscriber describes the interface of subscribers of new events.
type Subscriber interface {
	ID() string
	OnEventCallback(event any)
}

// EventQueue is used to divide [Subscriber]s between [Worker]s.
// This prevents one [Worker] of being very busy.
type EventQueue struct {
	workerPool    WorkerPool
	subscribers   []Subscriber
	subscribersMu *sync.RWMutex
}

// NewEventQueue creates a new [EventQueue].
func NewEventQueue(
	ctx context.Context,
	logger *slog.Logger,
	maxWorkers int,
	channelBufferSize int,
) *EventQueue {
	pool := &EventQueue{
		workerPool:    *NewWorkerPool(ctx, logger, maxWorkers, channelBufferSize),
		subscribers:   []Subscriber{},
		subscribersMu: &sync.RWMutex{},
	}

	return pool
}

// EnqueueEvent puts an event on the [Worker] channels.
func (q *EventQueue) EnqueueEvent(event any) {
	q.workerPool.EnqueueWork(func(ctx context.Context, logger *slog.Logger) error {
		q.processEvent(ctx, logger, event)
		return nil
	})
}

// AddSubscriber adds a [Subscriber] to the [EventQueue].
func (q *EventQueue) AddSubscriber(sub Subscriber) {
	q.subscribersMu.Lock()
	defer q.subscribersMu.Unlock()

	q.subscribers = append(q.subscribers, sub)
}

// RemoveSubscriber removes a [Subscriber] from the [EventQueue].
func (q *EventQueue) RemoveSubscriber(sub Subscriber) {
	q.subscribersMu.Lock()
	defer q.subscribersMu.Unlock()

	var i int
	for i = range q.subscribers {
		if q.subscribers[i].ID() != sub.ID() {
			continue
		}
		break
	}

	// delete subscriber
	q.subscribers[i] = q.subscribers[len(q.subscribers)-1]
	q.subscribers = q.subscribers[:len(q.subscribers)-1]
}

func (q *EventQueue) processEvent(_ context.Context, _ *slog.Logger, event any) {
	q.subscribersMu.RLock()

	for _, sub := range q.subscribers {
		sub.OnEventCallback(event)
	}

	q.subscribersMu.RUnlock()
}
