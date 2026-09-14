# GCP worker deploy — Cloud Run (worker mode)

Deploys a **long-running consumer** (Kafka/RabbitMQ) as an always-on Cloud Run
service — no public ingress, CPU always allocated, `min-instances=1`. Auth via
Workload Identity Federation. Copy `workflows/` into `.github/workflows/` as
`<id>-deploy-<env>.yml` (keep the env keyword).

## Container requirement

Cloud Run requires the container to listen on `$PORT`. A pure consumer has no
HTTP, so expose a tiny health/readiness server on `$PORT` alongside the consumer
loop (the `go-kafka-worker` boilerplate's `cmd` should start one). Without it the
revision never becomes Ready and the deploy fails.

## One-time setup

Artifact Registry repo · a Cloud Run **service** per env · WIF pool/provider +
deploy SA (`roles/run.admin`, `roles/artifactregistry.writer`,
`roles/iam.serviceAccountUser`).

## Environment vars (per `stage`/`preprod`/`production`)

| Var | Example |
|-----|---------|
| `GCP_PROJECT_ID` | `my-proj-stage` |
| `GCP_REGION` | `europe-west1` |
| `GCP_AR_REPO` | `apps` |
| `GCP_WIF_PROVIDER` | `projects/123/.../providers/gh` |
| `GCP_DEPLOY_SA` | `deployer@…gserviceaccount.com` |
| `SERVICE_NAME` | `orders-worker` |
| `APP_DIR` | `.` or `worker/go-kafka-worker` |
| `MIGRATE_COMMAND` | optional |

Secret `DATABASE_URL` per env only if migrating. No cloud key secrets (WIF).

Broker credentials / topic config are the worker's own env/secrets — wire them
onto the Cloud Run service (`--set-env-vars` / `--set-secrets`) in setup.
