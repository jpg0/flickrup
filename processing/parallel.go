package processing

import (
	"sync"
)

type Job struct {
	Index int
	Input interface{}
}

func ParallelMap(input []interface{}, mapper func(interface{}) interface{}, concurrency int) []interface{} {
	tasks := make(chan Job, 64)

	rv := make([]interface{}, len(input))

	// spawn N worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			for job := range tasks {
				rv[job.Index] = mapper(job.Input)
			}
			wg.Done()
		}()
	}

	for i, job := range input {
		tasks <- Job{i, job}
	}

	close(tasks)

	// wait for the workers to finish
	wg.Wait()

	return rv
}