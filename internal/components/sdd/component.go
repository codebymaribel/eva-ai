// Package sdd implements the Spec-Driven Development component.
//
// It embeds all SDD phase templates into the binary at compile time
// using Go's embed package — no external files needed at runtime.
//
// When injected into an agent, it writes:
//   - All phase agents to the agent's config (CLAUDE.md, .cursorrules, etc.)
//   - The workflow overview as a preamble
//   - The .eva/ README template to the project directory
package sdd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// — Embed all phase templates at compile time —
// The //go:embed directive tells the Go compiler to read the file
// and store its contents in the variable. No file I/O at runtime.

//go:embed templates/workflow.md
var workflowTemplate string

//go:embed templates/phases/sdd-scan-agent.md
var scanTemplate string

//go:embed templates/phases/sdd-directive-agent.md
var directiveTemplate string

//go:embed templates/phases/sdd-schematics-agent.md
var schematicsTemplate string

//go:embed templates/phases/sdd-sequence-agent.md
var sequenceTemplate string

//go:embed templates/phases/sdd-execute-agent.md
var executeTemplate string

//go:embed templates/phases/sdd-debrief-agent.md
var debriefTemplate string

//go:embed templates/eva-README-template.md
var evaREADMETemplate string

//go:embed templates/eva-orchestrator.md
var orchestratorTemplate string

// Phase represents a single SDD phase with its name and template content.
type Phase struct {
	Name     string
	Command  string
	Template string
}

// Phases returns all SDD phases in order.
// Each phase has a name, the slash-command the dev uses,
// and the full template content to inject.
func Phases() []Phase {
	return []Phase{
		{Name: "scan", Command: "/scan", Template: scanTemplate},
		{Name: "directive", Command: "/directive", Template: directiveTemplate},
		{Name: "schematics", Command: "/schematics", Template: schematicsTemplate},
		{Name: "sequence", Command: "/sequence", Template: sequenceTemplate},
		{Name: "execute", Command: "/execute", Template: executeTemplate},
		{Name: "debrief", Command: "/debrief", Template: debriefTemplate},
	}
}

// Component implements the SDD injectable component.
// It knows how to render and inject SDD content into any agent config.
type Component struct{}

// New constructs a new SDD Component.
func New() *Component {
	return &Component{}
}

// Name returns the component identifier.
func (c *Component) Name() string {
	return "sdd"
}

// RenderOrchestrator returns the EVA orchestrator template.
// This is injected separately from the phase agents — it goes first
// in the config file so the agent sees EVA before the subagents.
//
// EVA is intentionally restricted to plan mode with read-only tools
// to prevent it from consuming context by writing code directly.
func (c *Component) RenderOrchestrator() string {
	return orchestratorTemplate
}

// Render builds the full SDD block to inject into an agent config file.
// Order matters — EVA orchestrator goes first, then the phase agents.
// The agent reads top to bottom, so EVA needs to be the first thing it sees.
func (c *Component) Render() string {
	var b strings.Builder

	// 1. EVA orchestrator — always first
	b.WriteString("# EVA — Spec-Driven Development\n\n")
	b.WriteString(orchestratorTemplate)
	b.WriteString("\n\n---\n\n")

	// 2. Workflow overview
	b.WriteString("## Workflow Overview\n\n")
	b.WriteString(workflowTemplate)
	b.WriteString("\n\n---\n\n")

	// 3. Phase subagents — in order
	b.WriteString("## Phase Subagents\n\n")
	for _, phase := range Phases() {
		b.WriteString(fmt.Sprintf("### %s\n\n", phase.Command))
		b.WriteString(phase.Template)
		b.WriteString("\n\n---\n\n")
	}

	return b.String()
}

// InjectIntoFile appends the SDD block to an existing config file.
// If the file already contains SDD content (detected by the header),
// it replaces the existing block instead of duplicating it.
//
// This is the core operation — called by each agent adapter
// when the sdd component is selected for installation.
func (c *Component) InjectIntoFile(path string) error {
	// Read existing content — empty string if file doesn't exist yet
	existing, err := readFileOrEmpty(path)
	if err != nil {
		return fmt.Errorf("sdd: failed to read %s: %w", path, err)
	}

	rendered := c.Render()

	// If SDD block already exists, replace it
	// If not, append it to whatever is already in the file
	var updated string
	if containsSDD(existing) {
		updated = replaceSDD(existing, rendered)
	} else {
		updated = appendSDD(existing, rendered)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("sdd: failed to create directory for %s: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("sdd: failed to write %s: %w", path, err)
	}

	return nil
}

// InitEvaDir creates the .eva/ directory structure in the given project dir.
// This is called by `eva init` — it sets up the knowledge base
// that the SDD agents will read and write during missions.
//
// Structure created:
//
//	projectDir/.eva/
//	├── README.md         ← universal entry point for agents
//	├── memory.md         ← empty, agents will populate this
//	└── phases/           ← empty dir, agents write here during missions
func (c *Component) InitEvaDir(projectDir string) error {
	evaDir := filepath.Join(projectDir, ".eva")
	phasesDir := filepath.Join(evaDir, "phases")
	skillsDir := filepath.Join(evaDir, "skills")
	sharedSkillsDir := filepath.Join(evaDir, "shared-skills")
	docsDir := filepath.Join(evaDir, "docs")

	// Create all directories
	for _, dir := range []string{phasesDir, skillsDir, sharedSkillsDir, docsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("sdd: failed to create directory %s: %w", dir, err)
		}
	}

	// Write .eva/README.md — the universal entry point
	readmePath := filepath.Join(evaDir, "README.md")
	if err := writeIfNotExists(readmePath, evaREADMETemplate); err != nil {
		return fmt.Errorf("sdd: failed to write .eva/README.md: %w", err)
	}

	// Write .eva/memory.md — empty, agents populate this
	memoryPath := filepath.Join(evaDir, "memory.md")
	memoryContent := "# Mission Memory\n\nAccumulated context from SDD missions.\n"
	if err := writeIfNotExists(memoryPath, memoryContent); err != nil {
		return fmt.Errorf("sdd: failed to write .eva/memory.md: %w", err)
	}

	// Write .eva/.gitignore — phases and memory are local, skills are shared
	gitignorePath := filepath.Join(evaDir, ".gitignore")
	gitignoreContent := "# SDD phase outputs — local context, not shared\nphases/\nmemory.md\n"
	if err := writeIfNotExists(gitignorePath, gitignoreContent); err != nil {
		return fmt.Errorf("sdd: failed to write .eva/.gitignore: %w", err)
	}

	return nil
}

// — internal helpers —

// sddMarkerStart and sddMarkerEnd are the sentinel strings used to
// detect and replace existing SDD blocks in agent config files.
const sddMarkerStart = "<!-- eva:sdd:start -->"
const sddMarkerEnd = "<!-- eva:sdd:end -->"

// containsSDD returns true if the content already has an SDD block.
func containsSDD(content string) bool {
	return strings.Contains(content, sddMarkerStart)
}

// appendSDD adds the SDD block at the end of existing content,
// wrapped in markers so future runs can find and replace it.
func appendSDD(existing, rendered string) string {
	separator := "\n\n"
	if existing == "" {
		separator = ""
	}
	return existing + separator +
		sddMarkerStart + "\n" +
		rendered +
		sddMarkerEnd + "\n"
}

// replaceSDD finds the existing SDD block between markers and
// replaces it with the new rendered content.
func replaceSDD(existing, rendered string) string {
	start := strings.Index(existing, sddMarkerStart)
	end := strings.Index(existing, sddMarkerEnd)

	if start == -1 || end == -1 {
		// Markers not found — fall back to append
		return appendSDD(existing, rendered)
	}

	before := existing[:start]
	after := existing[end+len(sddMarkerEnd):]

	return before +
		sddMarkerStart + "\n" +
		rendered +
		sddMarkerEnd + "\n" +
		after
}

// readFileOrEmpty reads a file and returns its content.
// If the file doesn't exist, returns empty string without error —
// this is the first-install case.
func readFileOrEmpty(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}

// writeIfNotExists writes content to path only if the file doesn't exist.
// Used for init files that the dev will customize — we never overwrite.
func writeIfNotExists(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // file exists, skip
	}
	return os.WriteFile(path, []byte(content), 0644)
}