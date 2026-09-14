# android-native boilerplate

Kotlin + Jetpack Compose + Hilt + Retrofit Task CRUD app calling the
same REST shape every `backend/*` boilerplate exposes.

Read [.ai/architecture.md](.ai/architecture.md) first — it documents
three real bugs hit building this (a retrofit-kotlinx-serialization
package name gotcha, a Kotlin nested-comment syntax trap, and AGP 9's
compileSdk 36 requirement) so they don't regress.

## Commands

```
./gradlew assembleDebug   # -> app/build/outputs/apk/debug/
./gradlew test              # unit tests (MockWebServer, no device needed)
```

`BuildConfig.API_BASE_URL` defaults to the Android emulator's
localhost alias (`10.0.2.2`) — see architecture.md#config.
