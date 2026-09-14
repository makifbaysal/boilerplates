"use client";

import { useQueryClient } from "@tanstack/react-query";

import { TaskItem } from "@/components/molecules/TaskItem";
import {
  useListTasks,
  useSetTaskDone,
  useDeleteTask,
  getListTasksQueryKey,
} from "@/lib/api/generated/default/default";

export function TaskList() {
  const queryClient = useQueryClient();
  const { data: tasks, isLoading, isError } = useListTasks();

  const invalidate = () => queryClient.invalidateQueries({ queryKey: getListTasksQueryKey() });
  const { mutate: setDone, variables: setDoneVars, isPending: isSettingDone } = useSetTaskDone({
    mutation: { onSuccess: invalidate },
  });
  const { mutate: deleteTask, variables: deleteVars, isPending: isDeleting } = useDeleteTask({
    mutation: { onSuccess: invalidate },
  });

  if (isLoading) return <p className="text-sm text-gray-500">Loading tasks…</p>;
  if (isError) return <p className="text-sm text-red-600">Could not load tasks.</p>;
  if (!tasks || tasks.length === 0) return <p className="text-sm text-gray-500">No tasks yet.</p>;

  return (
    <ul className="flex flex-col gap-2">
      {tasks.map((task) => (
        <TaskItem
          key={task.id}
          task={task}
          isMutating={
            (isSettingDone && setDoneVars?.id === task.id) ||
            (isDeleting && deleteVars?.id === task.id)
          }
          onToggleDone={(done) => setDone({ id: task.id, data: { done } })}
          onDelete={() => deleteTask({ id: task.id })}
        />
      ))}
    </ul>
  );
}
