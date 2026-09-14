package ai.tasktrooper.androidnative.data

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

// Matches the Task shape every backend boilerplate in this repo exposes.
@Serializable
data class Task(
    val id: String,
    val title: String,
    val done: Boolean,
    @SerialName("created_at") val createdAt: String,
    @SerialName("updated_at") val updatedAt: String,
)

@Serializable
data class CreateTaskRequest(val title: String)

@Serializable
data class SetDoneRequest(val done: Boolean)
