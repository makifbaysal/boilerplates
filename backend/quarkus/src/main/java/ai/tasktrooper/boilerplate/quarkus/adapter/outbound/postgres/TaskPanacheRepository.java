package ai.tasktrooper.boilerplate.quarkus.adapter.outbound.postgres;

import jakarta.enterprise.context.ApplicationScoped;

import io.quarkus.hibernate.orm.panache.PanacheRepositoryBase;

// PanacheRepository<Entity> defaults to a Long id; our id is a String,
// so this needs the explicit-id-type base interface instead.
@ApplicationScoped
public class TaskPanacheRepository implements PanacheRepositoryBase<TaskPanacheEntity, String> {
}
