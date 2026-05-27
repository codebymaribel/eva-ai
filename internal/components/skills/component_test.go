package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCoreSkills(t *testing.T) {
	skills := CoreSkills()
	if len(skills) != 3 {
		t.Fatalf("CoreSkills() returned %d skills, want 3", len(skills))
	}

	names := make(map[string]bool)
	for _, s := range skills {
		names[s.Name] = true
		if s.Content == "" {
			t.Errorf("skill %q has empty Content", s.Name)
		}
		if s.Trigger == "" {
			t.Errorf("skill %q has empty Trigger", s.Name)
		}
	}

	for _, want := range []string{"architecture", "testing", "git-workflow"} {
		if !names[want] {
			t.Errorf("CoreSkills() missing skill %q", want)
		}
	}
}

func TestComponent_Name(t *testing.T) {
	c := New()
	if c.Name() != "skills" {
		t.Errorf("Name() = %q, want %q", c.Name(), "skills")
	}
}

func TestComponent_Template(t *testing.T) {
	c := New()
	tmpl := c.Template()
	if tmpl == "" {
		t.Error("Template() returned empty string")
	}
	if !contains(tmpl, "name:") {
		t.Error("Template() missing frontmatter 'name:' field")
	}
}

func TestComponent_Render(t *testing.T) {
	c := New()
	rendered := c.Render()

	for _, want := range []string{"architecture", "testing", "git-workflow"} {
		if !contains(rendered, want) {
			t.Errorf("Render() missing skill %q", want)
		}
	}
	if !contains(rendered, "## Skills") {
		t.Error("Render() missing '## Skills' header")
	}
}

func TestComponent_InitSkillsDir(t *testing.T) {
	dir := t.TempDir()
	c := New()

	if err := c.InitSkillsDir(dir); err != nil {
		t.Fatalf("InitSkillsDir() error: %v", err)
	}

	// Check core skills were created
	for _, name := range []string{"architecture", "testing", "git-workflow"} {
		path := filepath.Join(dir, ".eva", "skills", name, "SKILL.md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", path)
		}
	}

	// Check README was created
	readmePath := filepath.Join(dir, ".eva", "skills", "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Error("expected .eva/skills/README.md to exist")
	}
}

func TestComponent_AddFromPath(t *testing.T) {
	dir := t.TempDir()
	c := New()

	// Create a source skill file
	srcDir := t.TempDir()
	skillContent := "---\nname: my-skill\ndescription: test skill\n---\n\n# My Skill\n"
	writeFile(t, srcDir, "SKILL.md", skillContent)

	// Initialize .eva/skills/ first
	if err := c.InitSkillsDir(dir); err != nil {
		t.Fatalf("InitSkillsDir() error: %v", err)
	}

	name, err := c.AddFromPath(dir, filepath.Join(srcDir, "SKILL.md"))
	if err != nil {
		t.Fatalf("AddFromPath() error: %v", err)
	}

	if name != "my-skill" {
		t.Errorf("AddFromPath() name = %q, want %q", name, "my-skill")
	}

	// Verify file was copied
	destPath := filepath.Join(dir, ".eva", "skills", "my-skill", "SKILL.md")
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("expected skill file to be copied to .eva/skills/my-skill/SKILL.md")
	}
}

func TestNormalizeGitHubURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "github blob URL",
			input: "https://github.com/JuliusBrussee/caveman/blob/main/skills/caveman/SKILL.md",
			want:  "https://raw.githubusercontent.com/JuliusBrussee/caveman/main/skills/caveman/SKILL.md",
		},
		{
			name:  "already raw URL",
			input: "https://raw.githubusercontent.com/user/repo/main/SKILL.md",
			want:  "https://raw.githubusercontent.com/user/repo/main/SKILL.md",
		},
		{
			name:  "other URL",
			input: "https://example.com/SKILL.md",
			want:  "https://example.com/SKILL.md",
		},
		{
			name:  "http github URL",
			input: "http://github.com/user/repo/blob/dev/path/SKILL.md",
			want:  "http://raw.githubusercontent.com/user/repo/dev/path/SKILL.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeGitHubURL(tt.input)
			if got != tt.want {
				t.Errorf("normalizeGitHubURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractSkillName(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "with frontmatter",
			content: "---\nname: caveman\ndescription: test\n---\n",
			want:    "caveman",
		},
		{
			name:    "no name field",
			content: "---\ndescription: test\n---\n",
			want:    "",
		},
		{
			name:    "empty content",
			content: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSkillName(tt.content)
			if got != tt.want {
				t.Errorf("extractSkillName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestComponent_RefreshREADME(t *testing.T) {
	dir := t.TempDir()
	c := New()

	// Init first
	if err := c.InitSkillsDir(dir); err != nil {
		t.Fatalf("InitSkillsDir() error: %v", err)
	}

	// Add a custom skill
	writeFile(t, dir, ".eva/skills/custom/SKILL.md", "---\nname: custom\ndescription: my custom skill\ntrigger: custom, test\n---\n\n# Custom\n")

	// Refresh
	if err := c.RefreshREADME(dir); err != nil {
		t.Fatalf("RefreshREADME() error: %v", err)
	}

	// Verify README contains the custom skill
	readmePath := filepath.Join(dir, ".eva", "skills", "README.md")
	content := readFile(t, readmePath)

	if !contains(content, "custom") {
		t.Error("README missing 'custom' skill after refresh")
	}
	if !contains(content, "architecture") {
		t.Error("README missing 'architecture' skill after refresh")
	}
}

func TestComponent_InjectIntoAgentConfigs(t *testing.T) {
	dir := t.TempDir()
	c := New()

	// Init and add a skill
	if err := c.InitSkillsDir(dir); err != nil {
		t.Fatalf("InitSkillsDir() error: %v", err)
	}

	skillContent := "---\nname: test-skill\ndescription: test\n---\n\n# Test Skill\nDo the thing.\n"
	writeFile(t, dir, ".eva/skills/test-skill/SKILL.md", skillContent)

	// Create a fake agent config
	configPath := filepath.Join(dir, ".claude", "CLAUDE.md")
	writeFile(t, dir, ".claude/CLAUDE.md", "# Claude Config\n\nSome existing content.\n")

	// Inject
	injected := c.InjectIntoAgentConfigs(dir, "test-skill")

	if len(injected) != 1 {
		t.Fatalf("InjectIntoAgentConfigs() injected into %d configs, want 1", len(injected))
	}

	// Verify the config now contains the skill
	configContent := readFile(t, configPath)
	if !contains(configContent, "## Skill: test-skill (auto-loaded)") {
		t.Error("agent config missing skill injection marker")
	}
	if !contains(configContent, "Do the thing.") {
		t.Error("agent config missing skill content")
	}
}

func TestComponent_InjectIntoAgentConfigs_NoDuplicate(t *testing.T) {
	dir := t.TempDir()
	c := New()

	if err := c.InitSkillsDir(dir); err != nil {
		t.Fatalf("InitSkillsDir() error: %v", err)
	}

	writeFile(t, dir, ".eva/skills/test-skill/SKILL.md", "---\nname: test-skill\n---\n\n# Test\n")
	writeFile(t, dir, ".claude/CLAUDE.md", "# Config\n")

	// Inject twice
	c.InjectIntoAgentConfigs(dir, "test-skill")
	c.InjectIntoAgentConfigs(dir, "test-skill")

	// Count occurrences of the marker
	configContent := readFile(t, filepath.Join(dir, ".claude", "CLAUDE.md"))
	count := countOccurrences(configContent, "## Skill: test-skill (auto-loaded)")
	if count != 1 {
		t.Errorf("skill marker appears %d times, want 1 (duplicate injection)", count)
	}
}

// helpers

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", path, err)
	}
	return string(data)
}
