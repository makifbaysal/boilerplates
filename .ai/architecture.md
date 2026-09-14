# Shared conventions

Applies to every boilerplate in this repo unless its own `.ai/architecture.md`
overrides it for stack-specific reasons.

## Repo shape

Each boilerplate is self-contained: nothing outside its own folder should be
required to build, test, or run it. No shared root-level dependency manifest,
no root-level source. The root only holds docs, the catalog, and per-stack CI
workflow files (which are path-scoped to one folder each).

## CI pattern

One workflow file per boilerplate at `.github/workflows/<id>-ci.yml`, id
matching the `catalog.yaml` entry. Trigger:

```yaml
on:
  push:
    paths: ["<path>/**", ".github/workflows/<id>-ci.yml"]
  pull_request:
    paths: ["<path>/**", ".github/workflows/<id>-ci.yml"]
```

Jobs: `build`, `lint`, `test` at minimum, run with `working-directory:
<path>`. No shared reusable-workflow indirection for now — one job list per
file, copy-paste is fine at this scale and keeps each workflow readable
standalone.

## Docker

Multi-stage builds, non-root final user, no secrets baked into layers. Base
images pinned to a specific tag (not `latest`). Prefer `alpine` for the
final stage over `distroless`/`debian` where the ecosystem has good musl
support (Go, Node) — smaller image, still has a shell for debugging.

## Testing

Every CRUD example ships with tests that actually exercise it (unit +
one integration/handler-level test), not placeholder assertions.

## Commits

Conventional Commits (`feat:`, `fix:`, `chore:`, ...). One boilerplate per
commit/PR when adding a new stack.

## Code generation

Wired into the build, not just documented:

| stack | tool |
|---|---|
| go-fiber, go-kafka-worker | mockery (interfaces → mocks) |
| java-spring | Lombok + MapStruct (mappers) + Mockito |
| quarkus | Quarkus codegen extensions (Panache) + Mockito |
| ts-react-tailwind-next | orval (OpenAPI → typed client + React Query hooks) |
| flutter | build_runner + freezed/json_serializable |
| android-native | Hilt (DI codegen) + KSP |
| ios-native | Swift macros / sourcery where applicable |

## Catalog

[.ai/catalog.yaml](catalog.yaml) is the source of truth for what exists.
Every new boilerplate PR must add its entry there in the same commit that
adds the folder — never land the folder without the catalog entry.
