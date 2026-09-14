package ai.tasktrooper.boilerplate.javaspring.adapter.inbound.rest;

import static org.hamcrest.Matchers.notNullValue;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.delete;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.patch;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.context.annotation.Import;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import ai.tasktrooper.boilerplate.javaspring.adapter.outbound.memory.InMemoryTaskRepository;
import ai.tasktrooper.boilerplate.javaspring.core.task.TaskService;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Runs the REST adapter against a real InMemoryTaskRepository — no
 * mocks, same "adapter test over the real use case" pattern as
 * backend/go-fiber's rest.HandlerSuite. See .ai/testing.md.
 */
@WebMvcTest(TaskController.class)
@Import({TaskService.class, InMemoryTaskRepository.class})
class TaskControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @ParameterizedTest(name = "title=\"{0}\" -> {1}")
    @CsvSource({"buy milk,201", "'   ',400"})
    void create(String title, int wantStatus) throws Exception {
        mockMvc.perform(post("/tasks")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(new java.util.HashMap<>(java.util.Map.of("title", title)))))
                .andExpect(status().is(wantStatus));
    }

    @Test
    void crudFlow() throws Exception {
        String body = mockMvc.perform(post("/tasks")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"title\":\"buy milk\"}"))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id", notNullValue()))
                .andReturn().getResponse().getContentAsString();

        String id = objectMapper.readTree(body).get("id").asText();

        mockMvc.perform(get("/tasks/" + id)).andExpect(status().isOk());
        mockMvc.perform(get("/tasks")).andExpect(status().isOk());
        mockMvc.perform(patch("/tasks/" + id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"done\":true}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.done").value(true));
        mockMvc.perform(delete("/tasks/" + id)).andExpect(status().isNoContent());
        mockMvc.perform(get("/tasks/" + id)).andExpect(status().isNotFound());
    }
}
