package ai.tasktrooper.androidnative.data

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.GET
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.Path

/** The only I/O boundary — everything else in the app depends on this
 * interface, not on Retrofit/OkHttp directly. */
interface TaskApiService {
    @GET("tasks")
    suspend fun listTasks(): List<Task>

    @POST("tasks")
    suspend fun createTask(@Body request: CreateTaskRequest): Response<Task>

    @PATCH("tasks/{id}")
    suspend fun setTaskDone(@Path("id") id: String, @Body request: SetDoneRequest): Task

    @DELETE("tasks/{id}")
    suspend fun deleteTask(@Path("id") id: String): Response<Unit>
}
