import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/task.dart';
import 'resilience.dart';

class ApiException implements Exception {
  ApiException(this.statusCode, this.body);
  final int statusCode;
  final String body;
  @override
  String toString() => 'ApiException($statusCode): $body';
}

/// Calls the same REST shape every backend boilerplate exposes.
/// Point [baseUrl] at go-fiber/java-spring/quarkus.
class TaskApiClient {
  TaskApiClient({
    required this.baseUrl,
    http.Client? client,
    CircuitBreaker? breaker,
    RateLimiter? limiter,
  })  : _client = client ?? http.Client(),
        _breaker = breaker ?? CircuitBreaker(),
        _limiter = limiter ?? RateLimiter();

  final String baseUrl;
  final http.Client _client;
  final CircuitBreaker _breaker;
  final RateLimiter _limiter;

  /// Every call goes through the throttle and the breaker: one dead backend
  /// should cost one failed request, not one per widget rebuild.
  Future<T> _guard<T>(Future<T> Function() action) {
    if (!_limiter.tryAcquire()) {
      throw ApiException(429, 'client rate limit exceeded');
    }
    return _breaker.run(action);
  }

  Future<List<Task>> listTasks() async {
    final res = await _guard(() => _client.get(Uri.parse('$baseUrl/tasks')));
    _checkOk(res, 200);
    final list = jsonDecode(res.body) as List<dynamic>;
    return list.map((e) => Task.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Task> createTask(String title) async {
    final res = await _guard(() => _client.post(
      Uri.parse('$baseUrl/tasks'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'title': title}),
    ));
    _checkOk(res, 201);
    return Task.fromJson(jsonDecode(res.body) as Map<String, dynamic>);
  }

  Future<Task> setTaskDone(String id, bool done) async {
    final res = await _guard(() => _client.patch(
      Uri.parse('$baseUrl/tasks/$id'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'done': done}),
    ));
    _checkOk(res, 200);
    return Task.fromJson(jsonDecode(res.body) as Map<String, dynamic>);
  }

  Future<void> deleteTask(String id) async {
    final res = await _guard(() => _client.delete(Uri.parse('$baseUrl/tasks/$id')));
    _checkOk(res, 204);
  }

  void _checkOk(http.Response res, int want) {
    if (res.statusCode != want) {
      throw ApiException(res.statusCode, res.body);
    }
  }
}
