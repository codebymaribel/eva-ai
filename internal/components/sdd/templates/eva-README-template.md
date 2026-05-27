# .eva — AI Agent Knowledge Base

This directory contains everything an AI agent needs to work
on this project effectively.

---

## Quick Start — No SDD needed

Read the skills relevant to your task before writing any code:

| Task | Skills to load |
|---|---|
| New component | `.eva/skills/mobile-ui/SKILL.md` |
| API integration | `.eva/skills/api/SKILL.md` |
| State management | `.eva/skills/architecture/SKILL.md` |
| Navigation | `.eva/skills/navigation/SKILL.md` |
| Domain logic | `.eva/skills/domain/SKILL.md` |

## Available Skills

| Skill | Trigger |
|---|---|
| `architecture` | Stack, patterns, state, technical decisions |
| `domain` | Business logic, models, entities |
| `mobile-ui` | Components, styles, animations |
| `api` | HTTP client, auth, tokens |
| `navigation` | Routes, guards, tabs |

### Common combinations

| Task | Skills |
|---|---|
| New screen with API data | `mobile-ui` + `api` + `navigation` |
| New Zustand store | `architecture` + `domain` |
| Auth flow | `api` + `navigation` + `architecture` |

### Shared skills

| Skill | Trigger |
|---|---|
| `shared-skills/typescript` | Types, generics, Zod, strict mode |
| `shared-skills/react-native` | Performance, lists, animations, Expo |

---

## If you use SDD

The `.eva/phases/` folder contains outputs from each SDD mission phase.
Use the phase mapping below to know which skills each phase loads:

| Phase | Command | Skills loaded |
|---|---|---|
| 0 — Scan | `/scan` | `architecture`, `domain` |
| 1 — Directive | `/directive` | none — intent only |
| 2 — Schematics | `/schematics` | all relevant to the change |
| 3 — Sequence | `/sequence` | none — breakdown only |
| 4 — Execute | `/execute` | skills referenced in the task |
| 5 — Debrief | `/debrief` | all skills + phases |

---

## Directory Structure

```
.eva/
├── README.md                     ← this file
├── memory.md                     ← accumulated mission context (gitignored)
├── phases/                       ← SDD phase outputs (gitignored)
│   ├── scan.md
│   ├── directive.md
│   ├── schematics.md
│   ├── sequence.md
│   └── debrief.md
├── skills/                       ← project-specific skills (in git)
│   ├── README.md                 ← this file
│   ├── architecture/SKILL.md
│   ├── domain/SKILL.md
│   ├── mobile-ui/SKILL.md
│   ├── api/SKILL.md
│   └── navigation/SKILL.md
├── shared-skills/                ← general skills (in git)
│   ├── typescript/SKILL.md
│   └── react-native/SKILL.md
└── docs/                         ← feature documentation (in git, optional)
    └── example-feature.md
```

---

## Maintenance

- If a project convention changes → update the relevant `SKILL.md`
- If a new pattern emerges → add it to the right skill, don't create a new file
- If a feature is completed → document it in `.eva/docs/[feature].md`
- Run `eva sync` after manual changes to verify skills match the codebase
