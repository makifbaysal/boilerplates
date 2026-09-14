package ai.tasktrooper.androidnative.ui

import ai.tasktrooper.androidnative.data.CreateTaskRequest
import ai.tasktrooper.androidnative.data.SetDoneRequest
import ai.tasktrooper.androidnative.data.Task
import ai.tasktrooper.androidnative.data.TaskApiService
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

sealed interface TaskUiState {
    data object Loading : TaskUiState
    data class Loaded(val tasks: List<Task>, val mutatingId: String? = null) : TaskUiState
    data class Error(val message: String) : TaskUiState
}

@HiltViewModel
class TaskViewModel @Inject constructor(
    private val api: TaskApiService,
) : ViewModel() {

    private val _state = MutableStateFlow<TaskUiState>(TaskUiState.Loading)
    val state: StateFlow<TaskUiState> = _state.asStateFlow()

    init {
        refresh()
    }

    fun refresh() {
        viewModelScope.launch {
            _state.value = TaskUiState.Loading
            _state.value = runCatching { api.listTasks() }
                .fold({ TaskUiState.Loaded(it) }, { TaskUiState.Error(it.message ?: "Unknown error") })
        }
    }

    fun createTask(title: String) {
        viewModelScope.launch {
            runCatching { api.createTask(CreateTaskRequest(title)) }
                .onSuccess { if (it.isSuccessful) refresh() }
        }
    }

    fun setDone(id: String, done: Boolean) {
        viewModelScope.launch {
            runCatching { api.setTaskDone(id, SetDoneRequest(done)) }
                .onSuccess { refresh() }
        }
    }

    fun delete(id: String) {
        viewModelScope.launch {
            runCatching { api.deleteTask(id) }
                .onSuccess { refresh() }
        }
    }
}
