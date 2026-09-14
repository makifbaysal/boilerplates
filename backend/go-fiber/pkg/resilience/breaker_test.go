package resilience_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/pkg/resilience"
)

type BreakerSuite struct {
	suite.Suite
	now time.Time
}

func TestBreakerSuite(t *testing.T) { suite.Run(t, new(BreakerSuite)) }

func (s *BreakerSuite) SetupTest() {
	s.now = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
}

func (s *BreakerSuite) breaker() *resilience.Breaker {
	return resilience.NewBreaker(resilience.BreakerConfig{
		FailureThreshold: 2,
		OpenDuration:     10 * time.Second,
		HalfOpenMaxCalls: 1,
		SuccessThreshold: 1,
		Now:              func() time.Time { return s.now },
	})
}

var errBoom = errors.New("boom")

func fail(context.Context) error { return errBoom }
func ok(context.Context) error   { return nil }

func (s *BreakerSuite) TestOpensAfterThresholdAndFailsFast() {
	b := s.breaker()
	ctx := context.Background()

	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.Equal(resilience.StateClosed, b.State())
	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.Equal(resilience.StateOpen, b.State())

	called := false
	err := b.Do(ctx, func(context.Context) error { called = true; return nil })
	s.ErrorIs(err, resilience.ErrOpen)
	s.False(called, "an open breaker must not reach the dependency")
}

func (s *BreakerSuite) TestSuccessResetsFailureCount() {
	b := s.breaker()
	ctx := context.Background()

	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.NoError(b.Do(ctx, ok))
	s.ErrorIs(b.Do(ctx, fail), errBoom)

	s.Equal(resilience.StateClosed, b.State(), "failures must be consecutive to open the circuit")
}

func (s *BreakerSuite) TestHalfOpenProbeClosesOnSuccess() {
	b := s.breaker()
	ctx := context.Background()
	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.ErrorIs(b.Do(ctx, fail), errBoom)

	s.now = s.now.Add(11 * time.Second)
	s.Equal(resilience.StateHalfOpen, b.State())

	s.NoError(b.Do(ctx, ok))
	s.Equal(resilience.StateClosed, b.State())
}

func (s *BreakerSuite) TestHalfOpenProbeFailureReopens() {
	b := s.breaker()
	ctx := context.Background()
	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.ErrorIs(b.Do(ctx, fail), errBoom)
	s.now = s.now.Add(11 * time.Second)

	s.ErrorIs(b.Do(ctx, fail), errBoom)

	s.Equal(resilience.StateOpen, b.State())
	s.ErrorIs(b.Do(ctx, ok), resilience.ErrOpen)
}

// A caller that gave up says nothing about the dependency's health.
func (s *BreakerSuite) TestCancelledContextIsNotADependencyFailure() {
	b := s.breaker()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s.ErrorIs(b.Do(ctx, ok), context.Canceled)
	s.ErrorIs(b.Do(context.Background(), func(context.Context) error { return context.DeadlineExceeded }), context.DeadlineExceeded)

	s.Equal(resilience.StateClosed, b.State())
}

func (s *BreakerSuite) TestLimiterAllowsBurstThenThrottles() {
	l := resilience.NewLimiter(10, 2)

	s.True(l.Allow())
	s.True(l.Allow())
	s.False(l.Allow(), "burst is exhausted")
}

func (s *BreakerSuite) TestLimiterDisabledWhenRateIsZero() {
	l := resilience.NewLimiter(0, 0)
	for i := 0; i < 100; i++ {
		s.True(l.Allow())
	}
}
