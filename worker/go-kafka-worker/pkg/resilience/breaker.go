// Package resilience holds the two defaults every service that calls
// something else needs: a circuit breaker for outbound dependencies and a
// token-bucket rate limiter for inbound work. Both are dependency-free on
// purpose — a boilerplate should not force a resilience library on the
// project that copies it.
package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned instead of calling the dependency while the breaker is
// open. Fail fast: a dependency that is already down does not get a thundering
// herd of doomed requests on top of its outage.
var ErrOpen = errors.New("circuit breaker is open")

// State is the breaker's position in the closed → open → half-open cycle.
type State int

const (
	// StateClosed passes every call through and counts failures.
	StateClosed State = iota
	// StateOpen rejects immediately until the cooldown expires.
	StateOpen
	// StateHalfOpen lets a limited number of probes through; one failure
	// re-opens, enough successes close it again.
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "closed"
	}
}

// BreakerConfig tunes the breaker. Zero values fall back to the defaults, so
// `NewBreaker(BreakerConfig{})` is a usable production setting.
type BreakerConfig struct {
	// FailureThreshold is how many consecutive failures open the circuit.
	FailureThreshold int
	// OpenDuration is how long the circuit stays open before probing.
	OpenDuration time.Duration
	// HalfOpenMaxCalls is how many probes are allowed while half-open.
	HalfOpenMaxCalls int
	// SuccessThreshold is how many probe successes close the circuit.
	SuccessThreshold int
	// Now is injectable for tests; defaults to time.Now.
	Now func() time.Time
}

func (c BreakerConfig) withDefaults() BreakerConfig {
	if c.FailureThreshold <= 0 {
		c.FailureThreshold = 5
	}
	if c.OpenDuration <= 0 {
		c.OpenDuration = 30 * time.Second
	}
	if c.HalfOpenMaxCalls <= 0 {
		c.HalfOpenMaxCalls = 1
	}
	if c.SuccessThreshold <= 0 {
		c.SuccessThreshold = 1
	}
	if c.Now == nil {
		c.Now = time.Now
	}
	return c
}

// Breaker guards one outbound dependency (a database, an HTTP API, a broker).
// Use one breaker per dependency: sharing a breaker between two dependencies
// makes a healthy one unavailable because the other is broken.
type Breaker struct {
	cfg BreakerConfig

	mu        sync.Mutex
	state     State
	failures  int
	successes int
	halfOpen  int
	openedAt  time.Time
}

func NewBreaker(cfg BreakerConfig) *Breaker {
	return &Breaker{cfg: cfg.withDefaults()}
}

// State reports the current state (for health endpoints and metrics).
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refresh()
	return b.state
}

// Do runs fn unless the circuit is open. The context error is respected
// first: a caller that already gave up must not be counted as a dependency
// failure.
func (b *Breaker) Do(ctx context.Context, fn func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := b.allow(); err != nil {
		return err
	}
	err := fn(ctx)
	// A cancelled/expired caller context says nothing about the dependency.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		b.release()
		return err
	}
	if err != nil {
		b.onFailure()
		return err
	}
	b.onSuccess()
	return nil
}

// allow decides whether a call may proceed, moving open → half-open when the
// cooldown has elapsed.
func (b *Breaker) allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refresh()
	switch b.state {
	case StateOpen:
		return ErrOpen
	case StateHalfOpen:
		if b.halfOpen >= b.cfg.HalfOpenMaxCalls {
			return ErrOpen
		}
		b.halfOpen++
	}
	return nil
}

// release gives back a half-open slot taken by a call that told us nothing.
func (b *Breaker) release() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateHalfOpen && b.halfOpen > 0 {
		b.halfOpen--
	}
}

func (b *Breaker) refresh() {
	if b.state == StateOpen && b.cfg.Now().Sub(b.openedAt) >= b.cfg.OpenDuration {
		b.state = StateHalfOpen
		b.halfOpen = 0
		b.successes = 0
	}
}

func (b *Breaker) onSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateHalfOpen:
		b.successes++
		if b.halfOpen > 0 {
			b.halfOpen--
		}
		if b.successes >= b.cfg.SuccessThreshold {
			b.state = StateClosed
			b.failures = 0
			b.successes = 0
		}
	default:
		b.failures = 0
	}
}

func (b *Breaker) onFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateHalfOpen:
		// The probe failed: the dependency is still down, back to open.
		b.state = StateOpen
		b.openedAt = b.cfg.Now()
		b.failures = b.cfg.FailureThreshold
		b.successes = 0
		b.halfOpen = 0
	default:
		b.failures++
		if b.failures >= b.cfg.FailureThreshold {
			b.state = StateOpen
			b.openedAt = b.cfg.Now()
		}
	}
}
