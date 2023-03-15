package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/gbeletti/workers"
)

func TestDoJob(t *testing.T) {
	timeout := time.Millisecond * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() {
		cancel()
	})
	opt := workers.Options{
		TotalWorkers:  3, // sets up 3 workers for the test
		BufferChannel: 0, // channel will be unbuffered
	}
	workers.Start(ctx, opt)

	if !workers.Alive() {
		t.Fatal("workers should be alive")
	}
	var job workers.Job = func() {
		t.Log("doing something")
	}

	// sends a job to the workers to do
	if err := workers.DoJob(job); err != nil {
		t.Fatal("expected to send a job for the workers, instead there is no workers running")
	}

	// wait timeout and the workers to stop
	time.Sleep(time.Millisecond * 25)

	if err := workers.DoJob(job); err == nil {
		t.Error("all workers should have stopped, no job should have been done")
	}

	if workers.Alive() {
		t.Error("all workers should have stopped")
	}
}

func TestDoFunc(t *testing.T) {
	timeout := time.Millisecond * 20
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(func() {
		cancel()
	})
	opt := workers.Options{
		TotalWorkers:  3, // sets up 3 workers for the test
		BufferChannel: 0, // channel will be unbuffered
	}
	workers.Start(ctx, opt)

	if !workers.Alive() {
		t.Fatal("workers should be alive")
	}
	var job = func() {
		t.Log("doing something")
	}

	// sends a job to the workers to do
	if err := workers.DoFunc(job); err != nil {
		t.Fatal("expected to send a job for the workers, instead there is no workers running")
	}

	// wait timeout and the workers to stop
	time.Sleep(time.Millisecond * 25)

	if err := workers.DoFunc(job); err == nil {
		t.Error("all workers should have stopped, no job should have been done")
	}

	if workers.Alive() {
		t.Error("all workers should have stopped")
	}
}
