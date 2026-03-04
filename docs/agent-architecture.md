# Agent Platform Architecture

This document describes the high-level architecture for the autonomous agent platform MVP.

## Components

- Orchestrator
  - Receives GitHub Issue webhooks and scheduled scan triggers
  - Creates jobs and pushes to queue (Redis Streams)
  - Persists job metadata in Postgres

- Worker
  - Short-lived Kubernetes Job/Pod
  - Clones repo, creates git worktree + branch
  - Runs repo tests, asks Provider Manager for code patches
  - Applies patches, runs tests, commits, pushes, creates PR
  - Triggers ephemeral environment deployment for PR

- Provider Manager
  - Adapter layer for LLM providers (OpenRouter, Copilot)
  - Tracks usage, costs and fallbacks
  - Single API for generation requests

- Context Store
  - Redis for short-term checkpoints
  - Postgres + pgvector for embeddings and long-term context

- Ephemeral Environments
  - Use existing kustomize `template-ephemeral` overlay to spin up review apps per PR

- Observability
  - Prometheus metrics, Grafana dashboards, Loki for logs

## Deployment

- All components run in Kubernetes; ArgoCD manages deployments via GitOps.
- Use service accounts with least privilege; secrets stored in Vault or SealedSecrets.
