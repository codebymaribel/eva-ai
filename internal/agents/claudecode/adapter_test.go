package claudecode_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codebymaribel/eva-ai/internal/agents/claudecode"
)


// newTestAdapter creates an adapter pointing a temporal directory.
func newTestAdapter(t *testing.T) (*claudecode.Adapter, string) {
	t.Helper()
	dir := t.TempDir() // Go creates & cleans the directory when test is completed.
	return claudecode.New(dir), dir
}

// TestAdapter_Name verifies adapter returns the right name.
func TestAdapter_Name(t *testing.T) {
	a, _ := newTestAdapter(t)
	if got := a.Name(); got != "claude-code" {
		t.Errorf("Name() = %q, want %q", got, "claude-code")
	}
}

// TestAdapter_ConfigDir verifies route is equal to homeDir + ".claude".
func TestAdapter_ConfigDir(t *testing.T) {
	a, homeDir := newTestAdapter(t)
	want := filepath.Join(homeDir, ".claude")
	if got := a.ConfigDir(); got != want {
		t.Errorf("ConfigDir() = %q, want %q", got, want)
	}
}

// TestAdapter_MainConfigPath verifies complete route to CLAUDE.md.
func TestAdapter_MainConfigPath(t *testing.T) {
	a, homeDir := newTestAdapter(t)
	want := filepath.Join(homeDir, ".claude", "CLAUDE.md")
	if got := a.MainConfigPath(); got != want {
		t.Errorf("MainConfigPath() = %q, want %q", got, want)
	}
}

// TestAdapter_IsInstalled_WithConfigDir verifies IsInstalled returns true
// when config directory exists 
func TestAdapter_IsInstalled_WithConfigDir(t *testing.T) {
	a, homeDir := newTestAdapter(t)

	// Manually create .claude directory
	configDir := filepath.Join(homeDir, ".claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	if !a.IsInstalled() {
		t.Error("IsInstalled() = false, want true when config dir exists")
	}
}

// TestAdapter_IsInstalled_NotInstalled verifies IsInstalled returns false
// when binary or directory doesn't exist
func TestAdapter_IsInstalled_NotInstalled(t *testing.T) {
	a, _ := newTestAdapter(t)
	
	// tempDir doesn't have .claude inside, IsInstalled should be false
	if a.IsInstalled() {
		t.Skip("claude is installed on this machine, skipping not-installed test")
	}
}

// TestAdapter_Validate_NotInstalled verifies CanAdapterRun returns
// ErrNotInstalled when agent is not installed.
func TestAdapter_CanAdapterRun_NotInstalled(t *testing.T) {
	a, _ := newTestAdapter(t)

	if a.IsInstalled() {
		t.Skip("claude is installed on this machine, skipping")
	}

	err := a.CanAdapterRun()
	if err == nil {
		t.Fatal("Validate() should return error when not installed")
	}
}


// TestAdapter_Validate_ConfigDirIsFile verifies CanAdapterRun returns error
// when exists a file (not directory) on config route
func TestAdapter_CanAdapterRun_ConfigDirIsFile(t *testing.T) {
	a, homeDir := newTestAdapter(t)


	configPath := filepath.Join(homeDir, ".claude")
	if err := os.WriteFile(configPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := a.EnsureConfigDir()
	if err == nil {
		t.Error("EnsureConfigDir() should fail when path exists as a file")
	}
}


// TestAdapter_EnsureConfigDir_CreatesDir verifies that EnsureConfigDir
// creates directory if it doesn't exist
func TestAdapter_EnsureConfigDir_CreatesDir(t *testing.T) {
	a, homeDir := newTestAdapter(t)

	if err := a.EnsureConfigDir(); err != nil {
		t.Fatalf("EnsureConfigDir() unexpected error: %v", err)
	}

	expected := filepath.Join(homeDir, ".claude")
	info, err := os.Stat(expected)
	if err != nil {
		t.Fatalf("config dir was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("config dir should be a directory")
	}
}


//  TestAdapter_EnsureConfigDir_Idempotent verifies invoking EnsureConfigDir
//  two times doesn't file 
func TestAdapter_EnsureConfigDir_Idempotent(t *testing.T) {
	a, _ := newTestAdapter(t)

	if err := a.EnsureConfigDir(); err != nil {
		t.Fatalf("first call failed: %v", err)
	}
	if err := a.EnsureConfigDir(); err != nil {
		t.Fatalf("second call failed: %v", err)
	}
}

// TestAdapter_ReadMainMD_FileNotExist verifies an undefined CLAUDE.md
// returns an empty string w/o error (first install)
func TestAdapter_ReadMainMD_FileNotExist(t *testing.T) {
	a, _ := newTestAdapter(t)

	content, err := a.ReadMainMD()
	if err != nil {
		t.Fatalf("ReadMainMD() unexpected error: %v", err)
	}
	if content != "" {
		t.Errorf("ReadMainMD() = %q, want empty string", content)
	}
}

// TestAdapter_WriteAndReadMainConfig verifies the full cicle
// write -> read CLAUDE.md file
func TestAdapter_WriteAndReadMainConfig(t *testing.T) {
	a, _ := newTestAdapter(t)

	original := "# My AI Instructions\n\nAlways write clean code.\n"

	if err := a.WriteNewMD(original); err != nil {
		t.Fatalf("WriteMainConfig() unexpected error: %v", err)
	}

	got, err := a.ReadMainMD()
	if err != nil {
		t.Fatalf("ReadMainMD() unexpected error: %v", err)
	}

	if got != original {
		t.Errorf("ReadMainMD() = %q, want %q", got, original)
	}
}

// TestAdapter_WriteNewMD_CreatesDir verifies ReadMainMD creates
// directory if it doesn't exist
func TestAdapter_WriteNewMD_CreatesDir(t *testing.T) {
	a, homeDir := newTestAdapter(t)

	configDir := filepath.Join(homeDir, ".claude")
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Fatal("config dir should not exist yet")
	}

	if err := a.WriteNewMD("test content"); err != nil {
		t.Fatalf("WriteNewMD() unexpected error: %v", err)
	}

	if _, err := os.Stat(configDir); err != nil {
		t.Errorf("config dir should have been created: %v", err)
	}
}

// TestAdapter_ConfigFiles verifies ConfigFiles returns at least
// one path not empty.
func TestAdapter_ConfigFiles(t *testing.T) {
	a, _ := newTestAdapter(t)
	files := a.ConfigFiles()

	if len(files) == 0 {
		t.Fatal("ConfigFiles() should return at least one file")
	}

	for _, f := range files {
		if f.Path == "" {
			t.Error("ConfigFile.Path should not be empty")
		}
		if f.Description == "" {
			t.Error("ConfigFile.Description should not be empty")
		}
	}
}