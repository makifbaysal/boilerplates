FROM maven:3.9-eclipse-temurin-21-alpine AS build
WORKDIR /src
COPY pom.xml .
RUN mvn -q -B dependency:go-offline
COPY src ./src
RUN mvn -q -B package -DskipTests

FROM eclipse-temurin:21-jre-alpine
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /src/target/java-spring-*.jar /app.jar
COPY --from=build /src/src/main/resources/application*.yml /config/
ENV CONFIG_DIR=/config
# Load config from the external, editable copy (not just the jar's
# embedded classpath one) so LiveLogLevelWatcher's target file is what
# Spring actually reads at startup too — see .ai/architecture.md#config.
ENV SPRING_CONFIG_ADDITIONAL_LOCATION=file:/config/
# SPRING_PROFILES_ACTIVE: see .ai/architecture.md#config — deployer sets it, not this file.
USER app:app
EXPOSE 8080 9090
ENTRYPOINT ["java", "-jar", "/app.jar"]
