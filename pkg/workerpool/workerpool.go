package workerpool

import (
	"context"
	"errors"
)

// WorkerPool emplements worker pool pattern and allows to run limit amount of goroutines.
type WorkerPool struct {
	Pipeline chan func()
}

func New(ctx context.Context, chanCapacity int, numberOfWorkers int) (*WorkerPool, error) {
	if chanCapacity <= 0 {
		return nil, errors.New("chanCapacity must be grater than 0")
	}
	if numberOfWorkers <= 0 {
		return nil, errors.New("numberOfWorkers must be grater than 0")
	}

	wp := WorkerPool{
		Pipeline: make(chan func(), chanCapacity),
	}

	go wp.process(ctx, numberOfWorkers)
	return &wp, nil
}

func (wp *WorkerPool) process(ctx context.Context, numberOfWorkers int) {
	for range numberOfWorkers {
		go func() {
			select {
			case job := <-wp.Pipeline:
				if job != nil {
					job()
				}
			case <-ctx.Done():
				return
			}
		}()
	}

	<-ctx.Done()
}

func (wp *WorkerPool) AddJob(job func()) error {
	if job == nil {
		return errors.New("nil job")
	}

	select {
	case wp.Pipeline <- job:
		return nil
	default:
		return errors.New("channel is full")
	}
}
