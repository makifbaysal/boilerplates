package graphql_test

import (
	"testing"

	gqlclient "github.com/99designs/gqlgen/client"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/suite"

	graphqladapter "github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/graphql"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/graphql/generated"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// ResolverSuite drives the real generated schema/resolvers over HTTP (via
// gqlgen's test client) against a real in-memory repository — the
// GraphQL sibling of rest.HandlerSuite and grpc's ServerSuite.
type ResolverSuite struct {
	suite.Suite

	c *gqlclient.Client
}

func TestResolverSuite(t *testing.T) {
	suite.Run(t, new(ResolverSuite))
}

func (s *ResolverSuite) SetupTest() {
	uc := task.NewService(memory.NewTaskRepository())
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: graphqladapter.NewResolver(uc)})
	s.c = gqlclient.New(gqlhandler.NewDefaultServer(schema))
}

func (s *ResolverSuite) TestCreateTask() {
	cases := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{name: "valid title", title: "buy milk"},
		{name: "blank title rejected", title: "   ", wantErr: true},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()

			var resp struct {
				CreateTask struct{ ID, Title string }
			}
			err := s.c.Post(`mutation($title: String!) { createTask(title: $title) { id title } }`, &resp,
				gqlclient.Var("title", tc.title))

			if tc.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.Equal("buy milk", resp.CreateTask.Title)
			s.NotEmpty(resp.CreateTask.ID)
		})
	}
}

func (s *ResolverSuite) TestCRUDFlow() {
	var created struct {
		CreateTask struct{ ID string }
	}
	s.Require().NoError(s.c.Post(
		`mutation($title: String!) { createTask(title: $title) { id } }`, &created,
		gqlclient.Var("title", "buy milk")))
	id := created.CreateTask.ID
	s.Require().NotEmpty(id)

	var listResp struct{ Tasks []struct{ ID string } }
	s.NoError(s.c.Post(`{ tasks { id } }`, &listResp))
	s.NotEmpty(listResp.Tasks)

	var doneResp struct {
		SetTaskDone struct{ Done bool }
	}
	s.NoError(s.c.Post(`mutation($id: ID!) { setTaskDone(id: $id, done: true) { done } }`, &doneResp, gqlclient.Var("id", id)))
	s.True(doneResp.SetTaskDone.Done)

	var deleteResp struct{ DeleteTask bool }
	s.NoError(s.c.Post(`mutation($id: ID!) { deleteTask(id: $id) }`, &deleteResp, gqlclient.Var("id", id)))
	s.True(deleteResp.DeleteTask)

	var getResp struct {
		Task *struct{ ID string }
	}
	s.NoError(s.c.Post(`query($id: ID!) { task(id: $id) { id } }`, &getResp, gqlclient.Var("id", id)))
	s.Nil(getResp.Task)
}
