import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;
import 'package:http/testing.dart';

import 'package:flutter_app/api/task_api_client.dart';
import 'package:flutter_app/providers/task_providers.dart';
import 'package:flutter_app/screens/task_screen.dart';

void main() {
  testWidgets('TaskScreen renders tasks from the API', (tester) async {
    final mockClient = MockClient((request) async {
      if (request.method == 'GET' && request.url.path == '/tasks') {
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
      }
      return http.Response('not found', 404);
    });

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          taskApiClientProvider.overrideWithValue(
            TaskApiClient(baseUrl: 'http://test', client: mockClient),
          ),
        ],
        child: const MaterialApp(home: TaskScreen()),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.text('buy milk'), findsOneWidget);
  });
}
