# TaskTrooper Boilerplates

Monorepo of ready-to-use project skeletons, one per language/framework. Each
subfolder is a **standalone, buildable project** with its own dependency
manifest, tests, Dockerfile, CI workflow, `CLAUDE.md`, and `.ai/` docs.

This repo exists so agents (and humans) never scaffold a stack from scratch:
check [.ai/catalog.yaml](.ai/catalog.yaml) first, copy the matching folder,
rename, go. That lookup-before-generate step is the whole point — it saves
the tokens a full from-scratch scaffold would cost.

## Code comments

Do not add code comments unless truly necessary — a non-obvious invariant, a
workaround, or a WHY that isn't clear from the code itself. Never explain WHAT
the code does; well-named identifiers already do that. Applies to every
boilerplate in this repo.

## Layout

```
backend/    go-fiber, java-spring, quarkus
frontend/   ts-react-tailwind-next
mobile/     android-native, ios-native, flutter
worker/     go-kafka-worker
```

## Rules for every boilerplate

1. Must build, lint, and pass tests standalone (`cd <path> && <build cmd>`) —
   verified by its own path-filtered GitHub Actions workflow, not the root.
2. Must ship: a health/liveness endpoint or equivalent smoke check, one full
   CRUD example (so the pattern is visible, not just hello-world), tests for
   that CRUD example, and a Dockerfile.
3. Must include `<path>/CLAUDE.md` (short, points into `.ai/`) and
   `<path>/.ai/architecture.md` + `<path>/.ai/common-tasks.md` (walkthroughs
   for the things people actually ask for repeatedly — "add a new endpoint",
   "add a new background job", "add a new screen", etc).
4. Must include working code-gen tooling for that stack wired into the build
   (mocks, API clients, DI boilerplate — whatever is idiomatic there), not
   just a mention in docs.
5. Must be registered in [.ai/catalog.yaml](.ai/catalog.yaml) — id, path,
   type, language, framework, tags, description. This file is the machine-
   readable index agents search before generating anything.
6. Naming: folder = `<category>/<language-or-runtime>-<framework>`, kebab
   case, matches its `catalog.yaml` id.

See [.ai/architecture.md](.ai/architecture.md) for shared conventions
(commits, CI pattern, Docker base images, testing).
