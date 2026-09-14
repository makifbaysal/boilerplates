package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.graphql;

import org.springframework.graphql.execution.DataFetcherExceptionResolverAdapter;
import org.springframework.graphql.execution.ErrorType;
import org.springframework.stereotype.Component;

import graphql.GraphQLError;
import graphql.schema.DataFetchingEnvironment;

import ai.tasktrooper.boilerplate.javaspring.core.task.InvalidTaskTitleException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;

/**
 * Maps domain exceptions to GraphQL error types once, in one place —
 * the GraphQL sibling of rest.RestExceptionHandler / grpc's mapError.
 * Without this, Spring GraphQL logs them as "Unresolved ..." and
 * returns an opaque INTERNAL_ERROR.
 */
@Component
public class GraphQLExceptionResolver extends DataFetcherExceptionResolverAdapter {

    @Override
    protected GraphQLError resolveToSingleError(Throwable ex, DataFetchingEnvironment env) {
        if (ex instanceof TaskNotFoundException) {
            return toError(ex, ErrorType.NOT_FOUND, env);
        }
        if (ex instanceof InvalidTaskTitleException) {
            return toError(ex, ErrorType.BAD_REQUEST, env);
        }
        return null; // let the default resolver classify anything else
    }

    private static GraphQLError toError(Throwable ex, ErrorType type, DataFetchingEnvironment env) {
        return GraphQLError.newError()
                .errorType(type)
                .message(ex.getMessage())
                .path(env.getExecutionStepInfo().getPath())
                .location(env.getField().getSourceLocation())
                .build();
    }
}
