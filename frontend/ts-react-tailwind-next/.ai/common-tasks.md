# Common tasks

## Add a new endpoint's UI

1. Add the operation to `openapi/task-api.yaml`, `npm run generate`.
2. If it's a query: call the new `use<Name>` hook directly in an
   organism. If it's a mutation: call it, invalidate the relevant
   `get<Name>QueryKey()` on success — copy `CreateTaskForm.tsx`.
3. Compose atoms/molecules for any new UI; add a new atom only if
   nothing existing fits (check `components/atoms/` first).
4. Add a test — pure props test if it's a molecule/atom, a mocked-fetch
   `QueryClientProvider` test if it's an organism (see `.ai/testing.md`).

## Point this frontend at a different backend

Set `NEXT_PUBLIC_API_BASE_URL` — any of `backend/go-fiber`,
`backend/java-spring`, `backend/quarkus` expose the same `/tasks` REST
shape this client is generated from. Rebuild (it's a build-time var,
see architecture.md#config).

## Add a new resource end to end

Copy the `task-api.yaml` paths/schemas shape for the new resource,
`npm run generate`, then copy `TaskItem`/`TaskList`/`CreateTaskForm`'s
shape (molecule + organism + organism) for the new resource's UI.
