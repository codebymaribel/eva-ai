
package cli

import (
	"fmt"
	"strings"
)

type Agent string

const (
	AgentClaudeCode  Agent = "claude-code"
	AgentOpenCode    Agent = "opencode"
	AgentCursor      Agent = "cursor"
	AgentCopilot     Agent = "copilot"
	AgentWindsurf    Agent = "windsurf"
	AgentAntigravity Agent = "antigravity"
)

type Component string

const (
	ComponentSDD     Component = "sdd"
	ComponentSkills  Component = "skills"
	ComponentMCP     Component = "mcp"
	ComponentPersona Component = "persona"
)

type Preset string

const (
	PresetFull    Preset = "full"
	PresetMinimal Preset = "minimal"
	PresetCustom  Preset = "custom"
)


type InstallOptions struct {
	Agents     []Agent
	Components []Component
	Preset     Preset
	DryRun     bool
}


var SupportedAgents = map[Agent]bool{
	AgentClaudeCode: true,
	AgentOpenCode:   true,
	AgentCursor:     true,
	AgentCopilot:    true,
	AgentWindsurf:   true,
	AgentAntigravity: true,
}


var SupportedComponents = map[Component]bool{
	ComponentSDD:     true,
	ComponentSkills:  true,
	ComponentMCP:     true,
	ComponentPersona: true,
}


var PresetComponents = map[Preset][]Component{
	PresetFull:    {ComponentSDD, ComponentSkills, ComponentMCP, ComponentPersona},
	PresetMinimal: {ComponentPersona},
	PresetCustom:  {},
}

// ParseAgents converts a comma-separated string (e.g. "claude-code,cursor")
// to a validated slice of Agent values.
func ParseAgents(raw string) ([]Agent, error) {
	if raw == "" {
		return nil, fmt.Errorf("--agent is required: choose from %s", agentList())
	}

	parts := strings.Split(raw, ",")
	agents := make([]Agent, 0, len(parts))

	for _, p := range parts {
		a := Agent(strings.TrimSpace(p))
		if !SupportedAgents[a] {
			return nil, fmt.Errorf("unknown agent %q — supported: %s", a, agentList())
		}
		agents = append(agents, a)
	}

	return agents, nil
}

// ParseComponents converts a comma-separated string to a validated slice
// of Component values. Returns nil (not an error) when the string is empty,
// since components can be derived from a preset.
func ParseComponents(raw string) ([]Component, error) {
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	components := make([]Component, 0, len(parts))

	for _, p := range parts {
		c := Component(strings.TrimSpace(p))
		if !SupportedComponents[c] {
			return nil, fmt.Errorf("unknown component %q — supported: %s", c, componentList())
		}
		components = append(components, c)
	}

	return components, nil
}

// ParsePreset validates and returns a Preset value.
// Defaults to PresetFull when the raw string is empty.
func ParsePreset(raw string) (Preset, error) {
	if raw == "" {
		return PresetFull, nil
	}

	p := Preset(strings.ToLower(strings.TrimSpace(raw)))
	switch p {
	case PresetFull, PresetMinimal, PresetCustom:
		return p, nil
	default:
		return "", fmt.Errorf("unknown preset %q — choose: full, minimal, custom", raw)
	}
}

// Validate checks that the InstallOptions are consistent and complete.
// It resolves components from the preset when no explicit components are given.
func (o *InstallOptions) Validate() error {
	if len(o.Agents) == 0 {
		return fmt.Errorf("at least one --agent is required")
	}

	// If no explicit components, derive from preset
	if len(o.Components) == 0 {
		o.Components = PresetComponents[o.Preset]
	}

	if len(o.Components) == 0 {
		return fmt.Errorf("no components to install — specify --component or use a preset other than custom")
	}

	return nil
}

// agentList returns a human-readable list of supported agents.
func agentList() string {
	names := make([]string, 0, len(SupportedAgents))
	for a := range SupportedAgents {
		names = append(names, string(a))
	}
	return strings.Join(names, ", ")
}

// componentList returns a human-readable list of supported components.
func componentList() string {
	names := make([]string, 0, len(SupportedComponents))
	for c := range SupportedComponents {
		names = append(names, string(c))
	}
	return strings.Join(names, ", ")
}