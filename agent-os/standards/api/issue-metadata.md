# Issue Metadata Standard

Rule: When an issue is created or suggested by an automated agent, add the `agent-suggested` label and include agent metadata in the issue body.

Why
- Makes it easy to filter and review agent-originated suggestions.
- Provides traceability (which agent, why, and confidence).

Required fields (when agent-created/suggested)
- `Agent ID:` A stable identifier for the agent (e.g., `assistant-v1`).
- `Trigger:` Why the issue was created (e.g., `issue-created`, `suggested`, `agent-action`).
- `Confidence:` A numeric or textual confidence level (e.g., `0.78` or `medium`).
- `Ephemeral namespace:` (optional) If the agent created test resources, include the ephemeral namespace.

Labeling
- Add GitHub label: `agent-suggested`.
- Do NOT modify issue title automatically unless human-reviewed.

When to add metadata
- Agents should append metadata when they create or update an issue.
- Human-created issues do not need these fields unless an agent later modifies them.

Examples
- Feature request body snippet:

```
**Agent metadata (if the suggestion came from an agent):**
- Agent ID: assistant-v1
- Confidence: 0.86
- Suggested cadence: weekly
```

- Bug report body snippet:

```
**Agent metadata (if the issue was created/modified by an agent):**
- Agent ID: assistant-v1
- Trigger: issue-created
- Ephemeral namespace: pr-1234-abcde
```

Guidance for reviewers
- Review the suggestion on merits; remove the `agent-suggested` label after triage if accepted or not relevant.
- If metadata is missing but suspected agent-origin, add the `agent-suggested` label and fill the fields if possible.

Exceptions
- Small formatting or grammar edits by agents do not require `agent-suggested` labeling.

Related
- `.github/ISSUE_TEMPLATE/feature_request.md`
- `.github/ISSUE_TEMPLATE/bug_report.md`
