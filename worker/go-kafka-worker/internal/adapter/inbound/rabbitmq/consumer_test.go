package rabbitmq_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/inbound/rabbitmq"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/worker/go-kafka-worker/internal/core/job"
)

// ConsumerSuite exercises HandleMessage against a real in-memory Store —
// no broker needed. Run needs a live RabbitMQ channel, so it's covered
// by the skip-gated integration test instead (see .ai/testing.md).
type ConsumerSuite struct {
	suite.Suite

	store *memory.JobStore
	c     *rabbitmq.Consumer
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}

func (s *ConsumerSuite) SetupTest() {
	s.store = memory.NewJobStore()
	uc := job.NewService(s.store)
	s.c = rabbitmq.NewConsumer(uc, slog.Default())
}

func (s *ConsumerSuite) TestHandleMessage() {
	cases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid message is processed and stored", value: `{"task_id":"task-1","type":"task.created"}`},
		{name: "invalid json rejected", value: `not-json`, wantErr: true},
		{name: "missing task_id rejected by the use case", value: `{"type":"task.created"}`, wantErr: true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()

			err := s.c.HandleMessage(context.Background(), "task-events", []byte(tc.value))

			if tc.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)

			list, err := s.store.List(context.Background())
			s.NoError(err)
			s.Len(list, 1)
			s.Equal("task-1", list[0].TaskID)
		})
	}
}
