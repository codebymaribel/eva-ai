# 🤖 eva-ai


## 🚧 THIS PROJECT IS UNDER CONSTRUCTION. DO NOT PULL 🚧


### Customize your agentic development workflow from one central command unit.

eva-ai is a CLI tool that injects your custom Skills, SDD workflow, MCP servers, and Persona into your AI coding agents — so every agent works exactly the way you think.

Inspired by [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai), built from scratch with a focus on developer experience.

```bash
eva install --agent claude-code,cursor --preset full
```

---

## Why eva-ai?

Every AI coding agent has its own config format, its own folder, its own rules. Setting them up consistently across your machine is tedious and error-prone.

eva-ai solves that. Define your workflow once, inject it everywhere.

---

## Supported Agents

| Agent | Config Location |
|---|---|
| **Claude Code** | `~/.claude/CLAUDE.md` |
| **Cursor** | `~/.cursor/rules/.cursorrules` |
| **GitHub Copilot** | `.github/copilot-instructions.md` |
| **OpenCode** | `~/.config/opencode/AGENTS.md` |
| **Windsurf** | `~/.codeium/windsurf/memories/global_rules.md` |

---

## Components

| Component | Description |
|---|---|
| `sdd` | Spec-Driven Development — plan before you code |
| `skills` | Your curated coding patterns and best practices |
| `mcp` | MCP server registration (Context7, etc.) |
| `persona` | Agent behavior, teaching style, and tone |

---

## Installation

```bash
# Homebrew (macOS) — coming soon
brew tap codebymaribel/tap
brew install eva-ai

# Go install
go install github.com/codebymaribel/eva-ai/cmd/directive@latest
```

---

## Usage

```bash
# Full stack — install everything into Claude Code and Cursor
eva install --agent claude-code,cursor --preset full

# Minimal — persona only
eva install --agent claude-code --preset minimal

# Pick specific components
eva install --agent cursor --component sdd,skills,persona

# Preview without making any changes
eva install --agent claude-code --dry-run

# See what's supported
eva install --help
```

### Presets

| Preset | Components included |
|---|---|
| `full` | sdd, skills, mcp, persona |
| `minimal` | persona |
| `custom` | you pick with `--component` |

---

## How it works

```
eva install --agent claude-code --preset full
      │
      ├── 1. Detects your OS and platform
      ├── 2. Validates Claude Code is installed
      ├── 3. Backs up your existing config
      ├── 4. Injects SDD workflow → CLAUDE.md
      ├── 5. Injects your Skills → CLAUDE.md
      ├── 6. Registers MCP servers → settings.json
      ├── 7. Injects Persona → CLAUDE.md
      └── 8. Verifies everything looks good
```

If anything fails, it rolls back automatically. Your original config is never lost.

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
├── cmd/directive/          # CLI entrypoint
├── internal/
│   ├── agents/             # Agent interface + shared types
│   │   ├── claudecode/     # Claude Code adapter
│   │   ├── cursor/         # Cursor adapter
│   │   ├── copilot/        # GitHub Copilot adapter
│   │   ├── opencode/       # OpenCode adapter
│   │   └── windsurf/       # Windsurf adapter
│   ├── app/                # Cobra command wiring
│   ├── cli/                # Flags, validation, InstallOptions
│   ├── system/             # OS and distro detection
│   ├── components/         # SDD, Skills, MCP, Persona logic  [Phase 3]
│   ├── planner/            # Dependency graph                  [Phase 4]
│   ├── pipeline/           # Staged execution + rollback       [Phase 4]
│   ├── backup/             # Config snapshot + restore         [Phase 4]
│   ├── verify/             # Post-install health checks        [Phase 4]
│   └── tui/                # Interactive terminal UI           [Phase 5]
└── Makefile
```

---

## Roadmap

- [x] Phase 1 — CLI foundation, OS detection, flag validation
- [x] Phase 2 — Agent adapters (Claude Code, Cursor, Copilot, OpenCode, Windsurf)
- [ ] Phase 3 — Components (SDD, Skills, MCP, Persona)
- [ ] Phase 4 — Pipeline, backup, and rollback
- [ ] Phase 5 — Interactive TUI
- [ ] Phase 6 — Distribution (Homebrew, GoReleaser)

---

## License

MIT — built with 🤖 by [@codebymaribel](https://github.com/codebymaribel)