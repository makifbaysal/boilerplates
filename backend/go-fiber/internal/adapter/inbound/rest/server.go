package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// RateLimit is this adapter's own view of the inbound rate limit, built by
// cmd/api from platform/config — the same trick as the logger middleware, so
// this package keeps no dependency on the config package.
//
// It is on by default: an API without a rate limit is one loop away from being
// taken down by a single misbehaving client.
type RateLimit struct {
	// Disabled turns the limiter off entirely (local debugging, load tests).
	Disabled bool
	// Max requests per Window, per client IP.
	Max int
	// Window is the sliding window length.
	Window time.Duration
}

func (r RateLimit) withDefaults() RateLimit {
	if r.Max <= 0 {
		r.Max = 120
	}
	if r.Window <= 0 {
		r.Window = time.Minute
	}
	return r
}

// New builds the fiber app for the REST inbound adapter. extraMiddleware
// (e.g. platform/logger's fiber middleware) is installed first, ahead of
// recover/cors, so request logging wraps everything below it — passed in
// from cmd/api/main.go so this package stays free of a direct
// platform/logger dependency.
func New(uc task.UseCase, extraMiddleware ...fiber.Handler) *fiber.App {
	return NewWithRateLimit(uc, RateLimit{}, extraMiddleware...)
}

// NewWithRateLimit is New with an explicit inbound rate limit.
func NewWithRateLimit(uc task.UseCase, rl RateLimit, extraMiddleware ...fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "go-fiber-boilerplate",
	})

	for _, mw := range extraMiddleware {
		app.Use(mw)
	}
	app.Use(recover.New())
	app.Use(cors.New())

	if rl = rl.withDefaults(); !rl.Disabled {
		app.Use(limiter.New(limiter.Config{
			Max:        rl.Max,
			Expiration: rl.Window,
			// Health checks come from the platform (k8s, load balancer) and
			// must never be throttled — a rate-limited probe reads as an
			// outage and triggers a pointless restart.
			Next: func(c fiber.Ctx) bool {
				return c.Path() == "/healthz"
			},
			LimitReached: func(c fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).
					JSON(fiber.Map{"error": "rate limit exceeded"})
			},
		}))
	}

	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	NewHandler(uc).Register(app)

	return app
}
