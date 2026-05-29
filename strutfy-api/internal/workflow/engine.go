package workflow

import (
	"context"
	"log"
	"sync"
)

type Engine struct {
	queue chan Job
	mu    sync.Mutex
	runs  map[string]string
}

func NewEngine(bufferSize int) *Engine {
	return &Engine{
		queue: make(chan Job, bufferSize),
		runs:  make(map[string]string),
	}
}

func (e *Engine) Dispatch(job Job) {
	e.queue <- job
}

func (e *Engine) Start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		go e.worker(ctx, i)
	}
}

func (e *Engine) worker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("workflow worker %d stopped", workerID)
			return

		case job := <-e.queue:
			e.setStatus(job.ID, "running")

			log.Printf("worker %d processing job %s type=%s", workerID, job.ID, job.Type)

			// Aqui entram:
			// Stripe webhook
			// Postmark
			// conciliação bancária
			// IA
			// relatórios

			e.setStatus(job.ID, "done")
		}
	}
}

func (e *Engine) setStatus(jobID string, status string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.runs[jobID] = status
}
