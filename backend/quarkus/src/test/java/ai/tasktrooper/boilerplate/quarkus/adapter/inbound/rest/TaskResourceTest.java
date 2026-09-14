package ai.tasktrooper.boilerplate.quarkus.adapter.inbound.rest;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.notNullValue;
import static org.hamcrest.Matchers.is;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.api.Test;

import io.quarkus.test.junit.QuarkusTest;
import io.restassured.http.ContentType;

/**
 * Runs the REST adapter against the real (in-memory, since no
 * app.datastore profile is set for the test run) TaskUseCase — the
 * Quarkus equivalent of go-fiber's rest.HandlerSuite / java-spring's
 * TaskControllerTest. See .ai/testing.md.
 */
@QuarkusTest
class TaskResourceTest {

    @ParameterizedTest(name = "title=\"{0}\" -> {1}")
    @CsvSource({"buy milk,201", "'   ',400"})
    void create(String title, int wantStatus) {
        given().contentType(ContentType.JSON).body("{\"title\":\"" + title + "\"}")
                .when().post("/tasks")
                .then().statusCode(wantStatus);
    }

    @Test
    void crudFlow() {
        String id = given().contentType(ContentType.JSON).body("{\"title\":\"buy milk\"}")
                .when().post("/tasks")
                .then().statusCode(201).body("id", notNullValue())
                .extract().path("id");

        given().when().get("/tasks/" + id).then().statusCode(200);
        given().when().get("/tasks").then().statusCode(200);

        given().contentType(ContentType.JSON).body("{\"done\":true}")
                .when().patch("/tasks/" + id)
                .then().statusCode(200).body("done", is(true));

        given().when().delete("/tasks/" + id).then().statusCode(204);
        given().when().get("/tasks/" + id).then().statusCode(404);
    }
}
