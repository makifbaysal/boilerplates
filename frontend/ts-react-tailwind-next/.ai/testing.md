# Testing

Vitest + React Testing Library + `@testing-library/user-event`.

- **Molecules/atoms**: pure props-in/callbacks-out tests — see
  `TaskItem.test.tsx`. No network, no providers needed.
- **Organisms**: render under a real `QueryClientProvider` with
  `global.fetch` mocked via `vi.stubGlobal("fetch", ...)` — this
  exercises the *real* generated hooks and `mutator.ts`, not a mocked
  hook layer, the same "real adapter, mocked only at the true I/O
  boundary" principle the backend boilerplates use. See
  `TaskList.test.tsx`.
- `npm test` runs everything once (CI mode); there's no watch script
  wired up by default — add `"test:watch": "vitest"` if you want one.
