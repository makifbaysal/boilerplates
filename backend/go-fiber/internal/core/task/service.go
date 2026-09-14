package task

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service implements UseCase against a Repository. It's the only place
// business rules live — inbound adapters translate protocol <-> these
// calls, outbound adapters translate these calls <-> storage.
type Service struct {
	repo Repository
}

var _ UseCase = (*Service)(nil)

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, title string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidTitle
	}
	now := time.Now().UTC()
	t := &Task{
		ID:        uuid.NewString(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) SetDone(ctx context.Context, id string, done bool) (*Task, error) {
	t, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Done = done
	t.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
