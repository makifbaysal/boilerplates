package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest;

import java.util.List;

import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.DELETE;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.PATCH;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;

import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest.dto.CreateTaskRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest.dto.SetDoneRequest;
import ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest.dto.TaskResponse;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskUseCase;

/**
 * Adapts JAX-RS to the TaskUseCase input port. Depends on the
 * interface, not TaskService, so it never needs to know how the use
 * case is implemented. See .ai/common-tasks.md for adding a new
 * endpoint.
 */
@Path("/tasks")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class TaskResource {

    private final TaskUseCase useCase;

    @Inject
    public TaskResource(TaskUseCase useCase) {
        this.useCase = useCase;
    }

    @POST
    public Response create(CreateTaskRequest request) {
        return Response.status(Response.Status.CREATED).entity(TaskResponse.from(useCase.create(request.title()))).build();
    }

    @GET
    public List<TaskResponse> list() {
        return useCase.list().stream().map(TaskResponse::from).toList();
    }

    @GET
    @Path("/{id}")
    public TaskResponse get(@PathParam("id") String id) {
        return TaskResponse.from(useCase.get(id));
    }

    @PATCH
    @Path("/{id}")
    public TaskResponse setDone(@PathParam("id") String id, SetDoneRequest request) {
        return TaskResponse.from(useCase.setDone(id, request.done()));
    }

    @DELETE
    @Path("/{id}")
    public Response delete(@PathParam("id") String id) {
        useCase.delete(id);
        return Response.noContent().build();
    }
}
