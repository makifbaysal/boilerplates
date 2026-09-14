package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	grpcadapter "github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/grpc"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/grpc/taskpb"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/outbound/memory"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

// ServerSuite runs the real generated gRPC server (in-process, via
// bufconn — no real socket/port) against a real in-memory repository,
// the same "adapter test over the real use case" pattern as the REST
// adapter's HandlerSuite.
type ServerSuite struct {
	suite.Suite

	client taskpb.TaskServiceClient
	closer func()
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func (s *ServerSuite) SetupTest() {
	lis := bufconn.Listen(1024 * 1024)

	grpcServer := googlegrpc.NewServer()
	uc := task.NewService(memory.NewTaskRepository())
	taskpb.RegisterTaskServiceServer(grpcServer, grpcadapter.NewServer(uc))
	go func() { _ = grpcServer.Serve(lis) }()

	conn, err := googlegrpc.NewClient("passthrough:///bufnet",
		googlegrpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.client = taskpb.NewTaskServiceClient(conn)
	s.closer = func() {
		_ = conn.Close()
		grpcServer.Stop()
	}
}

func (s *ServerSuite) TearDownTest() {
	s.closer()
}

func (s *ServerSuite) TestCreate() {
	cases := []struct {
		name     string
		title    string
		wantCode codes.Code
	}{
		{name: "valid title", title: "buy milk", wantCode: codes.OK},
		{name: "blank title rejected", title: "   ", wantCode: codes.InvalidArgument},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.closer() // drop the suite-level connection from SetupTest, this case gets its own
			s.SetupTest()

			_, err := s.client.CreateTask(context.Background(), &taskpb.CreateTaskRequest{Title: tc.title})
			s.Equal(tc.wantCode, status.Code(err))
		})
	}
}

func (s *ServerSuite) TestCRUDFlow() {
	ctx := context.Background()

	created, err := s.client.CreateTask(ctx, &taskpb.CreateTaskRequest{Title: "buy milk"})
	s.Require().NoError(err)
	s.Require().NotEmpty(created.GetId())

	_, err = s.client.GetTask(ctx, &taskpb.GetTaskRequest{Id: created.GetId()})
	s.NoError(err)

	list, err := s.client.ListTasks(ctx, &taskpb.ListTasksRequest{})
	s.NoError(err)
	s.NotEmpty(list.GetTasks())

	done, err := s.client.SetTaskDone(ctx, &taskpb.SetTaskDoneRequest{Id: created.GetId(), Done: true})
	s.NoError(err)
	s.True(done.GetDone())

	_, err = s.client.DeleteTask(ctx, &taskpb.DeleteTaskRequest{Id: created.GetId()})
	s.NoError(err)

	_, err = s.client.GetTask(ctx, &taskpb.GetTaskRequest{Id: created.GetId()})
	s.Equal(codes.NotFound, status.Code(err))
}
