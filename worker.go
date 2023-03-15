package workers

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

const (
	defaultTotalWorkers int64 = 10
)

var (
	// ErrNoWorkers error for when there is no workers running and a job was sent to be done
	ErrNoWorkers error = errors.New("no workers running")
	doJob        chan Worker
	counter      int64
	totalWorkers int64
)

// Options configuration to the lib.
// TotalWorkers is the total workers available to do a job. Default is 10
// BufferChannel the size of the buffer for the channel to distribute the jobs for the workers
type Options struct {
	TotalWorkers  int
	BufferChannel int
}

// Worker interface to run worker
type Worker interface {
	Work()
}

// Job type to implement the Worker interface
type Job func()

// Work Job is a function that will be executed by workers
func (j Job) Work() {
	j()
}

// Start sets up the configuration and run the workers.
// If no option is set the default is 10 workers and unbuffered channel
func Start(ctx context.Context, opts ...Options) {
	opt := getOption(opts)
	totalWorkers = int64(opt.TotalWorkers)
	if totalWorkers == 0 {
		totalWorkers = defaultTotalWorkers
	}
	doJob = make(chan Worker, opt.BufferChannel)
	counter = 0
	var i int64
	for i = 0; i < totalWorkers; i++ {
		go runWorker(ctx)
	}
	go clean(ctx)
}

func clean(ctx context.Context) {
	<-ctx.Done()
	for {
		if Alive() {
			time.Sleep(time.Millisecond)
			continue
		}
		break
	}
	close(doJob)
}

// DoJob sends job for workers to do
func DoJob(w Worker) error {
	if !Alive() {
		return ErrNoWorkers
	}

	doJob <- w
	return nil
}

// Alive checks if the workers are alive
func Alive() bool {
	return atomic.LoadInt64(&counter) < totalWorkers
}

func runWorker(ctx context.Context) {
	defer func() {
		atomic.AddInt64(&counter, 1)
	}()

	for {
		select {
		case j, ok := <-doJob:
			if !ok {
				return
			}
			if j == nil {
				continue
			}
			j.Work()
		case <-ctx.Done():
			return
		}
	}
}

func getOption(opts []Options) (opt Options) {
	if len(opts) == 0 {
		return
	}
	opt = opts[0]
	return
}
