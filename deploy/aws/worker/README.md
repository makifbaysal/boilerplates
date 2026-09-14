# AWS worker deploy — ECS Fargate service

Deploys a **long-running consumer** as an ECS Fargate service with **no load
balancer**. Auth via GitHub OIDC → IAM role. Copy `workflows/` into
`.github/workflows/` as `<id>-deploy-<env>.yml` (keep the env keyword). The
workflow renders `task-def.json` (in this folder — copy it next to your app and
edit it) with the freshly pushed image, then waits for service stability, which
is the smoke test: a crash-looping task never stabilizes and fails the deploy.

## One-time setup

ECR repo · an ECS cluster + Fargate **service** (no target group) per env, with
a task execution role and the worker's runtime IAM · OIDC provider + deploy IAM
role (ECR push, `ecs:UpdateService`, `ecs:RegisterTaskDefinition`,
`iam:PassRole`).

## Environment vars (per env)

| Var | Example |
|-----|---------|
| `AWS_REGION` | `eu-central-1` |
| `AWS_ROLE_ARN` | `arn:aws:iam::123:role/gh-deploy` |
| `ECR_REPOSITORY` | `orders-worker` |
| `ECS_CLUSTER` | `workers-stage` |
| `ECS_SERVICE` | `orders-worker` |
| `CONTAINER_NAME` | `worker` (matches `task-def.json`) |
| `APP_DIR` | `.` or `worker/go-kafka-worker` |
| `MIGRATE_COMMAND` | optional |

Secret `DATABASE_URL` per env only if migrating. Broker creds go in the task
definition's `secrets` (from SSM/Secrets Manager), not here.
