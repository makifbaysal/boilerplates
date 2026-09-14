package rest

import (
	"time"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// taskResponse is the wire shape for this adapter. Never expose
// core.Task directly — the core has no idea REST/JSON exists, and this
// struct is where you add transport-specific fields (HATEOAS links,
// api version, etc.) without touching the domain.
type taskResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(t *task.Task) taskResponse {
	return taskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Done:      t.Done,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toResponseList(tasks []*task.Task) []taskResponse {
	out := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, toResponse(t))
	}
	return out
}

type createRequest struct {
	Title string `json:"title"`
}

type setDoneRequest struct {
	Done bool `json:"done"`
}
