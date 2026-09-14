import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { TaskList } from "./TaskList";
import type { Task } from "@/lib/api/generated/taskAPI.schemas";

// Runs the real generated hooks (useListTasks/useSetTaskDone/useDeleteTask)
// against a mocked fetch — the frontend equivalent of the backend
// boilerplates' "real adapter over a real-ish backing" adapter tests.
function renderWithQueryClient(ui: React.ReactElement) {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);
}

const task: Task = {
  id: "1",
  title: "buy milk",
  done: false,
  created_at: "2026-01-01T00:00:00Z",
  updated_at: "2026-01-01T00:00:00Z",
};

describe("TaskList", () => {
  beforeEach(() => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (url: string, options?: RequestInit) => {
        const method = options?.method ?? "GET";
        if (url.endsWith("/tasks") && method === "GET") {
          return new Response(JSON.stringify([task]), { status: 200 });
        }
        if (url.endsWith("/tasks/1") && method === "DELETE") {
          return new Response(null, { status: 204 });
        }
        throw new Error(`unmocked fetch: ${method} ${url}`);
      }),
    );
  });

  afterEach(() => vi.unstubAllGlobals());

  it("renders tasks returned by the API", async () => {
    renderWithQueryClient(<TaskList />);

    expect(await screen.findByText("buy milk")).toBeInTheDocument();
  });

  it("shows the empty state once tasks are gone", async () => {
    renderWithQueryClient(<TaskList />);
    await screen.findByText("buy milk");

    await userEvent.click(screen.getByRole("button", { name: "Delete" }));

    await waitFor(() => expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining("/tasks/1"),
      expect.objectContaining({ method: "DELETE" }),
    ));
  });
});
