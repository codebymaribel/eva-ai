# SDD — Spec-Driven Development

Before writing any code, always follow this workflow in order.
Never skip phases. Never code before the spec is approved.
All phases are orchestrated by EVA — invoke `/eva [task]` to start.

## Phases

| Phase | Command | Description |
|---|---|---|
| 0. Scan | `/scan` | Understand the existing codebase (read-only) |
| 1. Directive | `/directive` | Define WHAT we build and WHY |
| 2. Schematics | `/schematics` | Define HOW we build it — technical design |
| 3. Sequence | `/sequence` | Break into atomic, ordered tasks |
| 4. Execute | `/execute` | Implement ONE task at a time |
| 5. Debrief | `/debrief` | Verify implementation and sync .eva/ |

## Rules

- **Always start with `/eva [task]`** — never invoke phases directly
- `/scan` must be approved before `/directive`
- `/directive` must be approved before `/schematics`
- `/schematics` must be approved before `/sequence`
- Each `/execute` session handles ONE task only
- `/debrief` runs after ALL tasks are complete
- If requirements change mid-flight → back to `/directive`