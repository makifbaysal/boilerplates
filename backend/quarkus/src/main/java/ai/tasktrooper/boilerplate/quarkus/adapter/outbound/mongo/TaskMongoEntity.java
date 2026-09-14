package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.mongo;

import java.time.Instant;

import org.bson.codecs.pojo.annotations.BsonId;

import io.quarkus.mongodb.panache.common.MongoEntity;

/** This adapter's own persistence shape — never annotate core.Task with
 * Mongo mapping annotations directly. */
@MongoEntity(collection = "tasks")
public class TaskMongoEntity {

    @BsonId
    public String id;
    public String title;
    public boolean done;
    public Instant createdAt;
    public Instant updatedAt;
}
