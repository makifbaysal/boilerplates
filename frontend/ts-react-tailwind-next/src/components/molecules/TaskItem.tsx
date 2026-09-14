import { Button } from "@/components/atoms/Button";
import { Checkbox } from "@/components/atoms/Checkbox";
import type { Task } from "@/lib/api/generated/taskAPI.schemas";

export function TaskItem({
  task,
  onToggleDone,
  onDelete,
  isMutating,
}: {
  task: Task;
  onToggleDone: (done: boolean) => void;
  onDelete: () => void;
  isMutating: boolean;
}) {
  return (
    <li className="flex items-center gap-3 rounded-md border border-gray-200 px-3 py-2">
      <Checkbox
        checked={task.done}
        disabled={isMutating}
        onChange={(e) => onToggleDone(e.target.checked)}
        aria-label={`Mark "${task.title}" ${task.done ? "not done" : "done"}`}
      />
      <span className={`flex-1 text-sm ${task.done ? "text-gray-400 line-through" : "text-gray-900"}`}>
        {task.title}
      </span>
      <Button variant="danger" disabled={isMutating} onClick={onDelete}>
        Delete
      </Button>
    </li>
  );
}
