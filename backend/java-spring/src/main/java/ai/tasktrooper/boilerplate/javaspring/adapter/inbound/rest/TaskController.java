package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest.dto.CreateTaskRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest.dto.SetDoneRequest;
import ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest.dto.TaskResponse;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskUseCase;

/**
 * Adapts HTTP to the TaskUseCase input port. Depends on the interface,
 * not TaskService, so it never needs to know how the use case is
 * implemented. See .ai/common-tasks.md for adding a new endpoint.
 */
@RestController
@RequestMapping("/tasks")
public class TaskController {

    private final TaskUseCase useCase;

    public TaskController(TaskUseCase useCase) {
        this.useCase = useCase;
    }

    @PostMapping
    public ResponseEntity<TaskResponse> create(@RequestBody CreateTaskRequest request) {
        return ResponseEntity.status(HttpStatus.CREATED).body(TaskResponse.from(useCase.create(request.title())));
    }

    @GetMapping
    public List<TaskResponse> list() {
        return useCase.list().stream().map(TaskResponse::from).toList();
    }

    @GetMapping("/{id}")
    public TaskResponse get(@PathVariable String id) {
        return TaskResponse.from(useCase.get(id));
    }

    @PatchMapping("/{id}")
    public TaskResponse setDone(@PathVariable String id, @RequestBody SetDoneRequest request) {
        return TaskResponse.from(useCase.setDone(id, request.done()));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable String id) {
        useCase.delete(id);
        return ResponseEntity.noContent().build();
    }
}
