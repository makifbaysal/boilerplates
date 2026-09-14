# GCP web deploy — GCS + Cloud CDN (static)

Builds a static frontend and publishes it to a **GCS bucket behind Cloud CDN**.
Auth via Workload Identity Federation. Copy `workflows/` into
`.github/workflows/` as `<id>-deploy-<env>.yml` (keep the env keyword).

**Server-rendered app (Next.js SSR, API routes)?** Don't use this — deploy the
container with `deploy/gcp/backend` (Cloud Run) instead.

## One-time setup

A GCS bucket per env (website config, public read via the CDN backend) · an
HTTPS load balancer with a CDN-enabled backend bucket · WIF pool/provider + a
deploy SA with `roles/storage.admin` and `roles/compute.loadBalancerAdmin` (for
cache invalidation).

## Environment vars (per env)

| Var | Example |
|-----|---------|
| `GCP_PROJECT_ID` | `my-proj-stage` |
| `GCP_WIF_PROVIDER` | `projects/123/.../providers/gh` |
| `GCP_DEPLOY_SA` | `deployer@…gserviceaccount.com` |
| `GCS_BUCKET` | `stage.example.com` |
| `GCP_URL_MAP` | `web-lb` (leave empty to skip CDN invalidation) |
| `SITE_HOST` | `stage.example.com` (for the smoke check) |
| `BUILD_DIR` | `dist` (Vite) or `out` (`next export`) |
| `APP_DIR` | `.` or `frontend/ts-react-tailwind-next` |

Public build config (`NEXT_PUBLIC_*` / `VITE_*`) is **baked at build time** — set
it as per-env vars; the build job reads that env's values. Never put a secret in
a public build var. **Firebase Hosting** is a fine alternative to GCS+CDN (built
-in CDN + rollback) — swap the upload/invalidate steps for
`firebase deploy --only hosting`.
