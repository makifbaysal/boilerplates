plugins {
    id("com.android.application") version "9.3.1" apply false
    // AGP 9+ has built-in Kotlin support — see app/build.gradle.kts comment.
    id("org.jetbrains.kotlin.plugin.compose") version "2.3.10" apply false
    id("org.jetbrains.kotlin.plugin.serialization") version "2.3.10" apply false
    id("com.google.dagger.hilt.android") version "2.60.1" apply false
    id("com.google.devtools.ksp") version "2.3.10" apply false
}
