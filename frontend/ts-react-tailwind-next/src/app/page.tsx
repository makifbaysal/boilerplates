import { CreateTaskForm } from "@/components/organisms/CreateTaskForm";
import { TaskList } from "@/components/organisms/TaskList";

export default function Home() {
  return (
    <main className="mx-auto flex w-full max-w-md flex-1 flex-col gap-6 px-6 py-16">
      <h1 className="text-xl font-semibold text-gray-900">Tasks</h1>
      <CreateTaskForm />
      <TaskList />
    </main>
  );
}
