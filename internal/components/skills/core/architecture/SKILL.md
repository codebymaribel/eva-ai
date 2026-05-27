---
name: architecture
description: >
  Stack decisions, project structure, patterns, state management,
  and technical conventions. Load before making structural changes.
sdd_phases: [scan, schematics, execute, debrief]
trigger: "stack, patterns, architecture, structure, state, conventions, module, dependency"
---

# Architecture

## Stack
<!-- Populated by eva init — describe the detected stack here -->
_Detected during project initialization. Update this section with your stack details._

## Project Structure
How the project is organized — folder layout, module boundaries, entry points.

### Convention: folder purpose
Each folder has a single responsibility. If a file doesn't match the folder's
purpose, it belongs somewhere else.

```
✅ correct: src/auth/middleware.go handles auth middleware
❌ incorrect: src/auth/utils.go has generic string helpers
```

## Patterns
How things are done in this project. Concrete examples over abstract rules.

### Dependency direction
Dependencies point inward. Outer layers depend on inner layers, never the reverse.

### State management
Where state lives, how it flows, and who owns it.

### Error handling
How errors are created, propagated, and logged.

## Rules
- Follow existing patterns before introducing new ones
- Every new dependency must be justified with a reason
- If two valid approaches exist, present both with tradeoffs
- New modules need a clear boundary and a single responsibility

## Do not
- Create circular dependencies
- Put business logic in infrastructure code
- Mix concerns within a single module
- Introduce a new pattern without checking if one already exists

## References
- Link to architecture decision records (ADRs) if they exist
- Link to module boundary docs
