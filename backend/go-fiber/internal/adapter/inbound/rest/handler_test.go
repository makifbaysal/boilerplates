package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/rest"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// HandlerSuite runs the REST adapter against a real in-memory repository
// (no mocks here — the use case is trivial enough that an end-to-end
// adapter test is cheaper and more honest than mocking task.UseCase).
type HandlerSuite struct {
	suite.Suite

	app *fiber.App
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) SetupTest() {
	uc := task.NewService(memory.NewTaskRepository())
	s.app = rest.New(uc)
}

func (s *HandlerSuite) do(method, path string, body any) *http.Response {
	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		s.Require().NoError(err)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)
	s.Require().NoError(err)
	return resp
}

func (s *HandlerSuite) TestCreate() {
	cases := []struct {
		name       string
		title      string
		wantStatus int
	}{
		{name: "valid title", title: "buy milk", wantStatus: http.StatusCreated},
		{name: "blank title rejected", title: "   ", wantStatus: http.StatusBadRequest},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			resp := s.do(http.MethodPost, "/tasks", map[string]string{"title": tc.title})
			s.Equal(tc.wantStatus, resp.StatusCode)
		})
	}
}

func (s *HandlerSuite) TestCRUDFlow() {
	created := s.do(http.MethodPost, "/tasks", map[string]string{"title": "buy milk"})
	s.Require().Equal(http.StatusCreated, created.StatusCode)

	var body struct {
		ID string `json:"id"`
	}
	s.Require().NoError(json.NewDecoder(created.Body).Decode(&body))
	s.Require().NotEmpty(body.ID)

	s.Equal(http.StatusOK, s.do(http.MethodGet, "/tasks/"+body.ID, nil).StatusCode)
	s.Equal(http.StatusOK, s.do(http.MethodGet, "/tasks", nil).StatusCode)
	s.Equal(http.StatusOK, s.do(http.MethodPatch, "/tasks/"+body.ID, map[string]bool{"done": true}).StatusCode)
	s.Equal(http.StatusNoContent, s.do(http.MethodDelete, "/tasks/"+body.ID, nil).StatusCode)
	s.Equal(http.StatusNotFound, s.do(http.MethodGet, "/tasks/"+body.ID, nil).StatusCode)
}
