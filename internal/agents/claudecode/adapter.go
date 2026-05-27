// Package claudecode implements Claude Code adapter.

package claudecode

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codebymaribel/eva-ai/internal/agents"
)

const configDirName = ".claude"

const mainConfigFile = "CLAUDE.md"

type Adapter struct {
	homeDir string
}

func New(homeDir string) *Adapter {
	return &Adapter{homeDir: homeDir}
}

func (a *Adapter) Name() string {
	return "claude-code"
}

func (a *Adapter) ConfigDir() string {
	return filepath.Join(a.homeDir, configDirName)
}

func (a *Adapter) MainConfigPath() string {
	return filepath.Join(a.ConfigDir(), mainConfigFile)
}

func (a *Adapter) ConfigFiles() []agents.ConfigFile {
	return []agents.ConfigFile{
		{
			Path:        a.MainConfigPath(),
			Required:    false, // CLAUDE.md can't exist yet
			Description: "Claude Code global instructions",
		},
	}
}

// IsInstalled verifies if Claude Code is installed on the user system.
// Implements agents.Agent.
//
// Strategy:
//  1. Find "claude" binary on PATH (common installation)
//  2. If eva cant find it, verifies if claude config directory exists
//     (in case it's installed but not on path)

func (a *Adapter) IsInstalled() bool {
	if _, err := exec.LookPath("claude"); err == nil {
		return true
	}

	_, err := os.Stat(a.ConfigDir())
	return err == nil
}

// Validate verifica que el adapter pueda operar correctamente.
// CanAdapterRun verifies adapter can run.
// Implements agents.Agent.
//
// Verifies:
//   - That Claude Code is installed
//   - That the config directory is accesible

func (a *Adapter) CanAdapterRun() error {
	if !a.IsInstalled() {
		return &agents.ErrNotInstalled{AgentName: a.Name()}
	}
	return a.validateConfigDir()
}

// Validate implements agents.Agent.
func (a *Adapter) Validate() error {
	return a.CanAdapterRun()
}

func (a *Adapter) validateConfigDir() error {
	info, err := os.Stat(a.ConfigDir())
	if err != nil {
		if os.IsNotExist(err) {
			// If directory doesn't exist we'll create it on the pipeline
			return nil
		}
		return &agents.ErrInvalidConfig{
			AgentName: a.Name(),
			Reason:    fmt.Sprintf("cannot access config directory: %v", err),
		}
	}

	// Verifies that is a directory and not a file
	if !info.IsDir() {
		return &agents.ErrInvalidConfig{
			AgentName: a.Name(),
			Reason:    fmt.Sprintf("%s exists but is not a directory", a.ConfigDir()),
		}
	}

	return nil
}

// EnsureConfigDir creates the directory configuration if it doesn't exist
// Uses os.MkdirAll (mkdir -p equivalent in bash)
// creates all necessary directories

func (a *Adapter) EnsureConfigDir() error {
	if err := os.MkdirAll(a.ConfigDir(), 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", a.ConfigDir(), err)
	}
	return nil
}

// ReadMainMD reads CLAUDE.md current content
// Returns an empty string if file doesn't exist yet
// If it's first install - its not an error.

func (a *Adapter) ReadMainMD() (string, error) {
	content, err := os.ReadFile(a.MainConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // first install, file dont exist yet
		}
		return "", fmt.Errorf("failed to read %s: %w", a.MainConfigPath(), err)
	}
	return string(content), nil
}

// WriteNewMD writes content on CLAUDE.md
// Creates the file if it doesn't exist. If it does it replaces it.
// 0644 mode: owner can read/write, the rest can only read.

func (a *Adapter) WriteNewMD(content string) error {
	if err := a.EnsureConfigDir(); err != nil {
		return err
	}

	if err := os.WriteFile(a.MainConfigPath(), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", a.MainConfigPath(), err)
	}

	return nil
}
