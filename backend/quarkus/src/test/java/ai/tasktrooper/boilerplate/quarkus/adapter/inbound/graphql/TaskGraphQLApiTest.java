package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.graphql;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.notNullValue;
import static org.hamcrest.Matchers.nullValue;

import org.junit.jupiter.api.Test;

import io.quarkus.test.junit.QuarkusTest;
import io.restassured.http.ContentType;

/**
 * Runs the GraphQL adapter (POST /graphql, SmallRye's default endpoint)
 * against the real in-memory TaskUseCase — the GraphQL sibling of
 * TaskResourceTest / TaskGrpcServiceTest.
 */
@QuarkusTest
class TaskGraphQLApiTest {

    @Test
    void createRejectsBlankTitle() {
        given().contentType(ContentType.JSON)
                .body("{\"query\":\"mutation { createTask(title: \\\"   \\\") { id } }\"}")
                .when().post("/graphql")
                .then().statusCode(200)
                .body("errors[0].message", org.hamcrest.Matchers.equalTo("task title is required"));
    }

    @Test
    void crudFlow() {
        String id = given().contentType(ContentType.JSON)
                .body("{\"query\":\"mutation { createTask(title: \\\"buy milk\\\") { id } }\"}")
                .when().post("/graphql")
                .then().statusCode(200).body("data.createTask.id", notNullValue())
                .extract().path("data.createTask.id");

        given().contentType(ContentType.JSON).body("{\"query\":\"{ tasks { id } }\"}")
                .when().post("/graphql")
                .then().statusCode(200).body("data.tasks", notNullValue());

        given().contentType(ContentType.JSON)
                .body("{\"query\":\"mutation($id: String!) { setTaskDone(id: $id, done: true) { done } }\",\"variables\":{\"id\":\"" + id + "\"}}")
                .when().post("/graphql")
                .then().statusCode(200).body("data.setTaskDone.done", org.hamcrest.Matchers.is(true));

        given().contentType(ContentType.JSON)
                .body("{\"query\":\"mutation($id: String!) { deleteTask(id: $id) }\",\"variables\":{\"id\":\"" + id + "\"}}")
                .when().post("/graphql")
                .then().statusCode(200).body("data.deleteTask", org.hamcrest.Matchers.is(true));

        given().contentType(ContentType.JSON)
                .body("{\"query\":\"query($id: String!) { task(id: $id) { id } }\",\"variables\":{\"id\":\"" + id + "\"}}")
                .when().post("/graphql")
                .then().statusCode(200).body("data.task", nullValue());
    }
}
