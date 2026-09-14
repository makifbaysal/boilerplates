package task_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task/mocks"
)

// ServiceSuite is the standard shape for this repo: a testify suite per
// unit under test, one mocked Repository per test (fresh via SetupTest),
// table-driven cases inside each Test* method. See .ai/testing.md.
type ServiceSuite struct {
	suite.Suite

	repo *mocks.Repository
	svc  *task.Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.repo = mocks.NewRepository(s.T())
	s.svc = task.NewService(s.repo)
}

func (s *ServiceSuite) TestCreate() {
	cases := []struct {
		name      string
		title     string
		setupMock func()
		wantErr   error
		wantTitle string
	}{
		{
			name:  "trims and persists a valid title",
			title: "  buy milk  ",
			setupMock: func() {
				s.repo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(t *task.Task) bool { return t.Title == "buy milk" })).
					Return(nil).
					Once()
			},
			wantTitle: "buy milk",
		},
		{
			name:      "rejects an empty title without touching the repository",
			title:     "   ",
			setupMock: func() {},
			wantErr:   task.ErrInvalidTitle,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setupMock()

			got, err := s.svc.Create(context.Background(), tc.title)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				return
			}
			s.NoError(err)
			s.Equal(tc.wantTitle, got.Title)
			s.False(got.Done)
			s.NotEmpty(got.ID)
		})
	}
}

func (s *ServiceSuite) TestSetDone() {
	cases := []struct {
		name      string
		id        string
		done      bool
		setupMock func()
		wantErr   error
	}{
		{
			name: "marks an existing task done",
			id:   "abc",
			done: true,
			setupMock: func() {
				existing := &task.Task{ID: "abc", Title: "buy milk"}
				s.repo.EXPECT().Get(mock.Anything, "abc").Return(existing, nil).Once()
				s.repo.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(t *task.Task) bool { return t.Done })).
					Return(nil).
					Once()
			},
		},
		{
			name: "propagates not-found from the repository",
			id:   "missing",
			done: true,
			setupMock: func() {
				s.repo.EXPECT().Get(mock.Anything, "missing").Return(nil, task.ErrNotFound).Once()
			},
			wantErr: task.ErrNotFound,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setupMock()

			got, err := s.svc.SetDone(context.Background(), tc.id, tc.done)

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				return
			}
			s.NoError(err)
			s.True(got.Done)
		})
	}
}
