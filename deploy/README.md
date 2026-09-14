# Deploy boilerplates

CI/CD **deploy** templates, organized by `cloud / app-type`. The framework
boilerplates (`backend/`, `frontend/`, `mobile/`, `worker/`) ship a
`<id>-ci.yml` that lints/builds/tests. These add the **CD** half: environment
deploys to a real cloud.

Deploy config is cloud- and app-type-specific, **not** framework-specific — a
Cloud Run deploy is the same for `go-fiber`, `java-spring`, or `quarkus` (all
produce a container). So the deploy templates live here once, not duplicated
into every stack.

```
deploy/
  gcp/
    backend/   Cloud Run (container)
    worker/    Cloud Run job (no ingress)
    web/       Firebase Hosting / GCS+CDN (static)   — SSR? use gcp/backend
  aws/
    backend/   App Runner / ECS Fargate (container)
    worker/    ECS service (no ingress)
    web/       S3 + CloudFront (static)               — SSR? use aws/backend
  mobile/      fastlane → TestFlight + Play           (cloud-independent)
```

## How to use one

1. Pick the folder for your **cloud** and **app type**.
2. Copy its `workflows/deploy-stage.yml`, `deploy-preprod.yml`, `deploy-prod.yml`
   into your project's `.github/workflows/`, renamed to
   `<id>-deploy-stage.yml`, `<id>-deploy-preprod.yml`, `<id>-deploy-prod.yml`
   (keep the `stage` / `preprod` / `prod` keyword in the file name **and** the
   workflow `name:`).
3. Set the GitHub **environment vars/secrets** listed in that folder's README,
   and create the OIDC trust (GCP Workload Identity Federation / AWS IAM role)
   once.

## Why the naming matters

TaskTrooper reads your workflows and maps them onto its pipeline:

| Pipeline slot | Detected from |
|---------------|---------------|
| `validate` / `build` / `test` | a **job name** containing `lint`/`vet`, `build`/`docker`, `test` (in the `<id>-ci.yml`) |
| `stage_deploy` | a **workflow** whose name contains `stage` / `staging` / `preview` |
| `preprod_deploy` | a workflow whose name contains `preprod` / `pre-prod` |
| `prod_deploy` | a workflow whose name contains `prod` / `production` / `release` |

Every deploy workflow declares `on: workflow_dispatch` with **no required
inputs** — TaskTrooper's `trigger_release` dispatches with the branch ref only.

## Auth: OIDC, never stored keys

- **GCP:** GitHub OIDC → Workload Identity Federation → a deploy service
  account. No service-account JSON in secrets.
- **AWS:** GitHub OIDC → an IAM role via `aws-actions/configure-aws-credentials`.
  No access keys in secrets.
- **Mobile:** the exception — store signing material (App Store Connect API key,
  Play service-account JSON, keystore) as encrypted GitHub secrets; there is no
  OIDC path to the stores.

Non-sensitive config (project id, region, service name, bucket) goes in repo /
environment **vars**, not `secrets`.
