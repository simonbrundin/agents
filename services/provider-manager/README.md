Provider Manager

This service exposes a single endpoint `/v1/generate` that accepts generation requests and attempts them against a prioritized list of providers.

Adapters:
- openrouter
- copilot (stub)

Metrics:
- prometheus: provider_calls_total, provider_errors_total, provider_cost_total
