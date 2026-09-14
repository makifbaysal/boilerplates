package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.graphql;

import java.util.List;

import org.springframework.graphql.data.method.annotation.Argument;
import org.springframework.graphql.data.method.annotation.MutationMapping;
import org.springframework.graphql.data.method.annotation.QueryMapping;
import org.springframework.stereotype.Controller;

import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.graphql.dto.TaskDto;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskUseCase;

/** The GraphQL sibling of adapter.inbound.rest.TaskController — same
 * TaskUseCase, different wire format. */
@Controller
public class TaskGraphQLController {

    private final TaskUseCase useCase;

    public TaskGraphQLController(TaskUseCase useCase) {
        this.useCase = useCase;
    }

    @QueryMapping
    public List<TaskDto> tasks() {
        return useCase.list().stream().map(TaskDto::from).toList();
    }

    @QueryMapping
    public TaskDto task(@Argument String id) {
        try {
            return TaskDto.from(useCase.get(id));
        } catch (TaskNotFoundException ex) {
            return null; // GraphQL convention: missing-by-id is null, not an error
        }
    }

    @MutationMapping
    public TaskDto createTask(@Argument String title) {
        return TaskDto.from(useCase.create(title));
    }

    @MutationMapping
    public TaskDto setTaskDone(@Argument String id, @Argument boolean done) {
        return TaskDto.from(useCase.setDone(id, done));
    }

    @MutationMapping
    public boolean deleteTask(@Argument String id) {
        useCase.delete(id);
        return true;
    }
}
