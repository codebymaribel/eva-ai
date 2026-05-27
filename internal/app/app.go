package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codebymaribel/eva-ai/internal/agents"
	"github.com/codebymaribel/eva-ai/internal/agents/claudecode"
	"github.com/codebymaribel/eva-ai/internal/cli"
	"github.com/codebymaribel/eva-ai/internal/components/sdd"
	"github.com/codebymaribel/eva-ai/internal/components/skills"
	"github.com/codebymaribel/eva-ai/internal/scanner"
	"github.com/codebymaribel/eva-ai/internal/system"
	"github.com/spf13/cobra"
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
	root.AddCommand(newInitCmd())
	root.AddCommand(newSkillCmd())
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

	// 4. Run install for each agent
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	for _, agent := range agents {
		adapter := resolveAgentAdapter(agent, homeDir)
		if adapter == nil {
			fmt.Printf("\n⚠️  %s: adapter not implemented yet — skipping\n", agent)
			continue
		}

		if !adapter.IsInstalled() {
			fmt.Printf("\n⚠️  %s: not installed — skipping\n", agent)
			continue
		}

		fmt.Printf("\n🔧 Configuring %s...\n", agent)

		// Read existing config
		existing, err := adapter.ReadMainMD()
		if err != nil {
			return fmt.Errorf("failed to read %s config: %w", agent, err)
		}

		// Build content from selected components
		var blocks []string
		for _, comp := range opts.Components {
			block, err := renderComponent(comp)
			if err != nil {
				return fmt.Errorf("failed to render component %s: %w", comp, err)
			}
			if block != "" {
				blocks = append(blocks, block)
			}
		}

		// Combine: existing content + new component blocks
		combined := existing
		for _, block := range blocks {
			combined += "\n\n" + block
		}

		// Write to agent config
		if err := adapter.WriteNewMD(combined); err != nil {
			return fmt.Errorf("failed to write %s config: %w", agent, err)
		}

		fmt.Printf("   ✅ %s configured\n", adapter.Name())
	}

	fmt.Println("\n✅ Install complete!")
	return nil
}

// resolveAgentAdapter returns the adapter for a given agent, or nil if not implemented.
func resolveAgentAdapter(agent cli.Agent, homeDir string) agents.InjectableAgent {
	switch agent {
	case cli.AgentClaudeCode:
		return claudecode.New(homeDir)
	default:
		return nil
	}
}

// renderComponent renders the content for a given component.
func renderComponent(comp cli.Component) (string, error) {
	switch comp {
	case cli.ComponentSDD:
		c := sdd.New()
		return c.Render(), nil
	case cli.ComponentSkills:
		c := skills.New()
		return c.Render(), nil
	case cli.ComponentMCP:
		return "", nil // not implemented yet
	case cli.ComponentPersona:
		return "", nil // not implemented yet
	default:
		return "", fmt.Errorf("unknown component: %s", comp)
	}
}

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage skills — add, list, or remove",
	}

	cmd.AddCommand(newSkillAddCmd())

	return cmd
}

func newSkillAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path-or-url>",
		Short: "Add a skill from a local file or URL",
		Example: `  # From a local file
  eva skill add ./my-skill/SKILL.md

  # From a URL
  eva skill add https://raw.githubusercontent.com/user/repo/main/SKILL.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillAdd(args[0])
		},
	}

	return cmd
}

func runSkillAdd(source string) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Check if .eva/ exists
	evaDir := filepath.Join(projectDir, ".eva")
	if _, err := os.Stat(evaDir); os.IsNotExist(err) {
		return fmt.Errorf(".eva/ not found — run `eva init` first")
	}

	skillsComp := skills.New()
	var skillName string

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		fmt.Printf("⬇️  Downloading skill from %s...\n", source)
		skillName, err = skillsComp.AddFromURL(projectDir, source)
		if err != nil {
			return err
		}
	} else {
		skillName, err = skillsComp.AddFromPath(projectDir, source)
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ Skill saved: .eva/skills/%s/SKILL.md\n", skillName)

	// Update .eva/skills/README.md
	if err := skillsComp.RefreshREADME(projectDir); err != nil {
		fmt.Printf("⚠️  Could not update skills README: %v\n", err)
	} else {
		fmt.Println("📄 Updated .eva/skills/README.md")
	}

	// Inject into existing agent configs so the agent sees it automatically
	injected := skillsComp.InjectIntoAgentConfigs(projectDir, skillName)
	if len(injected) > 0 {
		fmt.Println("\n📌 Injected into agent configs:")
		for _, path := range injected {
			rel, _ := filepath.Rel(projectDir, path)
			fmt.Printf("   ✅ %s\n", rel)
		}
	} else {
		fmt.Println("\n💡 No agent configs found. Run `eva install --agent <agent>` first,")
		fmt.Println("   then re-run `eva skill add` to inject into the agent.")
	}

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

func newInitCmd() *cobra.Command {
	var (
		noSkills bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize .eva/ in the current project — scan codebase, create core skills, generate project skills",
		Example: `  # Full init — scan project and create all skills
  eva-ai init

  # Init without skills — only SDD structure
  eva-ai init --no-skills`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(noSkills)
		},
	}

	cmd.Flags().BoolVar(&noSkills, "no-skills", false, "Skip skill generation — only create .eva/ structure")

	return cmd
}

func runInit(noSkills bool) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	fmt.Println("🔍 Scanning project...")

	// 1. Detect platform (for future use)
	platform, err := system.Detect()
	if err != nil {
		return fmt.Errorf("platform detection failed: %w", err)
	}
	fmt.Printf("🖥️  Platform: %s\n", platform.String())

	// 2. Scan project stack and patterns
	info, err := scanner.Scan(projectDir)
	if err != nil {
		return fmt.Errorf("project scan failed: %w", err)
	}

	fmt.Printf("📦 Stack: %s", info.Stack.Language)
	if info.Stack.Framework != "" {
		fmt.Printf(" + %s", info.Stack.Framework)
	}
	fmt.Println()

	if len(info.Patterns) > 0 {
		fmt.Println("🏗️  Patterns:")
		for _, p := range info.Patterns {
			fmt.Printf("   - %s (%s)\n", p.Name, p.Confidence)
		}
	}

	// 3. Initialize .eva/ directory structure
	fmt.Println("\n📁 Creating .eva/ structure...")
	sddComp := sdd.New()
	if err := sddComp.InitEvaDir(projectDir); err != nil {
		return fmt.Errorf("failed to init .eva/: %w", err)
	}
	fmt.Println("   ✅ .eva/README.md")
	fmt.Println("   ✅ .eva/memory.md")
	fmt.Println("   ✅ .eva/phases/")

	// 4. Initialize core skills
	if !noSkills {
		fmt.Println("\n📚 Setting up skills...")
		skillsComp := skills.New()

		if err := skillsComp.InitSkillsDir(projectDir); err != nil {
			return fmt.Errorf("failed to init skills: %w", err)
		}
		fmt.Println("   ✅ .eva/skills/architecture/SKILL.md")
		fmt.Println("   ✅ .eva/skills/testing/SKILL.md")
		fmt.Println("   ✅ .eva/skills/git-workflow/SKILL.md")

		// 5. Generate project-specific skills from scan results
		fmt.Println("\n🔬 Generating project skills...")
		if err := generateProjectSkills(projectDir, info, skillsComp); err != nil {
			return fmt.Errorf("failed to generate project skills: %w", err)
		}

		// 6. Update .eva/skills/README.md with all skills
		if err := updateSkillsREADME(projectDir, info); err != nil {
			return fmt.Errorf("failed to update skills README: %w", err)
		}
	}

	fmt.Println("\n✅ .eva/ initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review .eva/skills/ — update generated skills with your conventions")
	fmt.Println("  2. Run `eva install --agent <agent>` to inject into your AI agent")
	fmt.Println("  3. Start a mission with `/eva [task]` in your agent")

	return nil
}

// generateProjectSkills creates project-specific SKILL.md files based on scan results.
func generateProjectSkills(projectDir string, info *scanner.ProjectInfo, comp *skills.Component) error {
	// Always generate a domain skill
	domainContent := scanner.GenerateSkillContent(info,
		"domain",
		"Business logic, models, entities, and domain-specific rules.",
		"business logic, domain, models, entities, validation, rules",
	)
	if err := writeSkill(projectDir, "domain", domainContent); err != nil {
		return err
	}
	fmt.Println("   ✅ .eva/skills/domain/SKILL.md")

	// Generate API/HTTP skill if the project uses HTTP
	if info.HTTPClient != "" {
		apiContent := scanner.GenerateSkillContent(info,
			"api",
			fmt.Sprintf("HTTP client (%s), API integration, auth, and request patterns.", info.HTTPClient),
			"api, http, client, request, auth, endpoint, rest, graphql",
		)
		if err := writeSkill(projectDir, "api", apiContent); err != nil {
			return err
		}
		fmt.Println("   ✅ .eva/skills/api/SKILL.md")
	}

	// Generate stack-specific skill if a framework was detected
	if info.Stack.Framework != "" {
		stackContent := scanner.GenerateSkillContent(info,
			info.Stack.Framework,
			fmt.Sprintf("Framework-specific patterns and conventions for %s.", info.Stack.Framework),
			fmt.Sprintf("%s, framework, components, routing, configuration", info.Stack.Framework),
		)
		if err := writeSkill(projectDir, info.Stack.Framework, stackContent); err != nil {
			return err
		}
		fmt.Printf("   ✅ .eva/skills/%s/SKILL.md\n", info.Stack.Framework)
	}

	return nil
}

// writeSkill writes a SKILL.md file to .eva/skills/<name>/
func writeSkill(projectDir, name, content string) error {
	skillDir := filepath.Join(projectDir, ".eva", "skills", name)
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644)
}

// updateSkillsREADME regenerates .eva/skills/README.md with all available skills.
func updateSkillsREADME(projectDir string, info *scanner.ProjectInfo) error {
	var b strings.Builder

	b.WriteString("# Project Skills\n\n")
	b.WriteString("Skills are loaded on demand by SDD agents.\n")
	b.WriteString("Each skill defines when it should be loaded via its `trigger` field.\n\n")
	b.WriteString("## Available Skills\n\n")
	b.WriteString("| Skill | Trigger |\n")
	b.WriteString("|---|---|\n")

	// Core skills
	b.WriteString("| `architecture` | Stack, patterns, state, technical decisions |\n")
	b.WriteString("| `testing` | Tests, coverage, mocks, fixtures |\n")
	b.WriteString("| `git-workflow` | Git, commits, branches, PRs, releases |\n")

	// Project skills
	b.WriteString("| `domain` | Business logic, models, entities |\n")
	if info.HTTPClient != "" {
		b.WriteString(fmt.Sprintf("| `api` | HTTP client (%s), API integration, auth |\n", info.HTTPClient))
	}
	if info.Stack.Framework != "" {
		b.WriteString(fmt.Sprintf("| `%s` | Framework-specific patterns for %s |\n", info.Stack.Framework, info.Stack.Framework))
	}

	b.WriteString("\n### Adding a skill\n\n")
	b.WriteString("```bash\n")
	b.WriteString("# From a local file\n")
	b.WriteString("eva skill add ./path/to/SKILL.md\n\n")
	b.WriteString("# From a URL\n")
	b.WriteString("eva skill add https://raw.githubusercontent.com/.../SKILL.md\n")
	b.WriteString("```\n")

	readmePath := filepath.Join(projectDir, ".eva", "skills", "README.md")
	return os.WriteFile(readmePath, []byte(b.String()), 0644)
}
