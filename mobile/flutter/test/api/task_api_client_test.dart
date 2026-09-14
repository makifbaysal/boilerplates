import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';
import 'package:flutter_app/api/task_api_client.dart';

void main() {
  group('TaskApiClient', () {
    test('listTasks parses the response body', () async {
      final client = TaskApiClient(
        baseUrl: 'http://test',
        client: MockClient((request) async {
          expect(request.url.path, '/tasks');
          return http.Response(
            jsonEncode([
              {
                'id': '1',
                'title': 'buy milk',
                'done': false,
                'created_at': '2026-01-01T00:00:00Z',
                'updated_at': '2026-01-01T00:00:00Z',
              },
            ]),
            200,
          );
        }),
      );

      final tasks = await client.listTasks();

      expect(tasks, hasLength(1));
      expect(tasks.first.title, 'buy milk');
    });

    test('createTask throws ApiException on a non-201 response', () async {
      final client = TaskApiClient(
        baseUrl: 'http://test',
        client: MockClient((request) async => http.Response('bad title', 400)),
      );

      await expectLater(client.createTask('   '), throwsA(isA<ApiException>()));
    });

    test('deleteTask succeeds on 204', () async {
      final client = TaskApiClient(
        baseUrl: 'http://test',
        client: MockClient((request) async {
          expect(request.method, 'DELETE');
          return http.Response('', 204);
        }),
      );

      await expectLater(client.deleteTask('1'), completes);
    });
  });
}
