# Tech Stack

## Frontend

Nuxt (Vue)

## Backend

Go

## Database

PostgreSQL

## ORM

Drizzle

## Deployment / Orchestration

Kubernetes

## Other

- CI: (e.g., GitHub Actions)
- Provider SDKs: OpenAI / other LLM providers
- Observability: Prometheus / Grafana (suggested)
- Messaging/Queueing: (e.g., RabbitMQ / NATS) if needed
- Ephemeral environment tooling for PR testing (e.g., ephemeral namespaces in Kubernetes)

## Dev & Ephemeral Environment Details

- Local & preview tooling: `tilt` with `local` and `kubernetes` modes (see project `plan/dev.nu` for example).
- Namespace convention: `pr-<pr-number>-<short-sha>` or `agent-<issue>-<id>` for ephemeral deploys.
- Dev workflow: `tilt up --namespace $namespace -- mode=<local|kubernetes>` for local development and previews.
- K8s login: check `kubectl cluster-info` and fallback to Teleport if cluster access is unavailable.
- Notifications: use `ntfy.sh` for quick dev notifications (primary). Example: `curl -d "http://$LINK" ntfy.sh/your-topic`.
  - Slack webhooks can be added later for richer team integration.
- CI: GitHub Actions for ephemeral deploys and tests (build, deploy to ephemeral namespace, run e2e, report URL back to PR).
- Observability: expose Prometheus metrics and Grafana dashboards in ephemeral environments when feasible.
