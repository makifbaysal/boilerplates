package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.graphql;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.graphql.GraphQlTest;
import org.springframework.context.annotation.Import;
import org.springframework.graphql.execution.ErrorType;
import org.springframework.graphql.test.tester.GraphQlTester;

import ai.tasktrooper.boilerplate.javaspring.adapter.outbound.memory.InMemoryTaskRepository;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskService;

/**
 * Runs the GraphQL adapter against a real InMemoryTaskRepository — the
 * GraphQL sibling of TaskControllerTest / TaskGrpcServiceTest.
 */
@GraphQlTest(TaskGraphQLController.class)
@Import({TaskService.class, InMemoryTaskRepository.class, GraphQLExceptionResolver.class})
class TaskGraphQLControllerTest {

    @Autowired
    private GraphQlTester graphQlTester;

    @ParameterizedTest(name = "title=\"{0}\"")
    @CsvSource({"buy milk,true", "'   ',false"})
    void createTask(String title, boolean wantOk) {
        GraphQlTester.Response response = graphQlTester.document(
                        "mutation($title: String!) { createTask(title: $title) { id title } }")
                .variable("title", title)
                .execute();

        if (wantOk) {
            response.path("createTask.title").entity(String.class).isEqualTo("buy milk");
        } else {
            response.errors().expect(err -> err.getErrorType() == ErrorType.BAD_REQUEST);
        }
    }

    @Test
    void crudFlow() {
        String id = graphQlTester.document("mutation($title: String!) { createTask(title: $title) { id } }")
                .variable("title", "buy milk")
                .execute()
                .path("createTask.id")
                .entity(String.class)
                .get();
        assertThat(id).isNotBlank();

        graphQlTester.document("{ tasks { id } }").execute().path("tasks").entityList(Object.class).hasSizeGreaterThan(0);

        graphQlTester.document("mutation($id: ID!) { setTaskDone(id: $id, done: true) { done } }")
                .variable("id", id)
                .execute()
                .path("setTaskDone.done")
                .entity(Boolean.class)
                .isEqualTo(true);

        graphQlTester.document("mutation($id: ID!) { deleteTask(id: $id) }")
                .variable("id", id)
                .execute()
                .path("deleteTask")
                .entity(Boolean.class)
                .isEqualTo(true);

        graphQlTester.document("query($id: ID!) { task(id: $id) { id } }")
                .variable("id", id)
                .execute()
                .path("task")
                .valueIsNull();
    }
}
