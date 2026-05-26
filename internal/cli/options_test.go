package cli_test

import (
	"testing"

	"github.com/codebymaribel/eva-ai/internal/cli"
)

// TestParseAgents_Valid checks that a valid comma-separated agent list is parsed.
func TestParseAgents_Valid(t *testing.T) {
	got, err := cli.ParseAgents("claude-code,cursor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 agents, got %d", len(got))
	}
	if got[0] != cli.AgentClaudeCode {
		t.Errorf("got[0] = %q, want %q", got[0], cli.AgentClaudeCode)
	}
	if got[1] != cli.AgentCursor {
		t.Errorf("got[1] = %q, want %q", got[1], cli.AgentCursor)
	}
}

// TestParseAgents_Unknown ensures an unknown agent returns an error.
func TestParseAgents_Unknown(t *testing.T) {
	_, err := cli.ParseAgents("claude-code,vim-ai")
	if err == nil {
		t.Error("expected error for unknown agent, got nil")
	}
}

// TestParseAgents_Empty ensures empty input returns an error.
func TestParseAgents_Empty(t *testing.T) {
	_, err := cli.ParseAgents("")
	if err == nil {
		t.Error("expected error for empty agents, got nil")
	}
}

// TestParseComponents_Valid checks that valid components are parsed correctly.
func TestParseComponents_Valid(t *testing.T) {
	got, err := cli.ParseComponents("sdd,skills,mcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 components, got %d", len(got))
	}
}

// TestParseComponents_Empty returns nil without error (preset-derived).
func TestParseComponents_Empty(t *testing.T) {
	got, err := cli.ParseComponents("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for empty components, got %v", got)
	}
}

// TestParseComponents_Unknown ensures unknown component returns error.
func TestParseComponents_Unknown(t *testing.T) {
	_, err := cli.ParseComponents("sdd,nonexistent")
	if err == nil {
		t.Error("expected error for unknown component, got nil")
	}
}

// TestParsePreset_Valid covers all valid preset values.
func TestParsePreset_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  cli.Preset
	}{
		{"full", cli.PresetFull},
		{"minimal", cli.PresetMinimal},
		{"custom", cli.PresetCustom},
		{"", cli.PresetFull}, // empty defaults to full
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := cli.ParsePreset(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParsePreset(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestParsePreset_Unknown ensures an unknown preset returns an error.
func TestParsePreset_Unknown(t *testing.T) {
	_, err := cli.ParsePreset("ultra")
	if err == nil {
		t.Error("expected error for unknown preset, got nil")
	}
}

// TestInstallOptions_Validate_ResolvesFromPreset checks that when no explicit
// components are provided, they are resolved from the preset.
func TestInstallOptions_Validate_ResolvesFromPreset(t *testing.T) {
	opts := &cli.InstallOptions{
		Agents:     []cli.Agent{cli.AgentClaudeCode},
		Components: nil,
		Preset:     cli.PresetFull,
	}

	if err := opts.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}

	if len(opts.Components) == 0 {
		t.Error("expected components to be resolved from preset, got none")
	}
}

// TestInstallOptions_Validate_NoAgents returns error when agents are empty.
func TestInstallOptions_Validate_NoAgents(t *testing.T) {
	opts := &cli.InstallOptions{
		Agents: []cli.Agent{},
		Preset: cli.PresetFull,
	}

	if err := opts.Validate(); err == nil {
		t.Error("expected error for no agents, got nil")
	}
}

// TestInstallOptions_Validate_CustomNoComponents returns error when
// preset is custom and no components are specified.
func TestInstallOptions_Validate_CustomNoComponents(t *testing.T) {
	opts := &cli.InstallOptions{
		Agents:     []cli.Agent{cli.AgentClaudeCode},
		Components: nil,
		Preset:     cli.PresetCustom,
	}

	if err := opts.Validate(); err == nil {
		t.Error("expected error for custom preset with no components, got nil")
	}
}