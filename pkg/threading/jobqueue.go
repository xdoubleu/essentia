package threading

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// CallbackFunc describes the interface for a func
// called before and after running a job.
type CallbackFunc = func(id string, isRunning bool, lastRunTime *time.Time)

// JobQueue is a queue of Jobs which will be executed by the workerpool.
type JobQueue struct {
	workerPool      WorkerPool
	logger          *slog.Logger
	recurringJobs   map[string]*jobContainer
	jobsMu          sync.RWMutex
	schedulerActive bool
	schedulerMu     sync.RWMutex
}

// Job describes the interface for a job executable by the jobqueue.
type Job interface {
	ID() string
	Run(context.Context, *slog.Logger) error
	RunEvery() time.Duration
}

type jobContainer struct {
	job         Job
	period      time.Duration
	lastRunTime *time.Time
	callback    CallbackFunc
	isPushed    bool
	mu          sync.RWMutex
}

// NewJobQueue creates a new jobqueue.
func NewJobQueue(
	ctx context.Context,
	logger *slog.Logger,
	amountWorkers int,
	size int,
) *JobQueue {
	jobQueue := &JobQueue{
		workerPool:      *NewWorkerPool(ctx, logger, amountWorkers, size),
		logger:          logger,
		recurringJobs:   make(map[string]*jobContainer),
		schedulerActive: false,
		jobsMu:          sync.RWMutex{},
		schedulerMu:     sync.RWMutex{},
	}

	return jobQueue
}

// Clear clears the JobQueue completely.
func (q *JobQueue) Clear() {
	if q.isSchedulerActive() {
		q.schedulerMu.Lock()
		q.schedulerActive = false
		q.schedulerMu.Unlock()
	}

	if q.isWorkerActive() {
		q.workerPool.Stop()
	}

	q.jobsMu.Lock()
	q.recurringJobs = make(map[string]*jobContainer)
	q.jobsMu.Unlock()
}

// AddJob adds a recurring job which should be executed by the workerpool.
// This will also execute the job.
func (q *JobQueue) AddJob(job Job, callback CallbackFunc) error {
	jobContainer := &jobContainer{
		job:         job,
		period:      job.RunEvery(),
		callback:    callback,
		lastRunTime: nil,
		isPushed:    false,
		mu:          sync.RWMutex{},
	}

	q.jobsMu.Lock()
	defer q.jobsMu.Unlock()

	_, ok := q.recurringJobs[job.ID()]
	if ok {
		return errors.New("a job with this ID already exists")
	}

	q.recurringJobs[job.ID()] = jobContainer

	q.push(jobContainer)
	return nil
}

// ForceRun forces a run of the specified job.
func (q *JobQueue) ForceRun(id string) {
	q.jobsMu.RLock()
	defer q.jobsMu.RUnlock()

	rj, ok := q.recurringJobs[id]
	if !ok {
		return
	}
	q.push(rj)
}

// FetchJobIDs fetches all IDs for all jobs.
func (q *JobQueue) FetchJobIDs() []string {
	q.jobsMu.RLock()
	defer q.jobsMu.RUnlock()

	result := []string{}
	for _, rj := range q.recurringJobs {
		result = append(result, rj.job.ID())
	}
	return result
}

// FetchState fetches the current state of the specified job.
func (q *JobQueue) FetchState(id string) (bool, *time.Time) {
	q.jobsMu.RLock()
	defer q.jobsMu.RUnlock()

	rj, ok := q.recurringJobs[id]
	if !ok {
		return false, nil
	}

	rj.mu.RLock()
	defer rj.mu.RUnlock()

	return rj.isPushed, rj.lastRunTime
}

func (q *JobQueue) push(jobContainer *jobContainer) {
	if !q.isSchedulerActive() {
		q.startScheduler()
	}

	if !q.isWorkerActive() {
		q.workerPool.Start()
	}

	jobContainer.mu.Lock()
	jobContainer.isPushed = true
	jobContainer.mu.Unlock()

	q.workerPool.EnqueueWork(jobContainer.run)
}

func (q *JobQueue) startScheduler() {
	q.schedulerMu.Lock()
	q.schedulerActive = true
	q.schedulerMu.Unlock()

	go func() {
		for q.isSchedulerActive() {
			q.jobsMu.RLock()
			for k := range q.recurringJobs {
				job := q.recurringJobs[k]
				if job.shouldRun() {
					q.push(job)
				}
			}

			//nolint:mnd //1/10th of period avoids boundary-condition flakiness
			sleep := getSmallestPeriod(q.recurringJobs) / 10
			if sleep < 50*time.Millisecond {
				//nolint:mnd //50ms minimum tick
				sleep = 50 * time.Millisecond
			}
			q.jobsMu.RUnlock()

			time.Sleep(sleep)
		}
	}()
}

func (q *JobQueue) isSchedulerActive() bool {
	q.schedulerMu.Lock()
	defer q.schedulerMu.Unlock()

	return q.schedulerActive
}

func (q *JobQueue) isWorkerActive() bool {
	return q.workerPool.Active()
}

func getSmallestPeriod(jobContainers map[string]*jobContainer) time.Duration {
	var smallestPeriod *time.Duration

	for _, c := range jobContainers {
		if smallestPeriod == nil ||
			c.period.Nanoseconds() < smallestPeriod.Nanoseconds() {
			smallestPeriod = &c.period
		}
	}

	if smallestPeriod == nil {
		//nolint:mnd //no magic number
		return 10 * time.Second
	}

	return *smallestPeriod
}

func (c *jobContainer) run(ctx context.Context, logger *slog.Logger) error {
	defer func() {
		c.mu.Lock()
		c.isPushed = false
		c.mu.Unlock()
	}()

	c.mu.RLock()
	c.callback(c.job.ID(), true, c.lastRunTime)
	c.mu.RUnlock()

	c.mu.Lock()
	nowUTC := time.Now().UTC()
	c.lastRunTime = &nowUTC
	c.mu.Unlock()

	logger.DebugContext(ctx, fmt.Sprintf("started job %s", c.job.ID()))
	err := c.job.Run(ctx, logger)
	if err != nil {
		return err
	}
	logger.DebugContext(ctx, fmt.Sprintf("successfully finished job %s", c.job.ID()))

	c.mu.RLock()
	c.callback(c.job.ID(), false, c.lastRunTime)
	c.mu.RUnlock()

	return nil
}

func (c *jobContainer) shouldRun() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isPushed {
		return false
	}

	if c.lastRunTime == nil {
		return true
	}

	now := time.Now().UTC()

	if c.period >= 24*time.Hour {
		return now.Day() > c.lastRunTime.Day()
	}

	if c.period >= time.Hour {
		return now.Day() > c.lastRunTime.Day() || now.Hour() > c.lastRunTime.Hour()
	}

	return time.Now().UTC().After(c.lastRunTime.Add(c.period))
}
