package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import ai.tasktrooper.boilerplate.javaspring.core.task.InvalidTaskTitleException;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskNotFoundException;

/** Maps domain exceptions to HTTP status once, in one place — extend
 * this rather than catching exceptions deep inside the controller. */
@RestControllerAdvice
public class RestExceptionHandler {

    @ExceptionHandler(TaskNotFoundException.class)
    public ResponseEntity<String> handleNotFound(TaskNotFoundException ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(ex.getMessage());
    }

    @ExceptionHandler(InvalidTaskTitleException.class)
    public ResponseEntity<String> handleInvalid(InvalidTaskTitleException ex) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(ex.getMessage());
    }
}
