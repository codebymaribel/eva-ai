package agents

import "fmt"

// Agent is the contract every adapter must fulfill.
type Agent interface {
	Name()        string
	ConfigDir()   string
	IsInstalled() bool
	Validate()    error
}

// ConfigFile represents a configuration file an adapter manages.
type ConfigFile struct {
	Path        string
	Required    bool
	Description string
}

// ErrNotInstalled is returned when an agent is not installed on the system.
type ErrNotInstalled struct {
	AgentName string
}

func (e *ErrNotInstalled) Error() string {
	return fmt.Sprintf("agent %q is not installed — install it first and re-run", e.AgentName)
}

// ErrInvalidConfig is returned when the agent config exists but is invalid.
type ErrInvalidConfig struct {
	AgentName string
	Reason    string
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("agent %q has invalid config: %s", e.AgentName, e.Reason)
}
