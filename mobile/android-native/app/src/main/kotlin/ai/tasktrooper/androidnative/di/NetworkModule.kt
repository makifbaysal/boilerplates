package ai.tasktrooper.androidnative.di

import ai.tasktrooper.androidnative.BuildConfig
import ai.tasktrooper.androidnative.data.TaskApiService
import ai.tasktrooper.androidnative.data.resilience.CircuitBreakerInterceptor
import ai.tasktrooper.androidnative.data.resilience.RateLimitInterceptor
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    fun provideJson(): Json = Json { ignoreUnknownKeys = true }

    @Provides
    @Singleton
    fun provideOkHttpClient(): OkHttpClient =
        OkHttpClient.Builder()
            // Throttle first (cheapest rejection), then the breaker, then the
            // logger so what is logged is what actually left the device.
            .addInterceptor(RateLimitInterceptor())
            .addInterceptor(CircuitBreakerInterceptor())
            .addInterceptor(HttpLoggingInterceptor().apply { level = HttpLoggingInterceptor.Level.BASIC })
            .build()

    @Provides
    @Singleton
    fun provideRetrofit(client: OkHttpClient, json: Json): Retrofit =
        Retrofit.Builder()
            .baseUrl("${BuildConfig.API_BASE_URL}/")
            .client(client)
            .addConverterFactory(json.asConverterFactory("application/json".toMediaType()))
            .build()

    @Provides
    @Singleton
    fun provideTaskApiService(retrofit: Retrofit): TaskApiService = retrofit.create(TaskApiService::class.java)
}
