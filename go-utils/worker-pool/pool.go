package workerpool

import (
	"context"
	"log"
	"sync"
)

type Pool[T any] struct {
	waitGroup  *sync.WaitGroup
	workersNum int
	jobChan    chan Job[T]
	outputChan chan T
	workers    []worker[T]
}

func NewPool[T any](workersNum int) Pool[T] {
	pool := Pool[T]{
		waitGroup:  &sync.WaitGroup{},
		workersNum: workersNum,
		jobChan:    make(chan Job[T], workersNum),
		outputChan: make(chan T, workersNum),
	}
	workers := make([]worker[T], workersNum)
	for i := 0; i < workersNum; i++ {
		workers[i] = newWorker(pool.waitGroup, pool.jobChan, pool.outputChan)
	}
	return pool
}

func (p Pool[T]) Start(ctx context.Context, jobs []Job[T]) []T {
	for i := range p.workers {
		go func(i int) {
			p.workers[i].start(ctx)
		}(i)
	}

	for i := range jobs {
		go func(i int) {
			p.jobChan <- jobs[i]
		}(i)
	}

	close(p.jobChan)
	p.waitGroup.Wait()

	outputs := make([]T, 0, len(jobs))
	for output := range p.outputChan {
		outputs = append(outputs, output)
	}
	return outputs
}

type worker[T any] struct {
	wg   *sync.WaitGroup
	jobs <-chan Job[T]
	out  chan<- T
}

func newWorker[T any](wg *sync.WaitGroup, jobs <-chan Job[T], out chan<- T) worker[T] {
	return worker[T]{wg: wg}
}

func (w worker[T]) start(ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("job finished because context was canceled")
			return
		case job, ok := <-w.jobs:
			if !ok {
				log.Println("job finished")
				return
			}
			w.out <- job.Do()
		}
	}
}

type Job[T any] interface {
	Do() T
}
