package threading_test

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/logging"
	"github.com/xdoubleu/essentia/v2/pkg/threading"
)

type TestJob struct {
}

func (j TestJob) ID() string {
	return "test"
}

func (j TestJob) Run(_ context.Context, _ *slog.Logger) error {
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (j TestJob) RunEvery() time.Duration {
	return 500 * time.Millisecond
}

func TestJobQueueSimple(t *testing.T) {
	jobQueue := threading.NewJobQueue(logging.NewNopLogger(), 1, 1)

	statesMu := sync.Mutex{}
	states := []bool{}

	err := jobQueue.AddJob(
		TestJob{},
		func(_ string, isRunning bool, _ *time.Time) {
			statesMu.Lock()
			states = append(states, isRunning)
			statesMu.Unlock()
		},
	)
	assert.Nil(t, err)

	jobIDs := jobQueue.FetchJobIDs()
	assert.Equal(t, []string{"test"}, jobIDs)

	state, _ := jobQueue.FetchState("test")
	assert.Equal(t, true, state)

	time.Sleep(400 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false}, states)
	statesMu.Unlock()

	time.Sleep(500 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false, true, false}, states)
	statesMu.Unlock()
}

func TestJobQueueSimpleAfterClear(t *testing.T) {
	jobQueue := threading.NewJobQueue(logging.NewNopLogger(), 1, 1)

	statesMu := sync.Mutex{}
	states := []bool{}

	err := jobQueue.AddJob(
		TestJob{},
		func(_ string, isRunning bool, _ *time.Time) {
			statesMu.Lock()
			states = append(states, isRunning)
			statesMu.Unlock()
		},
	)
	assert.Nil(t, err)

	time.Sleep(1 * time.Millisecond)
	jobQueue.Clear()

	time.Sleep(400 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false}, states)
	statesMu.Unlock()

	err = jobQueue.AddJob(
		TestJob{},
		func(_ string, isRunning bool, _ *time.Time) {
			statesMu.Lock()
			states = append(states, isRunning)
			statesMu.Unlock()
		},
	)
	assert.Nil(t, err)

	time.Sleep(400 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false, true, false}, states)
	statesMu.Unlock()
}

func TestJobQueueForce(t *testing.T) {
	jobQueue := threading.NewJobQueue(logging.NewNopLogger(), 1, 1)

	statesMu := sync.Mutex{}
	states := []bool{}

	err := jobQueue.AddJob(
		TestJob{},
		func(_ string, isRunning bool, _ *time.Time) {
			statesMu.Lock()
			states = append(states, isRunning)
			statesMu.Unlock()
		},
	)
	assert.Nil(t, err)

	time.Sleep(400 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false}, states)
	statesMu.Unlock()

	jobQueue.ForceRun("test")

	state, _ := jobQueue.FetchState("test")
	assert.Equal(t, true, state)

	time.Sleep(500 * time.Millisecond)
	statesMu.Lock()
	assert.Equal(t, []bool{true, false, true, false}, states)
	statesMu.Unlock()
}
