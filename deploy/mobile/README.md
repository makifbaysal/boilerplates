# Mobile deploy — fastlane (TestFlight + Play)

Mobile ships to **Apple's and Google's stores**, so there is no GCP/AWS split.
Each env workflow builds a signed artifact and uploads it to a track:

| Env | iOS | Android |
|-----|-----|---------|
| `stage` | TestFlight (internal) | Play `internal` |
| `preprod` | TestFlight (external group) | Play `beta` (closed) |
| `prod` | App Store (submit for review) | Play `production` |

Copy `fastlane/` (Fastfile + Appfile) and `workflows/` into your app.
Rename the workflows to `<id>-deploy-<env>.yml` (keep the env keyword — `prod`
is what `trigger_release` dispatches). Add a `Gemfile` with `gem "fastlane"`.

## Signing (the hard part)

- **iOS** uses fastlane **match** — certs/profiles live in a *separate* private
  git repo, decrypted at CI time. Never commit certificates to the app repo.
- **Android** signs with an upload keystore decoded from a base64 secret at
  runtime.

## GitHub configuration (per env)

**Vars:** `APP_DIR`, `MATCH_GIT_URL`, `ASC_KEY_ID`, `ASC_ISSUER_ID`,
`IOS_SCHEME`, `IOS_BUNDLE_ID`, `TF_EXTERNAL_GROUP`, `ANDROID_PACKAGE_NAME`,
`ANDROID_KEY_ALIAS`.

**Secrets (these ARE real secrets — no OIDC path to the stores):**
`MATCH_PASSWORD`, `ASC_KEY_P8` (base64 of the `.p8`), `PLAY_JSON_KEY`
(service-account JSON), `ANDROID_KEYSTORE_B64`, `ANDROID_KEYSTORE_PASSWORD`.

## Notes

- iOS builds run on `macos-14`; Android on `ubuntu-latest`.
- Bump the build/version number every upload or the store rejects it.
- A successful upload to the track is the pass signal; a failed/ rejected upload
  fails the job and sends the task to `need_revision`.
