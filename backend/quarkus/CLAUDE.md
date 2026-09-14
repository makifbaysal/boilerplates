# quarkus boilerplate

Java + Quarkus backend in hexagonal architecture: one `task` CRUD use
case exposed over REST (Quarkus REST), gRPC (quarkus-grpc), and GraphQL
(SmallRye GraphQL, code-first) — with Postgres(Panache)/MongoDB(Panache)/
Cassandra(plain DataStax driver) outbound adapters, MDC logging via a
JAX-RS filter + gRPC interceptor.

Read [.ai/architecture.md](.ai/architecture.md) first — it documents
five real Quarkus-specific bugs (CDI ambiguity, eager Hibernate ORM
activation, test-resource leaking across classes, a testcontainers
version mismatch, GraphQL's `set`-prefix method-name stripping) that
this boilerplate's structure exists specifically to avoid regressing.
[.ai/common-tasks.md](.ai/common-tasks.md) and
[.ai/testing.md](.ai/testing.md)/[.ai/coding-standards.md](.ai/coding-standards.md)
cover the rest, same shape as `backend/go-fiber`/`backend/java-spring`.

## Commands

```
mvn quarkus:dev     # REST+GraphQL on :8080, gRPC on :9090, live reload
mvn test            # unit/adapter tests (no DB needed)
mvn verify          # + Testcontainers DB integration tests — needs Docker
mvn package          # -> target/quarkus-app/
```

Activate a database: `app.datastore=postgres mvn quarkus:dev` (or
`mongo`/`cassandra`; omit for in-memory). `docker-compose.dev.yml`
starts local Postgres/MongoDB/Cassandra/Redis/Elasticsearch for that.
