package job_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job/mocks"
)

type ServiceSuite struct {
	suite.Suite

	store *mocks.Store
	svc   *job.Service
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) SetupTest() {
	s.store = mocks.NewStore(s.T())
	s.svc = job.NewService(s.store)
}

func (s *ServiceSuite) TestProcess() {
	cases := []struct {
		name      string
		taskID    string
		jobType   string
		setupMock func()
		wantErr   error
	}{
		{
			name:    "valid message is recorded",
			taskID:  "task-1",
			jobType: "task.created",
			setupMock: func() {
				s.store.EXPECT().
					Record(mock.Anything, mock.MatchedBy(func(j *job.Job) bool { return j.TaskID == "task-1" && j.Type == "task.created" })).
					Return(nil).
					Once()
			},
		},
		{
			name:      "missing task_id rejected without touching the store",
			taskID:    "  ",
			jobType:   "task.created",
			setupMock: func() {},
			wantErr:   job.ErrInvalidPayload,
		},
		{
			name:      "missing type rejected without touching the store",
			taskID:    "task-1",
			jobType:   "",
			setupMock: func() {},
			wantErr:   job.ErrInvalidPayload,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setupMock()

			got, err := s.svc.Process(context.Background(), tc.taskID, tc.jobType, "{}")

			if tc.wantErr != nil {
				s.ErrorIs(err, tc.wantErr)
				return
			}
			s.NoError(err)
			s.NotEmpty(got.ID)
			s.False(got.ProcessedAt.IsZero())
		})
	}
}
