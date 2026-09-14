package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import org.junit.jupiter.api.Test;

import io.grpc.StatusRuntimeException;
import io.quarkus.grpc.GrpcClient;
import io.quarkus.test.junit.QuarkusTest;

import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.CreateTaskRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.DeleteTaskRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.GetTaskRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.ListTasksRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.SetTaskDoneRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.Task;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.grpc.taskpb.TaskServiceGrpc;

/**
 * Runs the real gRPC service that @QuarkusTest starts (in-memory
 * TaskUseCase, since no app.datastore is set) — the gRPC sibling of
 * TaskResourceTest.
 */
@QuarkusTest
class TaskGrpcServiceTest {

    @GrpcClient
    TaskServiceGrpc.TaskServiceBlockingStub taskService;

    @Test
    void createRejectsBlankTitle() {
        assertThatThrownBy(() -> taskService.createTask(CreateTaskRequest.newBuilder().setTitle("   ").build()))
                .isInstanceOf(StatusRuntimeException.class)
                .hasMessageContaining("INVALID_ARGUMENT");
    }

    @Test
    void crudFlow() {
        Task created = taskService.createTask(CreateTaskRequest.newBuilder().setTitle("buy milk").build());
        assertThat(created.getId()).isNotBlank();

        taskService.getTask(GetTaskRequest.newBuilder().setId(created.getId()).build());
        assertThat(taskService.listTasks(ListTasksRequest.getDefaultInstance()).getTasksList()).isNotEmpty();

        Task done = taskService.setTaskDone(SetTaskDoneRequest.newBuilder().setId(created.getId()).setDone(true).build());
        assertThat(done.getDone()).isTrue();

        taskService.deleteTask(DeleteTaskRequest.newBuilder().setId(created.getId()).build());

        assertThatThrownBy(() -> taskService.getTask(GetTaskRequest.newBuilder().setId(created.getId()).build()))
                .isInstanceOf(StatusRuntimeException.class)
                .hasMessageContaining("NOT_FOUND");
    }
}
