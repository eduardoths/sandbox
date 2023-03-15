package main

import (
	"context"
	"fmt"
	"time"

	workerpool "github.com/eduardoths/sandbox/go-utils/worker-pool"
)

type Job struct {
	id int
}

func newJob(id int) workerpool.Job[Response] {
	return Job{
		id: id,
	}
}

func (j Job) Do() Response {
	time.Sleep(time.Second)
	return Response{
		Message: fmt.Sprintf("job %d finished", j.id),
	}
}

type Response struct {
	Message string
}

func main() {
	jobs := make([]workerpool.Job[Response], 0)
	for i := 0; i < 100; i++ {
		jobs = append(jobs, newJob(i))
	}

	pool := workerpool.NewPool[Response](100)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	tsInit := time.Now()
	responses := pool.Start(ctx, jobs)
	fmt.Println("job execution finished after %u", time.Now().Sub(tsInit))
	for _, resp := range responses {
		fmt.Println(resp.Message)
	}
}
