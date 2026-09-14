@AGENTS.md

# ts-react-tailwind-next boilerplate

Next.js (App Router) + TypeScript + Tailwind + TanStack Query frontend,
Atomic Design component layering, a typed API client + hooks generated
by orval from `openapi/task-api.yaml` (the same Task CRUD shape as the
`backend/*` boilerplates — point `NEXT_PUBLIC_API_BASE_URL` at any of
them).

Read [.ai/architecture.md](.ai/architecture.md) before changing
structure — it explains why `fetch.includeHttpResponseReturnType: false`
in `orval.config.ts` is load-bearing (caught by `next build`'s type
check, not a style preference). [.ai/common-tasks.md](.ai/common-tasks.md),
[.ai/coding-standards.md](.ai/coding-standards.md), and
[.ai/testing.md](.ai/testing.md) cover the rest.

## Commands

```
npm run dev         # :3000
npm run build        # production build (type-checks too)
npm run lint          # eslint
npm test               # vitest run
npm run generate        # regenerate src/lib/api/generated/ from openapi/task-api.yaml
```
