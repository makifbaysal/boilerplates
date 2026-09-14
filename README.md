# boilerplates

Ready-to-use project skeletons for TaskTrooper — one per stack, each
standalone and CI-verified. See [CLAUDE.md](CLAUDE.md) for the rules and
[.ai/catalog.yaml](.ai/catalog.yaml) for what's available.

| category | stacks |
|---|---|
| backend | go-fiber, java-spring, quarkus |
| frontend | ts-react-tailwind-next |
| mobile | android-native, ios-native, flutter |
| worker | go-kafka-worker |

## Operational baseline (every boilerplate)

| File | What it gives a copied project |
| --- | --- |
| `<boilerplate>/.github/workflows/ci.yml` | Standalone CI: build, lint/analyze, test, and a Sonar job that skips itself when `SONAR_TOKEN` is absent. It assumes the boilerplate folder is the repository root, so a copy has green CI without editing paths. The monorepo-scoped copies under `.github/workflows/<id>-ci.yml` keep this repo itself checked. |
| `<boilerplate>/sonar-project.properties` | Quality-gate configuration with the right sources/tests split and coverage report path per stack. |
| Rate limit | Backends limit inbound requests per client IP (health probes exempt); the worker limits its processing rate so a backlog replay cannot take the database down; web and mobile clients throttle outbound calls. |
| Circuit breaker | One breaker per outbound dependency, dependency-free implementations (`pkg/resilience`, `platform/resilience`, `lib/resilience`, `NetworkGuard`). Only transport errors and 5xx count as failures — a 404 is this request being wrong, not the dependency being down. |

Defaults are tuned to be invisible in normal use and are configurable per
environment (`rate_limit` / `app.rate-limit` config keys, env vars on the web
client). They exist so nobody ships a starter project with an unmetered
endpoint or a retry storm aimed at a dead dependency.

