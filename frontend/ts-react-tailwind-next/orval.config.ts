import { defineConfig } from "orval";

// Generates a typed fetch client + React Query hooks from
// openapi/task-api.yaml. Run `npm run generate` after editing the spec
// — see .ai/common-tasks.md. Output is committed (same call as
// go-fiber's mockery/mapstruct-generated code): it's small, reviewable,
// and this way `npm install` alone is enough to build, no codegen step
// required in CI.
export default defineConfig({
  taskApi: {
    input: "./openapi/task-api.yaml",
    output: {
      mode: "tags-split",
      target: "./src/lib/api/generated",
      client: "react-query",
      httpClient: "fetch",
      override: {
        mutator: {
          path: "./src/lib/api/mutator.ts",
          name: "customFetch",
        },
        // Our mutator (mutator.ts) returns the parsed body directly, not
        // a {data, status, headers} wrapper — this must stay false or
        // the generated types describe a wrapper that doesn't exist at
        // runtime (found by `next build`'s type check: `tasks.length`
        // failed because the inferred type was the wrapper).
        fetch: {
          includeHttpResponseReturnType: false,
        },
      },
    },
  },
});
