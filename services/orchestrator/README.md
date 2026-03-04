Orchestrator

Responsibilities:
- Receive GitHub Issue webhooks and scheduled scan triggers
- Create job entries in Postgres
- Push job IDs to Redis Streams queue
- Expose simple admin API to requeue jobs / inspect status

Config:
- POSTGRES_URL
- REDIS_URL
- GITHUB_TOKEN
