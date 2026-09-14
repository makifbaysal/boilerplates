package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import io.grpc.Server;
import io.grpc.StatusRuntimeException;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.ManagedChannel;

import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.CreateTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.DeleteTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.GetTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.SetTaskDoneRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.Task;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.TaskServiceGrpc;
import ai.tasktrooper.boilerplate.javaspring.adapter.outbound.memory.InMemoryTaskRepository;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskService;

/**
 * Runs the real generated gRPC service in-process (io.grpc.inprocess —
 * no real socket/port) against a real InMemoryTaskRepository — the gRPC
 * sibling of TaskControllerTest. See .ai/testing.md.
 */
class TaskGrpcServiceTest {

    private Server server;
    private ManagedChannel channel;
    private TaskServiceGrpc.TaskServiceBlockingStub client;

    @BeforeEach
    void setUp() throws IOException {
        String name = "task-grpc-test-" + System.nanoTime();
        TaskService useCase = new TaskService(new InMemoryTaskRepository());

        server = InProcessServerBuilder.forName(name)
                .directExecutor()
                .addService(new TaskGrpcService(useCase))
                .build()
                .start();
        channel = InProcessChannelBuilder.forName(name).directExecutor().build();
        client = TaskServiceGrpc.newBlockingStub(channel);
    }

    @AfterEach
    void tearDown() throws InterruptedException {
        channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
        server.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
    }

    @Test
    void createRejectsBlankTitle() {
        assertThatThrownBy(() -> client.createTask(CreateTaskRequest.newBuilder().setTitle("   ").build()))
                .isInstanceOf(StatusRuntimeException.class)
                .hasMessageContaining("INVALID_ARGUMENT");
    }

    @Test
    void crudFlow() {
        Task created = client.createTask(CreateTaskRequest.newBuilder().setTitle("buy milk").build());
        assertThat(created.getId()).isNotBlank();

        client.getTask(GetTaskRequest.newBuilder().setId(created.getId()).build());
        assertThat(client.listTasks(ai.tasktrooper.boilerplate.javaspring.adapter.inbound.grpc.taskpb.ListTasksRequest.getDefaultInstance())
                .getTasksList()).isNotEmpty();

        Task done = client.setTaskDone(SetTaskDoneRequest.newBuilder().setId(created.getId()).setDone(true).build());
        assertThat(done.getDone()).isTrue();

        client.deleteTask(DeleteTaskRequest.newBuilder().setId(created.getId()).build());

        assertThatThrownBy(() -> client.getTask(GetTaskRequest.newBuilder().setId(created.getId()).build()))
                .isInstanceOf(StatusRuntimeException.class)
                .hasMessageContaining("NOT_FOUND");
    }
}
