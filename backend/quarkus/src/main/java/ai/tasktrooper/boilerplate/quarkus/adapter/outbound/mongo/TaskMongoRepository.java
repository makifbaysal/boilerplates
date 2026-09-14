package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.mongo;

import jakarta.enterprise.context.ApplicationScoped;

import io.quarkus.mongodb.panache.PanacheMongoRepositoryBase;

@ApplicationScoped
public class TaskMongoRepository implements PanacheMongoRepositoryBase<TaskMongoEntity, String> {
}
