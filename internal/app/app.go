package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/codebymaribel/eva-ai/internal/cli"
	"github.com/codebymaribel/eva-ai/internal/system"
)


func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "eva-ai",
		Short: "🚀 Supercharge your AI coding agents with skills, SDD, MCP, and personalization",
		Long: `eva-ai is an ecosystem configurator for AI coding agents.

It injects your custom Skills, SDD workflow, MCP servers, and Personalization
into Claude Code, Cursor, Copilot, OpenCode, and Windsurf — with one command.`,
		SilenceUsage: true,
	}

	root.AddCommand(newInstallCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func newInstallCmd() *cobra.Command {
	var (
		agentFlag     string
		componentFlag string
		presetFlag    string
		dryRun        bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the AI stack into one or more agents",
		Example: `  # Full stack for Claude Code and Cursor
  eva-ai install --agent claude-code,cursor --preset full

  # Preview what would happen without making changes  
  eva-ai install --agent claude-code --dry-run

  # Pick specific components
  eva-ai install --agent cursor --component sdd,skills,personality`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstall(agentFlag, componentFlag, presetFlag, dryRun)
		},
	}

	cmd.Flags().StringVar(&agentFlag, "agent", "", "Agents to configure (comma-separated): claude-code, opencode, cursor, copilot, windsurf")
	cmd.Flags().StringVar(&componentFlag, "component", "", "Components to install (comma-separated): sdd, skills, mcp, persona")
	cmd.Flags().StringVar(&presetFlag, "preset", "full", "Preset: full, minimal, custom")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview the install plan without applying any changes")

	// --agent is required; Cobra will enforce this before RunE is called.
	_ = cmd.MarkFlagRequired("agent")

	return cmd
}

func runInstall(agentFlag, componentFlag, presetFlag string, dryRun bool) error {
	
	agents, err := cli.ParseAgents(agentFlag)
	if err != nil {
		return err
	}

	components, err := cli.ParseComponents(componentFlag)
	if err != nil {
		return err
	}

	preset, err := cli.ParsePreset(presetFlag)
	if err != nil {
		return err
	}

	opts := &cli.InstallOptions{
		Agents:     agents,
		Components: components,
		Preset:     preset,
		DryRun:     dryRun,
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	
	platform, err := system.Detect()
	if err != nil {
		return fmt.Errorf("platform detection failed: %w", err)
	}

	if !platform.IsSupported() {
		return fmt.Errorf("unsupported platform: %s", platform.String())
	}

	// 3. Print what we detected
	fmt.Printf("🖥️  Platform: %s\n", platform.String())
	fmt.Printf("🤖 Agents:   %v\n", opts.Agents)
	fmt.Printf("📦 Components: %v\n", opts.Components)

	if dryRun {
		fmt.Println("\n✅ Dry-run complete — no changes were made.")
		return nil
	}

	// 4. TODO (Phase 2+): run planner → pipeline
	fmt.Println("\n🚧 Pipeline coming in Phase 2...")
	return nil
}


func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("my-ai-stack v0.1.0")
		},
	}
}