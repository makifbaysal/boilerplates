# AWS web deploy — S3 + CloudFront (static)

Builds a static frontend, syncs it to **S3**, and invalidates **CloudFront**.
Auth via GitHub OIDC → IAM role. Copy `workflows/` into `.github/workflows/` as
`<id>-deploy-<env>.yml` (keep the env keyword).

**Server-rendered app (Next.js SSR, API routes)?** Don't use this — deploy the
container with `deploy/aws/backend` (App Runner/ECS) instead.

## One-time setup

An S3 bucket per env · a CloudFront distribution in front (OAC to the bucket) ·
OIDC provider + deploy IAM role (`s3:PutObject`/`s3:DeleteObject`/`s3:ListBucket`
on the bucket, `cloudfront:CreateInvalidation`).

## Environment vars (per env)

| Var | Example |
|-----|---------|
| `AWS_REGION` | `eu-central-1` |
| `AWS_ROLE_ARN` | `arn:aws:iam::123:role/gh-deploy` |
| `S3_BUCKET` | `stage-example-web` |
| `CLOUDFRONT_DISTRIBUTION_ID` | `E123ABC456DEF` |
| `SITE_HOST` | `stage.example.com` |
| `BUILD_DIR` | `dist` or `out` |
| `APP_DIR` | `.` or `frontend/ts-react-tailwind-next` |

Public build config (`NEXT_PUBLIC_*` / `VITE_*`) is baked at build time — set it
as per-env vars. Never put a secret in a public build var. `index.html` is
uploaded with `no-cache` so a deploy is picked up immediately; hashed assets can
cache long.
