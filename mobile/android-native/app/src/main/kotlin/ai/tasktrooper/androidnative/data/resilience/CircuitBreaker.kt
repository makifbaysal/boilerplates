package ai.tasktrooper.androidnative.data.resilience

import java.io.IOException
import java.util.concurrent.TimeUnit
import okhttp3.Interceptor
import okhttp3.Response

/**
 * Client-side circuit breaker as an OkHttp interceptor. A phone that keeps
 * retrying a dead API drains the battery and delays the moment the UI can
 * honestly say "offline", so failing fast is a UX decision here.
 *
 * Only transport errors and 5xx count as failures: a 404 means this request
 * was wrong, not that the backend is down.
 */
class CircuitBreakerInterceptor(
    private val failureThreshold: Int = 4,
    private val openDurationMs: Long = TimeUnit.SECONDS.toMillis(20),
    private val now: () -> Long = System::currentTimeMillis,
) : Interceptor {

    class OpenCircuitException(retryInMs: Long) :
        IOException("circuit breaker is open, retry in ${retryInMs / 1000}s")

    private enum class State { CLOSED, OPEN, HALF_OPEN }

    private val lock = Any()
    private var state = State.CLOSED
    private var failures = 0
    private var openedAt = 0L

    override fun intercept(chain: Interceptor.Chain): Response {
        synchronized(lock) {
            if (state == State.OPEN) {
                val elapsed = now() - openedAt
                if (elapsed < openDurationMs) {
                    throw OpenCircuitException(openDurationMs - elapsed)
                }
                state = State.HALF_OPEN
            }
        }

        val response = try {
            chain.proceed(chain.request())
        } catch (e: IOException) {
            onFailure()
            throw e
        }

        if (response.code >= 500) {
            onFailure()
        } else {
            onSuccess()
        }
        return response
    }

    private fun onSuccess() = synchronized(lock) {
        state = State.CLOSED
        failures = 0
    }

    private fun onFailure() = synchronized(lock) {
        if (state == State.HALF_OPEN) {
            state = State.OPEN
            openedAt = now()
            return@synchronized
        }
        failures++
        if (failures >= failureThreshold) {
            state = State.OPEN
            openedAt = now()
        }
    }
}
