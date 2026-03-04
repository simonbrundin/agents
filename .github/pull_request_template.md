# Pull Request Template

## Summary
A short description of what this PR changes.

## Related Issue
Closes #<issue-number>

## What I changed
- 

## Ephemeral preview
- Ephemeral namespace: `pr-<pr-number>-<short-sha>`
- Preview URL: (auto-populated by CI; example: `https://preview.example.com/pr-<pr-number>-<short-sha>`)
- Preview status: CI check `preview/health` (pass/fail). CI will also post a comment with the preview URL and a short smoke-test summary when available.
- Cost estimate for changes: 

## Checklist (agent-driven PRs should ensure these are automated where possible)
- [ ] Tests added for new behavior
- [ ] Existing tests pass
- [ ] Linter/format checks passed
- [ ] Ephemeral preview deployed and smoke-tested
- [ ] Security checks (secrets, permissions) verified

## Agent metadata
- Agent ID: 
- Actions performed: (e.g., created branch, ran tests, pushed changes)
- Model/provider used: 

## Notes for reviewer
- How to test manually
- Any known limitations
