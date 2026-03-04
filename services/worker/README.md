Worker

Kubernetes Job that pulls a job from Redis and performs the following steps:
1. Clone repository
2. Create Git worktree and branch: `agent/<issue-id>/<slug>`
3. Run tests
4. Request code patch from Provider Manager
5. Apply patch, run tests
6. Commit & push branch
7. Create PR and post ephemeral deployment URL

Config:
- GITHUB_TOKEN
- POSTGRES_URL
- REDIS_URL
