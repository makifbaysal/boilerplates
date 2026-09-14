package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.cassandra;

import java.time.Instant;

import org.springframework.data.cassandra.core.mapping.PrimaryKey;
import org.springframework.data.cassandra.core.mapping.Table;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

/** This adapter's own persistence shape. See schema.cql for the table
 * this expects. */
@Table("tasks")
@Getter
@NoArgsConstructor
@AllArgsConstructor
public class TaskCassandraEntity {

    @PrimaryKey
    private String id;
    private String title;
    private boolean done;
    private Instant createdAt;
    private Instant updatedAt;
}
