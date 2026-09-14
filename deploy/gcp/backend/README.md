# GCP backend deploy — Cloud Run

Container deploy to **Cloud Run**, one workflow per environment, authenticated by
**Workload Identity Federation** (no service-account key). Works for any
containerized backend boilerplate (`go-fiber`, `java-spring`, `quarkus`) — it
deploys the image your `Dockerfile` builds.

## Copy in

Copy the three files in `workflows/` into your repo's `.github/workflows/`,
renamed with your boilerplate id:

```
<id>-deploy-stage.yml
<id>-deploy-preprod.yml
<id>-deploy-prod.yml
```

Keep the `stage` / `preprod` / `prod` keyword in both the file name and the
workflow `name:` so TaskTrooper maps each to the right deploy slot.

## One-time GCP setup

1. **Artifact Registry** Docker repo in your region.
2. **Cloud Run service** per environment (or let the first deploy create it).
3. **Workload Identity Federation**: a pool + provider trusting this GitHub repo,
   and a deploy service account with `roles/run.admin`,
   `roles/artifactregistry.writer`, and `roles/iam.serviceAccountUser`.

## GitHub configuration

Create three environments — `stage`, `preprod`, `production` — and set these
**vars** on each (values differ per environment):

| Var | Example |
|-----|---------|
| `GCP_PROJECT_ID` | `my-proj-stage` |
| `GCP_REGION` | `europe-west1` |
| `GCP_AR_REPO` | `apps` (Artifact Registry repo) |
| `GCP_WIF_PROVIDER` | `projects/123/locations/global/workloadIdentityPools/gh/providers/gh` |
| `GCP_DEPLOY_SA` | `deployer@my-proj-stage.iam.gserviceaccount.com` |
| `SERVICE_NAME` | `myapi` |
| `APP_DIR` | `.` (repo-root app) or `backend/go-fiber` (monorepo) |
| `HEALTH_PATH` | `/health` (default) |
| `MIGRATE_COMMAND` | `go run ./cmd/migrate up` (optional; empty = skip) |

**Secrets** (only if you run migrations): `DATABASE_URL` per environment. With
WIF there are no cloud credential secrets.

## Notes

- Images are tagged with the commit SHA — never `:latest` — so rollback is
  deterministic (`gcloud run services update-traffic --to-revisions`).
- `stage` also auto-deploys on push to `main`; `preprod`/`prod` are
  dispatch-only. `prod` is what `trigger_release` fires.
- Worker (no HTTP)? use `deploy/gcp/worker`. Static site? use `deploy/gcp/web`.
