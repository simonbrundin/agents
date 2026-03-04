# Product Roadmap

## Phase 1: MVP

- Auto-start one agent per GitHub issue
- Worktree/branch automation + PR creation
- Basic implementation loop: code change, unit tests, run CI
- Ephemeral PR environment for manual testing
  - Namespace convention: `pr-<pr-number>-<short-sha>` or `agent-<issue>-<id>` for previews
  - Use `tilt` for local developers and GitHub Actions for remote ephemeral deploys
  - Notify PR with ephemeral URL and cost estimate via `ntfy.sh` (primary). Example: `curl -d "Ephemeral ready: $URL" ntfy.sh/your-topic`
- Notifications to GitHub (issues/PR comments) and Slack (optional)
- Basic cost reporting per issue
- Security model: agents use isolated accounts and environments

## Phase 2: Post-Launch

- Agents suggest new issues periodically (configurable cadence)
- Model-priority list and provider-fallback when credits exhaust
- Automatic provider switching and context migration between providers
- Browser-driven UI testing agents (headless or real browsers in isolated env)
- Multi-agent orchestration dashboard with cost and status tracking
- Deeper observability and alerting (Prometheus/Grafana)
- Advanced policies for question routing and when to ask humans (in issues or Slack)
