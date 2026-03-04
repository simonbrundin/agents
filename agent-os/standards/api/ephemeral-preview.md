# Ephemeral Preview Standard

Rule: Create an ephemeral preview environment for every Pull Request using the namespace pattern `pr-<pr-number>-<short-sha>`. Run automated smoke tests against the preview and report the result to the PR. Tear down the preview when the PR is closed.

Why
- Fast feedback: reviewers can test changes in a realistic environment.
- Traceability: namespace naming ties environment to the PR.
- Safety and cost control: automatic teardown prevents resource leakage.

When to create
- Automatically create ephemeral previews for all PRs targeting the repository.

Naming
- Use the namespace pattern: `pr-<pr-number>-<short-sha>`
  - `<pr-number>`: the GitHub PR number
  - `<short-sha>`: first 7 characters of the commit SHA used in the preview
  - Example: `pr-123-1a2b3c4`

Build & deploy
- CI job builds the preview image/artifact and deploys to an isolated environment (k8s namespace, preview cluster, or cloud preview service).
- Artifacts must be tagged with PR number and short SHA.

Smoke tests
- Run a small, fast smoke test suite against the preview (basic routing, auth, API health, key user flows).
- Smoke tests should complete within a short, configurable timeout (e.g., 5 minutes).
- Report results back to the PR as a status check or comment with pass/fail and a short log.

Teardown / Retention
- Destroy the preview when the PR is closed (merged or closed without merge).
- If an automated teardown fails, alert the PR author via comment and retry.

Access & Security
- Previews should use test/staging credentials or short-lived tokens, never production secrets.
- Network access should be restricted according to the repository's security posture (e.g., IP allowlist for internal services).

Cost & Quotas
- Keep preview resources lightweight (minimal replicas, scaled-down databases) to control costs.
- If quotas are reached, the system should fail gracefully and notify the PR.

Examples
- Ephemeral namespace for PR 45 with short SHA `d34db33`: `pr-45-d34db33`

Related
- `.github/pull_request_template.md` (ephemeral preview fields expected)

Guidance for failures
- If smoke tests fail, post a clear summary with log excerpts to the PR and mark the preview as failing.
- If the deployment cannot be created, post a comment linking to CI logs and suggest next steps.
