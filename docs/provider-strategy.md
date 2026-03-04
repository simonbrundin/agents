# Provider Strategy

Priority list (MVP):
1. `openrouter` (primary)
2. `copilot` (secondary; enabled when credential provided)

Behaviors:
- Provider Manager attempts providers in order until a successful response is obtained or list exhausted.
- Track per-provider metrics: calls, errors, cost, latency.
- Implement preflight cost estimation; abort if estimated cost > per-job-budget.

Credentials:
- `OPENROUTER_API_KEY`
- `COPILOT_API_KEY` (optional)

Costs:
- Per-job limit: 2 SEK (configurable)
- Daily budget guard: configured per deployment

Adapters:
- `adapters/openrouter` - REST client to OpenRouter
- `adapters/copilot` - stubbed adapter to be enabled when API key available
