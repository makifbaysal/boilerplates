import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { TaskItem } from "./TaskItem";
import type { Task } from "@/lib/api/generated/taskAPI.schemas";

const task: Task = {
  id: "1",
  title: "buy milk",
  done: false,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

describe("TaskItem", () => {
  it("calls onToggleDone with the new checked state", async () => {
    const onToggleDone = vi.fn();
    render(<TaskItem task={task} onToggleDone={onToggleDone} onDelete={vi.fn()} isMutating={false} />);

    await userEvent.click(screen.getByRole("checkbox"));

    expect(onToggleDone).toHaveBeenCalledWith(true);
  });

  it("calls onDelete when the delete button is clicked", async () => {
    const onDelete = vi.fn();
    render(<TaskItem task={task} onToggleDone={vi.fn()} onDelete={onDelete} isMutating={false} />);

    await userEvent.click(screen.getByRole("button", { name: "Delete" }));

    expect(onDelete).toHaveBeenCalledOnce();
  });

  it("disables controls while a mutation is in flight", () => {
    render(<TaskItem task={task} onToggleDone={vi.fn()} onDelete={vi.fn()} isMutating={true} />);

    expect(screen.getByRole("checkbox")).toBeDisabled();
    expect(screen.getByRole("button", { name: "Delete" })).toBeDisabled();
  });

  it("renders a done task with strikethrough styling", () => {
    render(<TaskItem task={{ ...task, done: true }} onToggleDone={vi.fn()} onDelete={vi.fn()} isMutating={false} />);

    expect(screen.getByText("buy milk")).toHaveClass("line-through");
  });
});
