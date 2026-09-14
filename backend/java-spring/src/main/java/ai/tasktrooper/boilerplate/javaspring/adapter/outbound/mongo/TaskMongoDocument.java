package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.mongo;

import java.time.Instant;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

/** This adapter's own persistence shape — never annotate core.Task with
 * Mongo mapping annotations directly. */
@Document(collection = "tasks")
@Getter
@NoArgsConstructor
@AllArgsConstructor
public class TaskMongoDocument {

    @Id
    private String id;
    private String title;
    private boolean done;
    private Instant createdAt;
    private Instant updatedAt;
}
