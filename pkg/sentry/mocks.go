package sentry

import (
	"context"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

const testDsn = "http://whatever@example.com/1337"

// MockedSentryClientOptions returns a mocked version of [sentry.ClientOptions].
func MockedSentryClientOptions() sentry.ClientOptions {
	//nolint:exhaustruct //other fields are optional
	return sentry.ClientOptions{
		Dsn:       testDsn,
		Transport: newTransportMock(),
	}
}

// MockedSentryHub returns a mocked version of [*sentry.Hub].
func MockedSentryHub() *sentry.Hub {
	clientOptions := MockedSentryClientOptions()

	client, err := sentry.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	scope := sentry.NewScope()
	hub := sentry.NewHub(client, scope)
	return hub
}

type transportMock struct {
	mu        sync.Mutex
	events    []*sentry.Event
	lastEvent *sentry.Event
}

func newTransportMock() sentry.Transport {
	return &transportMock{
		mu:        sync.Mutex{},
		events:    []*sentry.Event{},
		lastEvent: nil,
	}
}

func (t *transportMock) FlushWithContext(_ context.Context) bool { return true }
func (t *transportMock) Configure(_ sentry.ClientOptions)        {}
func (t *transportMock) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}
func (t *transportMock) Flush(_ time.Duration) bool {
	return true
}
func (t *transportMock) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}
func (t *transportMock) Close() {
}
