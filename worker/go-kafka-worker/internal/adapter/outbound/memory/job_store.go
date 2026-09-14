// Package memory implements the job.Store port in-process. Default for
// cmd/{kafka,rabbitmq}, and what adapter tests run against.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*job.Job
}

func NewJobStore() *JobStore {
	return &JobStore{jobs: make(map[string]*job.Job)}
}

func (s *JobStore) Record(_ context.Context, j *job.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.ID] = j
	return nil
}

func (s *JobStore) Get(_ context.Context, id string) (*job.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil, job.ErrNotFound
	}
	return j, nil
}

func (s *JobStore) List(_ context.Context) ([]*job.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*job.Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		out = append(out, j)
	}
	sort.Slice(out, func(i, k int) bool { return out[i].ReceivedAt.Before(out[k].ReceivedAt) })
	return out, nil
}
