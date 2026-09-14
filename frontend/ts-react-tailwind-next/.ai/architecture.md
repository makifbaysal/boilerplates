# Architecture

Next.js App Router, TypeScript, Tailwind, TanStack Query. Atomic Design
component layering (see the root repo's own `apps/web` rule — this
boilerplate follows the same convention so it's a drop-in match if this
repo's own frontend rule applies to whatever project you copy it into):

```
src/
  app/
    layout.tsx       root layout, wraps children in Providers
    providers.tsx      QueryClientProvider (client component)
    page.tsx            assembles organisms — no business logic here
  components/
    atoms/            Button, Input, Checkbox — no app knowledge, pure props
    molecules/         TaskItem — composes atoms, still no data-fetching
    organisms/          TaskList, CreateTaskForm — the only files that call
                         generated hooks / own state
  lib/
    api/
      mutator.ts         the one fetch implementation every generated hook
                           calls through — base URL + headers live here
      generated/          orval output — committed, not gitignored
openapi/
  task-api.yaml        source of truth for the generated client
```

## Code generation: orval

`openapi/task-api.yaml` describes the same Task REST API as
`backend/go-fiber`/`backend/java-spring`/`backend/quarkus` — point
`NEXT_PUBLIC_API_BASE_URL` at any of them and this frontend works
unmodified. `npm run generate` (orval) regenerates
`src/lib/api/generated/` after editing the spec; the output is
committed (same call as the backends' mockery/MapStruct-generated code)
so `npm install` alone is enough to build — no codegen step required in
CI.

**`fetch.includeHttpResponseReturnType: false` in `orval.config.ts` is
load-bearing, not cosmetic.** Without it, orval's fetch client generates
types describing a `{data, status, headers}` wrapper regardless of what
the mutator actually returns — our `mutator.ts` returns the parsed body
directly, so the wrapped type was simply wrong (`next build`'s type
check caught this: `tasks.length` failed because the inferred type
had no `.length` under the mismatched wrapper).

## Data fetching

TanStack Query owns all server state — no `useEffect` fetch code
anywhere. Mutations invalidate `getListTasksQueryKey()` on success
rather than manually patching cache, keeping the organisms simple at
the cost of one extra request per mutation — fine at this scale.

## Config

`NEXT_PUBLIC_API_BASE_URL` (read once in `lib/api/mutator.ts`) is
inlined into the client bundle at **build time** — unlike the backend
boilerplates' runtime env vars, changing it means rebuilding the image,
not just restarting the container. Bake a separate image per
environment, or fall back to relative paths + a reverse proxy if you
need the same image across environments.
