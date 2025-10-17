package workers

import (
	"context"
	"fmt"
	errWrap "task_queue/common/error"
	"task_queue/services"
	"time"
)

type WorkerImpl struct {
	queueService services.QueueService
}

func NewWorker(queueService services.QueueService) Worker {
	return &WorkerImpl{queueService: queueService}
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
					data, err := w.queueService.GetQueue(ctx)
					if err != nil {
						time.Sleep(50 * time.Millisecond)
						continue
					}
					if data == nil {
						time.Sleep(50 * time.Millisecond)
						continue
					}

					err = w.queueService.PublishPredictionToS3AndMQTT(ctx, data)
					if err != nil {
						time.Sleep(50 * time.Millisecond)
						continue
					}

					time.Sleep(30 * time.Millisecond)
				}
			}
		}(i + 1)
	}
	return nil
}
