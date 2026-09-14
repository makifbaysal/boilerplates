# Architecture

Kotlin + Jetpack Compose, Hilt DI (KSP, not kapt), Retrofit +
OkHttp + kotlinx.serialization for the only I/O boundary. AGP 9+'s
built-in Kotlin support — no `org.jetbrains.kotlin.android` plugin (see
`app/build.gradle.kts` comment; that plugin is gone as of AGP 9, and
removing it also removes the `kotlinOptions {}` DSL — JVM target now
comes from `compileOptions`).

```
data/Task.kt              @Serializable models (request/response DTOs)
data/TaskApiService.kt      Retrofit interface — the only I/O boundary
di/NetworkModule.kt          Hilt module providing OkHttp/Retrofit/the service
ui/TaskViewModel.kt          @HiltViewModel, StateFlow<TaskUiState>
ui/TaskScreen.kt              Compose screen, hiltViewModel()
ui/components/TaskItem.kt      pure Composable, props in/callbacks out
```

## Real bugs found building this

- **`retrofit2-kotlinx-serialization-converter`'s actual package is
  `com.jakewharton.retrofit2.converter.kotlinx.serialization`**, not
  `retrofit2.converter.kotlinx.serialization` — an easy mistake since
  the Maven coordinate itself doesn't have the `com.jakewharton` prefix.
- **A KDoc comment containing a literal `backend/*` breaks compilation.**
  Kotlin block comments nest, so `/** ... backend/* ... */` opens a
  second unclosed comment. Avoid `/*` inside `/** */` docs — use `//`
  or reword.
- **`androidx.core`/`androidx.activity` 1.11.x require `compileSdk 36`**
  — `compileSdk 35` fails `checkDebugAarMetadata` with a clear message;
  this repo installs `platforms;android-36` for exactly that reason.

## Config

`BuildConfig.API_BASE_URL` (`app/build.gradle.kts` →
`buildConfigField`) defaults to `http://10.0.2.2:8080` — the Android
emulator's alias for the host machine's localhost, so it reaches a
backend/* boilerplate running on your dev machine unmodified.

## Testing

Unit tests (`app/src/test`) mock only the true I/O boundary
(`MockWebServer`, real Retrofit service on top) — same principle as
every other boilerplate here. No instrumented/Compose UI tests are
included (they need a device/emulator, out of scope for what CI can run
here) — add them under `app/src/androidTest` following the same
mocked-server pattern if your project needs them.
