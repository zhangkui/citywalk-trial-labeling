package scheduler

import (
	"context"
	"sync"
	"time"
)

type Task func(context.Context)

type Scheduler struct {
	mu    sync.Mutex
	tasks map[string]context.CancelFunc
}

func New() *Scheduler {
	return &Scheduler{tasks: make(map[string]context.CancelFunc)}
}

func (s *Scheduler) Every(name string, interval time.Duration, task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cancel, ok := s.tasks[name]; ok {
		cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.tasks[name] = cancel
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				task(ctx)
			}
		}
	}()
}

func (s *Scheduler) Stop(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cancel, ok := s.tasks[name]; ok {
		cancel()
		delete(s.tasks, name)
	}
}

