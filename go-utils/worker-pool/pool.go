package workerpool

import (
	"context"
	"fmt"
)

type Pool[T any] struct {
	jobChan    chan Job[T]
	outputChan chan T
	workers    []worker[T]
}

func NewPool[T any](workersNum int) Pool[T] {
	p := Pool[T]{
		jobChan:    make(chan Job[T], workersNum),
		outputChan: make(chan T, workersNum),
		workers:    make([]worker[T], workersNum),
	}
	for i := 0; i < workersNum; i++ {
		p.workers[i] = newWorker(p.jobChan, p.outputChan)
	}
	return p
}

func (p Pool[T]) Start(ctx context.Context, jobs []Job[T]) []T {
	for i := range p.workers {
		go p.workers[i].start(ctx)
	}

	go func() {
		for i := range jobs {
			p.jobChan <- jobs[i]
		}
		close(p.jobChan)
	}()

	outputs := make([]T, 0, len(jobs))
	for i := 0; i < len(jobs); i++ {
		select {
		case output := <-p.outputChan:
			outputs = append(outputs, output)
		case <-ctx.Done():
			break
		}
	}
	return outputs
}

type worker[T any] struct {
	jobs <-chan Job[T]
	out  chan<- T
}

func newWorker[T any](jobs <-chan Job[T], out chan<- T) worker[T] {
	return worker[T]{jobs: jobs, out: out}
}

func (w worker[T]) start(ctx context.Context) {
	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				fmt.Println("worker finished")
				return
			}
			w.out <- job.Do()
		case <-ctx.Done():
			fmt.Println("worker finished because context was canceled")
			return
		}
	}
}

type Job[T any] interface {
	Do() T
}
