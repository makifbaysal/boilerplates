package resilience

import (
	"sync"
	"time"
)

// Limiter is a token bucket: Rate tokens are added per second up to Burst.
// It is what keeps a consumer (or any non-HTTP entry point) from pulling more
// work than the downstream dependencies can absorb — the inbound counterpart
// of the breaker.
type Limiter struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	tokens float64
	last   time.Time
	now    func() time.Time
}

// NewLimiter builds a limiter allowing `rate` events per second with a burst
// of `burst`. A non-positive rate disables limiting (Allow always true), so a
// zero-valued config never accidentally stops all work.
func NewLimiter(rate, burst float64) *Limiter {
	if burst < 1 && rate > 0 {
		burst = rate
	}
	return &Limiter{rate: rate, burst: burst, tokens: burst, now: time.Now}
}

// Allow consumes one token if available.
func (l *Limiter) Allow() bool {
	if l == nil || l.rate <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if l.last.IsZero() {
		l.last = now
	}
	l.tokens += now.Sub(l.last).Seconds() * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}
	l.last = now
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Wait blocks until a token is available or the deadline passes, returning
// whether it got one. Use it in a consumer loop where dropping the message is
// not an option — unlike an HTTP request, a queued job can simply wait.
func (l *Limiter) Wait(timeout time.Duration) bool {
	if l == nil || l.rate <= 0 {
		return true
	}
	deadline := time.Now().Add(timeout)
	for {
		if l.Allow() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(time.Duration(float64(time.Second) / l.rate / 4))
	}
}
