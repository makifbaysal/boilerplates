package ai.tasktrooper.boilerplate.javaspring.core.task;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.argThat;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.time.Instant;
import java.util.List;
import java.util.Optional;
import java.util.stream.Stream;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

/**
 * Mockito (mock) + JUnit 5 {@code @Nested} (suite grouping) +
 * {@code @ParameterizedTest} (table-driven) — the Java equivalent of
 * backend/go-fiber's mockery + testify-suite + table-test pattern. See
 * .ai/testing.md.
 */
@ExtendWith(MockitoExtension.class)
class TaskServiceTest {

    @Mock
    private TaskRepositoryPort repository;

    private TaskService service;

    @BeforeEach
    void setUp() {
        service = new TaskService(repository);
    }

    @Nested
    class Create {

        static Stream<Arguments> titles() {
            return Stream.of(
                    Arguments.of("  buy milk  ", "buy milk", true),
                    Arguments.of("   ", null, false));
        }

        @ParameterizedTest(name = "title=\"{0}\"")
        @MethodSource("titles")
        void create(String input, String wantTitle, boolean wantOk) {
            if (wantOk) {
                when(repository.create(any())).thenAnswer(inv -> inv.getArgument(0));

                Task got = service.create(input);

                assertThat(got.title()).isEqualTo(wantTitle);
                assertThat(got.done()).isFalse();
                assertThat(got.id()).isNotBlank();
            } else {
                assertThatThrownBy(() -> service.create(input)).isInstanceOf(InvalidTaskTitleException.class);
                verify(repository, never()).create(any());
            }
        }
    }

    @Nested
    class SetDone {

        @Test
        void marksAnExistingTaskDone() {
            Instant now = Instant.now();
            Task existing = new Task("abc", "buy milk", false, now, now);
            when(repository.get("abc")).thenReturn(Optional.of(existing));
            when(repository.update(any())).thenAnswer(inv -> inv.getArgument(0));

            Task got = service.setDone("abc", true);

            assertThat(got.done()).isTrue();
            verify(repository).update(argThat(Task::done));
        }

        @Test
        void propagatesNotFoundFromTheRepository() {
            when(repository.get("missing")).thenReturn(Optional.empty());

            assertThatThrownBy(() -> service.setDone("missing", true))
                    .isInstanceOf(TaskNotFoundException.class);
        }
    }

    @Test
    void listDelegatesToRepository() {
        Instant now = Instant.now();
        List<Task> tasks = List.of(new Task("a", "one", false, now, now));
        when(repository.list()).thenReturn(tasks);

        assertThat(service.list()).isEqualTo(tasks);
    }
}
