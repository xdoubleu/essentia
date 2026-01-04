package threading_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xdoubleu/essentia/v2/pkg/logging"
	"github.com/xdoubleu/essentia/v2/pkg/threading"
)

func doWork(_ context.Context, _ *slog.Logger) error {
	time.Sleep(1 * time.Second)
	return nil
}

func TestBasicWorkerPool(t *testing.T) {
	workerpool := threading.NewWorkerPool(logging.NewNopLogger(), 1, 2)

	workerpool.EnqueueWork(doWork)
	workerpool.EnqueueWork(doWork)

	workerpool.WaitUntilDone()
	assert.False(t, workerpool.IsDoingWork())
}
