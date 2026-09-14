"use client";

import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";

import { Button } from "@/components/atoms/Button";
import { Input } from "@/components/atoms/Input";
import { useCreateTask, getListTasksQueryKey } from "@/lib/api/generated/default/default";

export function CreateTaskForm() {
  const [title, setTitle] = useState("");
  const [error, setError] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const { mutate, isPending } = useCreateTask({
    mutation: {
      onSuccess: () => {
        setTitle("");
        setError(null);
        queryClient.invalidateQueries({ queryKey: getListTasksQueryKey() });
      },
      onError: () => setError("Could not create task — title may be blank."),
    },
  });

  return (
    <form
      className="flex gap-2"
      onSubmit={(e) => {
        e.preventDefault();
        mutate({ data: { title } });
      }}
    >
      <Input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Buy milk"
        disabled={isPending}
        aria-label="New task title"
      />
      <Button type="submit" disabled={isPending}>
        Add task
      </Button>
      {error && <p className="text-sm text-red-600">{error}</p>}
    </form>
  );
}
