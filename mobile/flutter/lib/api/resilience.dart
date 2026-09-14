import 'dart:async';

/// Thrown instead of calling the API while the breaker is open.
class CircuitOpenException implements Exception {
  CircuitOpenException(this.retryAfter);
  final Duration retryAfter;
  @override
  String toString() => 'CircuitOpenException(retry in ${retryAfter.inSeconds}s)';
}

enum BreakerState { closed, open, halfOpen }

/// Client-side circuit breaker. A phone that keeps retrying a dead API drains
/// the battery and delays the moment the UI can honestly say "offline", so
/// failing fast is a UX decision here, not just an infrastructure one.
class CircuitBreaker {
  CircuitBreaker({
    this.failureThreshold = 4,
    this.openDuration = const Duration(seconds: 20),
    DateTime Function()? now,
  }) : _now = now ?? DateTime.now;

  final int failureThreshold;
  final Duration openDuration;
  final DateTime Function() _now;

  BreakerState _state = BreakerState.closed;
  int _failures = 0;
  DateTime? _openedAt;

  BreakerState get state {
    _refresh();
    return _state;
  }

  Future<T> run<T>(Future<T> Function() action) async {
    _refresh();
    if (_state == BreakerState.open) {
      final elapsed = _now().difference(_openedAt!);
      throw CircuitOpenException(openDuration - elapsed);
    }
    try {
      final result = await action();
      _onSuccess();
      return result;
    } catch (_) {
      _onFailure();
      rethrow;
    }
  }

  void _refresh() {
    if (_state == BreakerState.open &&
        _now().difference(_openedAt!) >= openDuration) {
      _state = BreakerState.halfOpen;
    }
  }

  void _onSuccess() {
    _state = BreakerState.closed;
    _failures = 0;
  }

  void _onFailure() {
    if (_state == BreakerState.halfOpen) {
      _state = BreakerState.open;
      _openedAt = _now();
      return;
    }
    _failures++;
    if (_failures >= failureThreshold) {
      _state = BreakerState.open;
      _openedAt = _now();
    }
  }
}

/// Token-bucket throttle for outbound calls: it keeps a chatty screen (or a
/// rebuild loop) from turning into a request storm the server has to absorb.
class RateLimiter {
  RateLimiter({this.perSecond = 5, this.burst = 10, DateTime Function()? now})
      : _tokens = burst.toDouble(),
        _now = now ?? DateTime.now;

  final double perSecond;
  final int burst;
  final DateTime Function() _now;

  double _tokens;
  DateTime? _last;

  bool tryAcquire() {
    if (perSecond <= 0) return true;
    final now = _now();
    _last ??= now;
    _tokens += now.difference(_last!).inMicroseconds / 1e6 * perSecond;
    if (_tokens > burst) _tokens = burst.toDouble();
    _last = now;
    if (_tokens < 1) return false;
    _tokens -= 1;
    return true;
  }
}
