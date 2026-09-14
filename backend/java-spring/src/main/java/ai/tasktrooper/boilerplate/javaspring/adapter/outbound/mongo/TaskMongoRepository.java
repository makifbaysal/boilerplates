package ai.tasktrooper.boilerplate.javaspring.adapter.outbound.mongo;

import org.springframework.data.mongodb.repository.MongoRepository;

interface TaskMongoRepository extends MongoRepository<TaskMongoDocument, String> {
}
