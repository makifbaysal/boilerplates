// Package grpc adapts google.golang.org/grpc to the task.UseCase input
// port — the gRPC sibling of adapter/inbound/rest. Regenerate taskpb
// from proto/task.proto with `make proto` after editing it.
package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/adapter/inbound/grpc/taskpb"
	"github.com/makifbaysal/boilerplates/backend/go-fiber/internal/core/task"
)

type Server struct {
	taskpb.UnimplementedTaskServiceServer
	uc task.UseCase
}

func NewServer(uc task.UseCase) *Server {
	return &Server{uc: uc}
}

func toProto(t *task.Task) *taskpb.Task {
	return &taskpb.Task{
		Id:        t.ID,
		Title:     t.Title,
		Done:      t.Done,
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}

func (s *Server) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.Task, error) {
	t, err := s.uc.Create(ctx, req.GetTitle())
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(t), nil
}

func (s *Server) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	t, err := s.uc.Get(ctx, req.GetId())
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(t), nil
}

func (s *Server) ListTasks(ctx context.Context, _ *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := s.uc.List(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	out := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, toProto(t))
	}
	return &taskpb.ListTasksResponse{Tasks: out}, nil
}

func (s *Server) SetTaskDone(ctx context.Context, req *taskpb.SetTaskDoneRequest) (*taskpb.Task, error) {
	t, err := s.uc.SetDone(ctx, req.GetId(), req.GetDone())
	if err != nil {
		return nil, mapError(err)
	}
	return toProto(t), nil
}

func (s *Server) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := s.uc.Delete(ctx, req.GetId()); err != nil {
		return nil, mapError(err)
	}
	return &taskpb.DeleteTaskResponse{}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, task.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, task.ErrInvalidTitle):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
