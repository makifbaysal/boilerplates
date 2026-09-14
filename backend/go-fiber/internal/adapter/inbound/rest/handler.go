package rest

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// Handler adapts HTTP to the task.UseCase input port. It depends on the
// interface, not task.Service, so it never needs to know how the use
// case is implemented.
type Handler struct {
	uc task.UseCase
}

func NewHandler(uc task.UseCase) *Handler {
	return &Handler{uc: uc}
}

// Register wires the task routes onto router. See .ai/common-tasks.md for
// the walkthrough on adding a new endpoint.
func (h *Handler) Register(router fiber.Router) {
	router.Get("/tasks", h.list)
	router.Post("/tasks", h.create)
	router.Get("/tasks/:id", h.get)
	router.Patch("/tasks/:id", h.setDone)
	router.Delete("/tasks/:id", h.delete)
}

func (h *Handler) create(c fiber.Ctx) error {
	var req createRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	t, err := h.uc.Create(c.Context(), req.Title)
	if err != nil {
		return mapError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(toResponse(t))
}

func (h *Handler) list(c fiber.Ctx) error {
	tasks, err := h.uc.List(c.Context())
	if err != nil {
		return mapError(err)
	}
	return c.JSON(toResponseList(tasks))
}

func (h *Handler) get(c fiber.Ctx) error {
	t, err := h.uc.Get(c.Context(), c.Params("id"))
	if err != nil {
		return mapError(err)
	}
	return c.JSON(toResponse(t))
}

func (h *Handler) setDone(c fiber.Ctx) error {
	var req setDoneRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	t, err := h.uc.SetDone(c.Context(), c.Params("id"), req.Done)
	if err != nil {
		return mapError(err)
	}
	return c.JSON(toResponse(t))
}

func (h *Handler) delete(c fiber.Ctx) error {
	if err := h.uc.Delete(c.Context(), c.Params("id")); err != nil {
		return mapError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func mapError(err error) error {
	switch {
	case errors.Is(err, task.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, task.ErrInvalidTitle):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}
}
