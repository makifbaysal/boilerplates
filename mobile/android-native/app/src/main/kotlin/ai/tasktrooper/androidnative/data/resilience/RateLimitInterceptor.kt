package ai.tasktrooper.androidnative.data.resilience

import java.io.IOException
import okhttp3.Interceptor
import okhttp3.Response

/**
 * Token-bucket throttle for outbound requests: it keeps a recomposition loop
 * or an over-eager screen from turning into a request storm the server has to
 * absorb. Fails the call rather than queueing, so the bug is visible in
 * development instead of silently doubling latency.
 */
class RateLimitInterceptor(
    private val perSecond: Double = 8.0,
    private val burst: Double = 16.0,
    private val now: () -> Long = System::currentTimeMillis,
) : Interceptor {

    class RateLimitedException : IOException("client rate limit exceeded")

    private val lock = Any()
    private var tokens = burst
    private var last = 0L

    override fun intercept(chain: Interceptor.Chain): Response {
        if (!tryAcquire()) throw RateLimitedException()
        return chain.proceed(chain.request())
    }

    private fun tryAcquire(): Boolean = synchronized(lock) {
        if (perSecond <= 0) return true
        val current = now()
        if (last == 0L) last = current
        tokens = minOf(burst, tokens + (current - last) / 1000.0 * perSecond)
        last = current
        if (tokens < 1) return false
        tokens -= 1
        true
    }
}
