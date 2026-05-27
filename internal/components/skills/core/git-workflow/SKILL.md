---
name: git-workflow
description: >
  Git conventions — branching, commits, PRs, and release workflow.
  Load before committing, branching, or creating PRs.
sdd_phases: [execute, debrief]
trigger: "git, commit, branch, PR, pull request, merge, release, changelog, conventional commits"
---

# Git Workflow

## Branching Strategy
<!-- Populated by eva init — describe the detected branching model here -->
_Detected during project initialization. Update this section with your branching model._

### Convention: branch naming
Branch names follow the pattern: `type/short-description`

```
✅ correct: feat/user-auth, fix/login-redirect, refactor/api-client
❌ incorrect: my-stuff, fix-bug, maria-changes
```

## Commit Messages
Follow Conventional Commits format:

```
type(scope): short description

Optional body — explain WHY, not WHAT.

Optional footer — references to issues.
```

### Types
| Type | When to use |
|---|---|
| `feat` | New feature |
| `fix` | Bug fix |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `chore` | Build, CI, tooling changes |

### Rules
- Subject line: max 72 characters, imperative mood
- Body: wrap at 72 characters, explain the WHY
- Footer: reference issues with `Closes #123` or `Fixes #123`

```
✅ correct: feat(auth): add JWT refresh token rotation
❌ incorrect: added some auth stuff
```

## Pull Requests
- PR title follows the same Conventional Commits format
- PR description summarizes what changed and why
- Link to the issue or task that motivated the PR
- Keep PRs focused — one concern per PR

## Rules
- Never force push to shared branches
- Never commit directly to main/master
- Squash merge feature branches
- Delete branches after merge

## Do not
- Commit generated files, build artifacts, or secrets
- Use `git add .` without reviewing what's staged
- Mix unrelated changes in a single commit
- Rewrite history on shared branches

## References
- Link to the project's CONTRIBUTING.md if it exists
- Link to the project's CI/CD pipeline docs
