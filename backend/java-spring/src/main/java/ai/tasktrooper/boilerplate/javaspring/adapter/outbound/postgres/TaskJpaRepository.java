package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.postgres;

import org.springframework.data.jpa.repository.JpaRepository;

interface TaskJpaRepository extends JpaRepository<TaskJpaEntity, String> {
}
