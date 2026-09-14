package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.graphql;

import java.util.List;

import jakarta.inject.Inject;

import org.eclipse.microprofile.graphql.GraphQLApi;
import org.eclipse.microprofile.graphql.Mutation;
import org.eclipse.microprofile.graphql.Name;
import org.eclipse.microprofile.graphql.Query;

import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.graphql.dto.TaskDto;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskUseCase;

/**
 * Code-first SmallRye GraphQL API — the GraphQL sibling of
 * adapter.inbound.rest.TaskResource, same TaskUseCase. Quarkus'
 * MicroProfile GraphQL is code-first (annotations generate the schema),
 * unlike go-fiber/java-spring's schema-first approach — there's no
 * schema.graphqls to hand-write here; run with
 * `quarkus.smallrye-graphql.print-data-fetcher-exception=true` in dev
 * to see the generated schema at /q/graphql-schema.
 */
@GraphQLApi
public class TaskGraphQLApi {

    private final TaskUseCase useCase;

    @Inject
    public TaskGraphQLApi(TaskUseCase useCase) {
        this.useCase = useCase;
    }

    @Query
    public List<TaskDto> tasks() {
        return useCase.list().stream().map(TaskDto::from).toList();
    }

    @Query
    public TaskDto task(@Name("id") String id) {
        try {
            return TaskDto.from(useCase.get(id));
        } catch (TaskNotFoundException ex) {
            return null; // GraphQL convention: missing-by-id is null, not an error
        }
    }

    @Mutation
    public TaskDto createTask(@Name("title") String title) {
        return TaskDto.from(useCase.create(title));
    }

    // Explicit @Name: MicroProfile GraphQL strips a leading "set" from
    // mutation method names as if it were a bean setter (setTaskDone ->
    // field "taskDone"), which silently produced a different field name
    // than intended — caught by TaskGraphQLApiTest.crudFlow failing with
    // "Field 'setTaskDone' ... is undefined".
    @Mutation
    @Name("setTaskDone")
    public TaskDto setTaskDone(@Name("id") String id, @Name("done") boolean done) {
        return TaskDto.from(useCase.setDone(id, done));
    }

    @Mutation
    public boolean deleteTask(@Name("id") String id) {
        useCase.delete(id);
        return true;
    }
}
