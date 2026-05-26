# 🤖 eva-ai

> Personaliza tu flujo de desarrollo agéntico desde una unidad de comando central.

## 🚧 PROYECTO EN CONSTRUCCION. NO HAGAS PULL 🚧

eva-ai es una herramienta CLI que inyecta tus Skills personalizadas, flujo de trabajo SDD, servidores MCP y Persona en todos tus agentes de AI.
Vos decile a EVA qué necesitás y ella se encarga de que tus agentes trabajen como tú necesitas.

Inspirado en [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai), construido desde cero con foco en la facilidad de configuración al desarrollador.

```bash
eva install --agent claude-code,cursor --preset full
```

---

## ¿Por qué eva-ai?

Cada agente de AI tiene su propio formato de configuración, su propia carpeta, sus propias reglas. Configurarlos de forma consistente en tu máquina es tedioso y propenso a errores.

Además, requiere de una cantidad de conocimiento y contexto que no todo developer tiene, sobretodo a la hora de trabajar en proyectos de equipo.

eva-ai resuelve eso. Definís tu flujo de trabajo una vez, lo inyectás en todos lados.

---

## Agentes soportados

| Agente | Ubicación de configuración |
|---|---|
| **Claude Code** | `~/.claude/CLAUDE.md` |
| **Cursor** | `~/.cursor/rules/.cursorrules` |
| **GitHub Copilot** | `.github/copilot-instructions.md` |
| **OpenCode** | `~/.config/opencode/AGENTS.md` |
| **Windsurf** | `~/.codeium/windsurf/memories/global_rules.md` |

---

## Componentes

| Componente | Descripción |
|---|---|
| `sdd` | Spec-Driven Development — planificá antes de codear |
| `skills` | Tus patrones de código y buenas prácticas |
| `mcp` | Registro de servidores MCP (Context7, etc.) |
| `persona` | Comportamiento del agente, estilo y tono |

---

## Instalación

```bash
# Homebrew (macOS) — próximamente
brew tap codebymaribel/tap
brew install eva-ai

# Go install
go install github.com/codebymaribel/eva-ai/cmd/directive@latest
```

---

## Uso

```bash
# Stack completo — instala todo en Claude Code y Cursor
eva install --agent claude-code,cursor --preset full

# Minimal — solo persona
eva install --agent claude-code --preset minimal

# Elegís los componentes específicos
eva install --agent cursor --component sdd,skills,persona

# Vista previa sin hacer ningún cambio
eva install --agent claude-code --dry-run

# Ver qué está soportado
eva install --help
```

### Presets

| Preset | Componentes incluidos |
|---|---|
| `full` | sdd, skills, mcp, persona |
| `minimal` | persona |
| `custom` | vos elegís con `--component` |

---

## Cómo funciona

```
eva install --agent claude-code --preset full
      │
      ├── 1. Detecta tu OS y plataforma
      ├── 2. Valida que Claude Code esté instalado
      ├── 3. Hace backup de tu configuración existente
      ├── 4. Inyecta el flujo SDD → CLAUDE.md
      ├── 5. Inyecta tus Skills → CLAUDE.md
      ├── 6. Registra los servidores MCP → settings.json
      ├── 7. Inyecta la Persona → CLAUDE.md
      └── 8. Verifica que todo quedó bien
```

Si algo falla, hace rollback automáticamente. Tu configuración original nunca se pierde.

---

## Soporte de plataformas

| OS | Estado |
|---|---|
| macOS | ✅ Soportado |
| Linux | ✅ Soportado |
| Windows | 🚧 Próximamente |

---

## Desarrollo

```bash
# Clonar
git clone https://github.com/codebymaribel/eva-ai.git
cd eva-ai

# Instalar dependencias
go mod tidy

# Correr tests
go test ./... -v

# Compilar binario
go build -o bin/eva ./cmd/directive

# Ejemplo dry-run
./bin/eva install --agent claude-code --dry-run
```

### Estructura del proyecto

```
eva-ai/
├── cmd/directive/          # Punto de entrada del CLI
├── internal/
│   ├── agents/             # Interface Agent + tipos compartidos
│   │   ├── claudecode/     # Adapter de Claude Code
│   │   ├── cursor/         # Adapter de Cursor
│   │   ├── copilot/        # Adapter de GitHub Copilot
│   │   ├── opencode/       # Adapter de OpenCode
│   │   └── windsurf/       # Adapter de Windsurf
│   ├── app/                # Wiring de comandos Cobra
│   ├── cli/                # Flags, validación, InstallOptions
│   ├── system/             # Detección de OS y distro
│   ├── components/         # Lógica de SDD, Skills, MCP, Persona  [Fase 3]
│   ├── planner/            # Grafo de dependencias                 [Fase 4]
│   ├── pipeline/           # Ejecución staged + rollback           [Fase 4]
│   ├── backup/             # Snapshot y restore de configs         [Fase 4]
│   ├── verify/             # Health checks post-instalación        [Fase 4]
│   └── tui/                # Terminal UI interactiva               [Fase 5]
└── Makefile
```

---

## Roadmap

- [x] Fase 1 — Base del CLI, detección de OS, validación de flags
- [x] Fase 2 — Adapters de agentes (Claude Code, Cursor, Copilot, OpenCode, Windsurf)
- [ ] Fase 3 — Componentes (SDD, Skills, MCP, Persona)
- [ ] Fase 4 — Pipeline, backup y rollback
- [ ] Fase 5 — TUI interactiva
- [ ] Fase 6 — Distribución (Homebrew, GoReleaser)

---

## Licencia

MIT — construido con 🤖 por [@codebymaribel](https://github.com/codebymaribel)