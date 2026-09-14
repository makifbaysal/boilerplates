# java-spring boilerplate

Java + Spring Boot backend in hexagonal (ports & adapters) architecture:
one `task` CRUD use case exposed over three inbound transports — REST
(Spring MVC), gRPC (net.devh starter), and GraphQL (spring-boot-starter-graphql)
— with three outbound storage implementations (Postgres/JPA, MongoDB,
Cassandra, plus the in-memory default), MDC-based structured logging
(REST filter + gRPC interceptor), and profile-based config with a live
log-level file watch.

Read [.ai/architecture.md](.ai/architecture.md) before changing
structure — it explains the layering, why all-three-DB-autoconfig is
excluded by default, and why this is one process instead of go-fiber's
three binaries. Read [.ai/common-tasks.md](.ai/common-tasks.md) before
adding an endpoint/RPC/query or resource. [.ai/coding-standards.md](.ai/coding-standards.md)
and [.ai/testing.md](.ai/testing.md) cover style and test conventions —
same shape as `backend/go-fiber`'s, adapted for Spring/JUnit5/Mockito.

## Commands

```
mvn spring-boot:run              # REST+GraphQL on :8080, gRPC on :9090
mvn test                         # unit tests (surefire, *Test.java)
mvn verify                       # + integration tests (failsafe, *IT.java) — needs Docker
mvn package                      # -> target/java-spring-*.jar
```

Activate a database: `SPRING_PROFILES_ACTIVE=dev,postgres mvn spring-boot:run`
(`postgres`|`mongo`|`cassandra`; omit for the in-memory default).
`docker-compose.dev.yml` starts local Postgres/MongoDB/Cassandra/Redis/
Elasticsearch for that.
