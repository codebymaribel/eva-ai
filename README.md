# 🤖 eva-ai

## 🚧 THIS PROJECT IS UNDER CONSTRUCTION. DO NOT PULL 🚧

### Customize your agentic development workflow from one central command unit.

eva-ai is a CLI tool that injects your custom Skills, SDD workflow, MCP servers, and Persona into your AI coding agents — so every agent works exactly the way you think.

Inspired by [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai), built from scratch with a focus on developer experience.

```bash
# Initialize your project
eva init

# Inject into your AI agent
eva install --agent claude-code --preset full

# Add external skills
eva skill add https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md
```

---

## Why eva-ai?

Every AI coding agent has its own config format, its own folder, its own rules. Setting them up consistently across your machine is tedious and error-prone.

eva-ai solves that. Define your workflow once, inject it everywhere.

---

## Supported Agents

| Agent | Config Location | Status |
|---|---|---|
| **Claude Code** | `~/.claude/CLAUDE.md` | ✅ Supported |
| **Cursor** | `~/.cursor/rules/.cursorrules` | 🚧 Adapter pending |
| **GitHub Copilot** | `.github/copilot-instructions.md` | 🚧 Adapter pending |
| **OpenCode** | `~/.config/opencode/AGENTS.md` | 🚧 Adapter pending |
| **Windsurf** | `~/.codeium/windsurf/memories/global_rules.md` | 🚧 Adapter pending |

---

## Components

| Component | Description | Status |
|---|---|---|
| `sdd` | Spec-Driven Development — plan before you code | ✅ |
| `skills` | Curated coding patterns and best practices | ✅ |
| `mcp` | MCP server registration (Context7, etc.) | 🚧 |
| `persona` | Agent behavior, teaching style, and tone | 🚧 |

---

## Installation

```bash
# Go install
go install github.com/codebymaribel/eva-ai/cmd/directive@latest

# Or build from source
git clone https://github.com/codebymaribel/eva-ai.git
cd eva-ai
go build -o bin/eva ./cmd/directive
```

---

## Usage

### `eva init` — Set up your project

Scans your codebase, detects the stack, and creates `.eva/` with core skills and project-specific skills.

```bash
# Full init — scan project and create all skills
eva init

# Init without skills — only SDD structure
eva init --no-skills
```

**What it creates:**

```
.eva/
├── README.md                     ← agent entry point
├── memory.md                     ← mission context (gitignored)
├── phases/                       ← SDD phase outputs (gitignored)
├── skills/
│   ├── README.md                 ← skill index with triggers
│   ├── architecture/SKILL.md     ← core: stack, patterns, structure
│   ├── testing/SKILL.md          ← core: test conventions
│   ├── git-workflow/SKILL.md     ← core: git conventions
│   └── domain/SKILL.md           ← generated: business logic
└── docs/                         ← feature documentation
```

### `eva install` — Inject into agents

```bash
# Full stack for Claude Code
eva install --agent claude-code --preset full

# Pick specific components
eva install --agent claude-code --component sdd,skills

# Preview without making changes
eva install --agent claude-code --dry-run
```

**Presets:**

| Preset | Components included |
|---|---|
| `full` | sdd, skills, mcp, persona |
| `minimal` | persona |
| `custom` | you pick with `--component` |

### `eva skill add` — Add external skills

```bash
# From a local file
eva skill add ./path/to/SKILL.md

# From a GitHub URL
eva skill add https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md

# From a raw URL
eva skill add https://raw.githubusercontent.com/user/repo/main/SKILL.md
```

**What happens:**
1. Downloads/copies the SKILL.md to `.eva/skills/<name>/`
2. Updates `.eva/skills/README.md` with the new entry
3. Injects the skill into existing agent configs (CLAUDE.md, .cursorrules, etc.)

The agent sees the skill rules from message one — no manual invocation needed.

---

## How it works

```
eva init
  ├── 1. Scans project (language, framework, patterns, tooling)
  ├── 2. Creates .eva/ directory structure
  ├── 3. Writes core skills (architecture, testing, git-workflow)
  └── 4. Generates project-specific skills from scan results

eva install --agent claude-code --component sdd,skills
  ├── 1. Detects your OS and platform
  ├── 2. Validates Claude Code is installed
  ├── 3. Reads existing CLAUDE.md
  ├── 4. Renders selected components
  ├── 5. Appends to CLAUDE.md
  └── 6. Done — agent sees new config next session

eva skill add <url>
  ├── 1. Downloads SKILL.md (normalizes GitHub URLs)
  ├── 2. Saves to .eva/skills/<name>/SKILL.md
  ├── 3. Refreshes .eva/skills/README.md
  └── 4. Injects into existing agent configs
```

---

## Platform support

| OS | Status |
|---|---|
| macOS | ✅ Supported |
| Linux | ✅ Supported |
| Windows | 🚧 Coming soon |

---

## Development

```bash
# Clone
git clone https://github.com/codebymaribel/eva-ai.git
cd eva-ai

# Install dependencies
go mod tidy

# Run tests
go test ./... -v

# Build binary
go build -o bin/eva ./cmd/directive

# Dry-run example
./bin/eva install --agent claude-code --dry-run
```

### Project structure

```
eva-ai/
├── cmd/directive/              # CLI entrypoint
├── internal/
│   ├── agents/                 # Agent interface + adapters
│   │   ├── agent.go            # Agent + InjectableAgent interfaces
│   │   └── claudecode/         # Claude Code adapter
│   ├── app/                    # Cobra command wiring
│   ├── cli/                    # Flags, validation, InstallOptions
│   ├── components/
│   │   ├── sdd/                # SDD component (orchestrator + 6 phases)
│   │   └── skills/             # Skills component (core + template)
│   ├── scanner/                # Project stack and pattern detection
│   └── system/                 # OS and distro detection
└── Makefile
```

---

## Roadmap

- [x] Phase 1 — CLI foundation, OS detection, flag validation
- [x] Phase 2 — Agent adapters (Claude Code)
- [x] Phase 3 — Components (SDD, Skills)
- [ ] Phase 3 — Components (MCP, Persona)
- [ ] Phase 2 — Agent adapters (Cursor, Copilot, OpenCode, Windsurf)
- [ ] Phase 4 — Pipeline, backup, and rollback
- [ ] Phase 5 — Interactive TUI
- [ ] Phase 6 — Distribution (Homebrew, GoReleaser)

---

## License

MIT — built with 🤖 by [@codebymaribel](https://github.com/codebymaribel)
