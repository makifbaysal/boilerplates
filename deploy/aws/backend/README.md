# AWS backend deploy — App Runner

Container deploy to **AWS App Runner**, one workflow per environment,
authenticated by **GitHub OIDC → an IAM role** (no access keys). Works for any
containerized backend boilerplate — it deploys the image your `Dockerfile`
builds. Prefer **ECS Fargate**? see the note at the bottom.

## Copy in

Copy `workflows/` into your repo's `.github/workflows/`, renamed:

```
<id>-deploy-stage.yml
<id>-deploy-preprod.yml
<id>-deploy-prod.yml
```

Keep the `stage` / `preprod` / `prod` keyword in the file name and workflow
`name:`.

## One-time AWS setup

1. **ECR** repository for the image.
2. **App Runner service** per environment, source = that ECR repo.
3. **GitHub OIDC provider** (`token.actions.githubusercontent.com`) and an
   **IAM role** trusting this repo, with permissions for ECR push and
   `apprunner:UpdateService` / `apprunner:DescribeService`.

## GitHub configuration

Environments `stage`, `preprod`, `production`, each with these **vars**:

| Var | Example |
|-----|---------|
| `AWS_REGION` | `eu-central-1` |
| `AWS_ROLE_ARN` | `arn:aws:iam::123456789012:role/gh-deploy` |
| `ECR_REPOSITORY` | `myapi` |
| `APP_RUNNER_SERVICE_ARN` | `arn:aws:apprunner:eu-central-1:123:service/myapi/abc` |
| `SERVICE_URL` | `xxxx.eu-central-1.awsapprunner.com` (no scheme) |
| `APP_PORT` | `8080` (container port) |
| `APP_DIR` | `.` or `backend/go-fiber` |
| `HEALTH_PATH` | `/health` |
| `MIGRATE_COMMAND` | `go run ./cmd/migrate up` (optional) |

**Secrets** (only for migrations): `DATABASE_URL` per environment. With OIDC
there are no AWS credential secrets.

## Notes

- Images are SHA-tagged, never `:latest` — rollback = re-point the service at the
  previous tag (see the comment in `deploy-prod.yml`).
- `stage` auto-deploys on push to `main`; `preprod`/`prod` are dispatch-only.
  `prod` is what `trigger_release` fires.
- **ECS Fargate instead of App Runner**: replace the "Deploy to App Runner" step
  with `aws-actions/amazon-ecs-render-task-definition` +
  `aws-actions/amazon-ecs-deploy-task-definition` (`wait-for-service-stability:
  true`) and add `ECS_CLUSTER` / `ECS_SERVICE` / `ECS_TASK_FAMILY` vars.
- Worker (no HTTP)? use `deploy/aws/worker`. Static site? use `deploy/aws/web`.
