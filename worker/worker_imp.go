package worker

import (
	"context"
	"fmt"
	"task_queue/repositories"

	errWrap "task_queue/common/error"
	"time"
)

type WorkerImpl struct {
	repository repositories.QueueRepository
}

func NewWorker(repository repositories.QueueRepository) Worker {
	return &WorkerImpl{repository: repository}
}

func (w *WorkerImpl) Run(ctx context.Context, worker int) error {
	for i := 0; i < worker; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					errWrap.WrapError(fmt.Errorf("worker %d stopped", workerID))
					return
				default:
					data, err := w.repository.GetQueue(ctx)
					if err != nil {
						errWrap.WrapError(err)
						time.Sleep(1 * time.Second)
						continue
					}
					if data == nil {
						time.Sleep(1 * time.Second)
						continue
					}
					fmt.Printf("Worker %d - data: %+v\n", workerID, data)
				}
			}
		}(i + 1)
	}
	return nil
}
