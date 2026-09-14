package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc;

import com.google.protobuf.Timestamp;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;

import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.CreateTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.DeleteTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.DeleteTaskResponse;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.GetTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.ListTasksRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.ListTasksResponse;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.SetTaskDoneRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.Task;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.TaskServiceGrpc;
import ai.tasktrooper.boilerplate.javaspring.core.task.InvalidTaskTitleException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskUseCase;

/** The gRPC sibling of adapter.inbound.rest.TaskController — same
 * TaskUseCase, different wire format. */
@GrpcService
public class TaskGrpcService extends TaskServiceGrpc.TaskServiceImplBase {

    private final TaskUseCase useCase;

    public TaskGrpcService(TaskUseCase useCase) {
        this.useCase = useCase;
    }

    private static Task toProto(ai.tasktrooper.boilerplate.javaspring.core.task.Task t) {
        return Task.newBuilder()
                .setId(t.id())
                .setTitle(t.title())
                .setDone(t.done())
                .setCreatedAt(toTimestamp(t.createdAt()))
                .setUpdatedAt(toTimestamp(t.updatedAt()))
                .build();
    }

    private static Timestamp toTimestamp(java.time.Instant instant) {
        return Timestamp.newBuilder().setSeconds(instant.getEpochSecond()).setNanos(instant.getNano()).build();
    }

    @Override
    public void createTask(CreateTaskRequest request, StreamObserver<Task> responseObserver) {
        try {
            responseObserver.onNext(toProto(useCase.create(request.getTitle())));
            responseObserver.onCompleted();
        } catch (RuntimeException ex) {
            responseObserver.onError(mapError(ex));
        }
    }

    @Override
    public void getTask(GetTaskRequest request, StreamObserver<Task> responseObserver) {
        try {
            responseObserver.onNext(toProto(useCase.get(request.getId())));
            responseObserver.onCompleted();
        } catch (RuntimeException ex) {
            responseObserver.onError(mapError(ex));
        }
    }

    @Override
    public void listTasks(ListTasksRequest request, StreamObserver<ListTasksResponse> responseObserver) {
        ListTasksResponse.Builder builder = ListTasksResponse.newBuilder();
        useCase.list().forEach(t -> builder.addTasks(toProto(t)));
        responseObserver.onNext(builder.build());
        responseObserver.onCompleted();
    }

    @Override
    public void setTaskDone(SetTaskDoneRequest request, StreamObserver<Task> responseObserver) {
        try {
            responseObserver.onNext(toProto(useCase.setDone(request.getId(), request.getDone())));
            responseObserver.onCompleted();
        } catch (RuntimeException ex) {
            responseObserver.onError(mapError(ex));
        }
    }

    @Override
    public void deleteTask(DeleteTaskRequest request, StreamObserver<DeleteTaskResponse> responseObserver) {
        try {
            useCase.delete(request.getId());
            responseObserver.onNext(DeleteTaskResponse.getDefaultInstance());
            responseObserver.onCompleted();
        } catch (RuntimeException ex) {
            responseObserver.onError(mapError(ex));
        }
    }

    private static Throwable mapError(RuntimeException ex) {
        if (ex instanceof TaskNotFoundException) {
            return Status.NOT_FOUND.withDescription(ex.getMessage()).asRuntimeException();
        }
        if (ex instanceof InvalidTaskTitleException) {
            return Status.INVALID_ARGUMENT.withDescription(ex.getMessage()).asRuntimeException();
        }
        return Status.INTERNAL.withDescription("internal error").asRuntimeException();
    }
}
