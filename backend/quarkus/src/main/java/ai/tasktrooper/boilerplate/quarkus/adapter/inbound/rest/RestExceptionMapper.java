package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest;

import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;

import ai.tasktrooper.boilerplate.quarkus.core.task.InvalidTaskTitleException;
import ai.tasktrooper.boilerplate.quarkus.core.task.TaskNotFoundException;

/** Maps domain exceptions to HTTP status once, in one place — extend
 * this rather than catching exceptions deep inside the resource. */
@Provider
public class RestExceptionMapper implements ExceptionMapper<RuntimeException> {

    @Override
    public Response toResponse(RuntimeException ex) {
        if (ex instanceof TaskNotFoundException) {
            return Response.status(Response.Status.NOT_FOUND).entity(ex.getMessage()).build();
        }
        if (ex instanceof InvalidTaskTitleException) {
            return Response.status(Response.Status.BAD_REQUEST).entity(ex.getMessage()).build();
        }
        return Response.serverError().entity("internal error").build();
    }
}
