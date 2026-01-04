package threading_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/logging"
	"github.com/xdoubleu/essentia/v2/pkg/threading"
)

type TestSubscriber struct {
	id       string
	output   string
	outputMu *sync.RWMutex
}

func NewTestSubscriber() *TestSubscriber {
	return &TestSubscriber{
		id:       uuid.NewString(),
		output:   "",
		outputMu: &sync.RWMutex{},
	}
}

func (sub TestSubscriber) ID() string {
	return sub.id
}

func (sub *TestSubscriber) OnSubscribeCallback() any {
	return nil
}

func (sub *TestSubscriber) OnEventCallback(event any) {
	sub.outputMu.Lock()
	defer sub.outputMu.Unlock()

	if v, ok := event.(string); ok {
		sub.output = v
	}
}

func (sub *TestSubscriber) Output() string {
	sub.outputMu.RLock()
	defer sub.outputMu.RUnlock()

	return sub.output
}

const sleep = 100 * time.Millisecond

func TestBasic(t *testing.T) {
	logger := logging.NewNopLogger()

	wp := threading.NewEventQueue(logger, 1, 10)

	tSub := NewTestSubscriber()
	wp.AddSubscriber(tSub)
	time.Sleep(sleep)

	event := "Hello, World!"
	wp.EnqueueEvent(event)
	time.Sleep(sleep)
	assert.Equal(t, event, tSub.Output())

	wp.RemoveSubscriber(tSub)
}
