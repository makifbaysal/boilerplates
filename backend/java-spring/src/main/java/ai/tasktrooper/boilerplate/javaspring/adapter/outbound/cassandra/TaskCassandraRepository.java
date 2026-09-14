package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.cassandra;

import org.springframework.data.cassandra.repository.CassandraRepository;

interface TaskCassandraRepository extends CassandraRepository<TaskCassandraEntity, String> {
}
