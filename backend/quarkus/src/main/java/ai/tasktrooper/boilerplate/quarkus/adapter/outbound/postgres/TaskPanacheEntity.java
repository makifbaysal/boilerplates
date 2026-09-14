package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres;

import java.time.Instant;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

/** This adapter's own persistence shape — never annotate core.Task with
 * JPA annotations directly. See schema.sql for the table this expects. */
@Entity
@Table(name = "tasks")
public class TaskPanacheEntity {

    @Id
    public String id;
    public String title;
    public boolean done;
    public Instant createdAt;
    public Instant updatedAt;
}
