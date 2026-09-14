package ai.tasktrooper.androidnative

import ai.tasktrooper.androidnative.data.CreateTaskRequest
import ai.tasktrooper.androidnative.data.TaskApiService
import com.google.common.truth.Truth.assertThat
import kotlinx.coroutines.test.runTest
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Before
import org.junit.Test
import retrofit2.Retrofit
import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory

// Mocks only the true I/O boundary (the HTTP server, via MockWebServer)
// and exercises the real Retrofit service — same principle as every
// other boilerplate in this repo.
class TaskApiServiceTest {

    private lateinit var server: MockWebServer
    private lateinit var api: TaskApiService

    @Before
    fun setUp() {
        server = MockWebServer()
        server.start()
        val json = Json { ignoreUnknownKeys = true }
        api = Retrofit.Builder()
            .baseUrl(server.url("/"))
            .addConverterFactory(json.asConverterFactory("application/json".toMediaType()))
            .build()
            .create(TaskApiService::class.java)
    }

    @After
    fun tearDown() {
        server.shutdown()
    }

    @Test
    fun listTasks_parsesResponseBody() = runTest {
        server.enqueue(
            MockResponse().setBody(
                """[{"id":"1","title":"buy milk","done":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]""",
            ).setHeader("Content-Type", "application/json"),
        )

        val tasks = api.listTasks()

        assertThat(tasks).hasSize(1)
        assertThat(tasks.first().title).isEqualTo("buy milk")
    }

    @Test
    fun createTask_sendsTitleAsJsonBody() = runTest {
        server.enqueue(
            MockResponse().setResponseCode(201).setBody(
                """{"id":"1","title":"buy milk","done":false,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}""",
            ).setHeader("Content-Type", "application/json"),
        )

        val response = api.createTask(CreateTaskRequest("buy milk"))

        assertThat(response.isSuccessful).isTrue()
        val request = server.takeRequest()
        assertThat(request.path).isEqualTo("/tasks")
        assertThat(request.body.readUtf8()).contains("buy milk")
    }
}
