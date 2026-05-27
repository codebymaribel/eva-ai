# 🤖 eva-ai

## 🚧 PROYECTO EN CONTRUCCIÓN. 🚧

> Personaliza tu flujo de desarrollo agéntico desde una unidad de comando central.

eva-ai es una herramienta CLI que inyecta tus Skills personalizadas, flujo de trabajo SDD, servidores MCP y Persona en todos tus agentes de AI.
Vos decile a EVA qué necesitás y ella se encarga de que tus agentes trabajen como tú necesitas.

Inspirado en [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai), construido desde cero con foco en la facilidad de configuración al desarrollador.

```bash
# Inicializá tu proyecto
eva init

# Inyectá en tu agente de AI
eva install --agent claude-code --preset full

# Agregá skills externos
eva skill add https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md
```

---

## ¿Por qué eva-ai?

Cada agente de AI tiene su propio formato de configuración, su propia carpeta, sus propias reglas. Configurarlos de forma consistente en tu máquina es tedioso y propenso a errores.

Además, requiere de una cantidad de conocimiento y contexto que no todo developer tiene, sobretodo a la hora de trabajar en proyectos de equipo.

eva-ai resuelve eso. Definís tu flujo de trabajo una vez, lo inyectás en todos lados.

---

## Agentes soportados

| Agente | Ubicación de configuración | Estado |
|---|---|---|
| **Claude Code** | `~/.claude/CLAUDE.md` | ✅ Soportado |
| **Cursor** | `~/.cursor/rules/.cursorrules` | 🚧 Adapter pendiente |
| **GitHub Copilot** | `.github/copilot-instructions.md` | 🚧 Adapter pendiente |
| **OpenCode** | `~/.config/opencode/AGENTS.md` | 🚧 Adapter pendiente |
| **Windsurf** | `~/.codeium/windsurf/memories/global_rules.md` | 🚧 Adapter pendiente |

---

## Componentes

| Componente | Descripción | Estado |
|---|---|---|
| `sdd` | Spec-Driven Development — planificá antes de codear | ✅ |
| `skills` | Tus patrones de código y buenas prácticas | ✅ |
| `mcp` | Registro de servidores MCP (Context7, etc.) | 🚧 |
| `persona` | Comportamiento del agente, estilo y tono | 🚧 |

---

## Instalación

```bash
# Go install
go install github.com/codebymaribel/eva-ai/cmd/directive@latest

# O compilar desde source
git clone https://github.com/codebymaribel/eva-ai.git
cd eva-ai
go build -o bin/eva ./cmd/directive
```

---

## Uso

### `eva init` — Configurá tu proyecto

Escaneá tu codebase, detecta el stack, y crea `.eva/` con skills core y skills específicos del proyecto.

```bash
# Init completo — escanea proyecto y crea todos los skills
eva init

# Init sin skills — solo estructura SDD
eva init --no-skills
```

**Lo que crea:**

```
.eva/
├── README.md                     ← punto de entrada del agente
├── memory.md                     ← contexto de misiones (gitignored)
├── phases/                       ← outputs de fases SDD (gitignored)
├── skills/
│   ├── README.md                 ← índice de skills con triggers
│   ├── architecture/SKILL.md     ← core: stack, patrones, estructura
│   ├── testing/SKILL.md          ← core: convenciones de testing
│   ├── git-workflow/SKILL.md     ← core: convenciones de git
│   └── domain/SKILL.md           ← generado: lógica de negocio
└── docs/                         ← documentación de features
```

### `eva install` — Inyectá en agentes

```bash
# Stack completo para Claude Code
eva install --agent claude-code --preset full

# Elegís componentes específicos
eva install --agent claude-code --component sdd,skills

# Vista previa sin hacer cambios
eva install --agent claude-code --dry-run
```

**Presets:**

| Preset | Componentes incluidos |
|---|---|
| `full` | sdd, skills, mcp, persona |
| `minimal` | persona |
| `custom` | vos elegís con `--component` |

### `eva skill add` — Agregá skills externos

```bash
# Desde un archivo local
eva skill add ./path/to/SKILL.md

# Desde una URL de GitHub
eva skill add https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md

# Desde una URL raw
eva skill add https://raw.githubusercontent.com/user/repo/main/SKILL.md
```

**Qué pasa:**
1. Descarga/copia el SKILL.md a `.eva/skills/<nombre>/`
2. Actualiza `.eva/skills/README.md` con la nueva entrada
3. Inyecta el skill en configs de agentes existentes (CLAUDE.md, .cursorrules, etc.)

El agente ve las reglas del skill desde el primer mensaje — sin necesidad de invocación manual.

---

## Cómo funciona

```
eva init
  ├── 1. Escanea proyecto (lenguaje, framework, patrones, tooling)
  ├── 2. Crea estructura de directorios .eva/
  ├── 3. Escribe skills core (architecture, testing, git-workflow)
  └── 4. Genera skills específicos del proyecto

eva install --agent claude-code --component sdd,skills
  ├── 1. Detecta tu OS y plataforma
  ├── 2. Valida que Claude Code esté instalado
  ├── 3. Lee CLAUDE.md existente
  ├── 4. Renderiza componentes seleccionados
  ├── 5. Appendea a CLAUDE.md
  └── 6. Listo — el agente ve la nueva config en la próxima sesión

eva skill add <url>
  ├── 1. Descarga SKILL.md (normaliza URLs de GitHub)
  ├── 2. Guarda en .eva/skills/<nombre>/SKILL.md
  ├── 3. Refresca .eva/skills/README.md
  └── 4. Inyecta en configs de agentes existentes
```

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
├── cmd/directive/              # Punto de entrada del CLI
├── internal/
│   ├── agents/                 # Interface Agent + adapters
│   │   ├── agent.go            # Interfaces Agent + InjectableAgent
│   │   └── claudecode/         # Adapter de Claude Code
│   ├── app/                    # Wiring de comandos Cobra
│   ├── cli/                    # Flags, validación, InstallOptions
│   ├── components/
│   │   ├── sdd/                # Componente SDD (orchestrator + 6 fases)
│   │   └── skills/             # Componente Skills (core + template)
│   ├── scanner/                # Detección de stack y patrones del proyecto
│   └── system/                 # Detección de OS y distro
└── Makefile
```

---

## Roadmap

- [x] Fase 1 — Base del CLI, detección de OS, validación de flags
- [x] Fase 2 — Adapters de agentes (Claude Code)
- [x] Fase 3 — Componentes (SDD, Skills)
- [ ] Fase 3 — Componentes (MCP, Persona)
- [ ] Fase 2 — Adapters de agentes (Cursor, Copilot, OpenCode, Windsurf)
- [ ] Fase 4 — Pipeline, backup y rollback
- [ ] Fase 5 — TUI interactiva
- [ ] Fase 6 — Distribución (Homebrew, GoReleaser)

---

## Licencia

MIT — construido con 🤖 por [@codebymaribel](https://github.com/codebymaribel)
