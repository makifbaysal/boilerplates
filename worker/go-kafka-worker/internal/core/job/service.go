package job

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

var _ Processor = (*Service)(nil)

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Process(ctx context.Context, taskID, jobType, payload string) (*Job, error) {
	taskID = strings.TrimSpace(taskID)
	jobType = strings.TrimSpace(jobType)
	if taskID == "" || jobType == "" {
		return nil, ErrInvalidPayload
	}

	now := time.Now().UTC()
	j := &Job{
		ID:          uuid.NewString(),
		TaskID:      taskID,
		Type:        jobType,
		Payload:     payload,
		ReceivedAt:  now,
		ProcessedAt: now,
	}
	if err := s.store.Record(ctx, j); err != nil {
		return nil, err
	}
	return j, nil
}
