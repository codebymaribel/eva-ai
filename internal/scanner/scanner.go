// Package scanner detects the project stack, patterns, and conventions
// by reading files in the project directory. Used by eva init to generate
// project-specific skills.
package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Stack represents the detected technology stack.
type Stack struct {
	Language    string // go, typescript, javascript, python, rust
	Framework   string // next, react, react-native, angular, vue, gin, echo, fastapi
	Runtime     string // node, deno, bun (empty for non-JS)
	PackageMngr string // npm, yarn, pnpm, pip, cargo, go
}

// Pattern represents a detected architectural pattern.
type Pattern struct {
	Name        string // e.g. "clean-architecture", "feature-based"
	Description string
	Confidence  string // high, medium, low
}

// ProjectInfo holds all detected information about a project.
type ProjectInfo struct {
	RootDir    string
	Stack      Stack
	Patterns   []Pattern
	HasTests   bool
	TestRunner string // jest, vitest, go test, pytest, cargo test
	HasCI      bool
	HasLint    bool
	Linter     string // eslint, golangci-lint, ruff, clippy
	HasTypeChk bool   // TypeScript strict, mypy, etc.
	HTTPClient string // fetch, axios, net/http, requests
	StateMgmt  string // zustand, redux, context, pinia
	Routing    string // next-router, react-router, vue-router, gin, echo
}

// Scan analyzes the project directory and returns detected information.
func Scan(projectDir string) (*ProjectInfo, error) {
	info := &ProjectInfo{
		RootDir: projectDir,
	}

	// Detect stack from config files
	info.Stack = detectStack(projectDir)

	// Detect patterns from directory structure
	info.Patterns = detectPatterns(projectDir)

	// Detect tooling
	info.HasTests, info.TestRunner = detectTesting(projectDir)
	info.HasCI = detectCI(projectDir)
	info.HasLint, info.Linter = detectLinting(projectDir)
	info.HasTypeChk = detectTypeChecking(projectDir)

	// Detect specific patterns from dependencies
	if info.Stack.Language == "typescript" || info.Stack.Language == "javascript" {
		info.HTTPClient, info.StateMgmt, info.Routing = detectJSPatterns(projectDir, info.Stack.Framework)
	} else if info.Stack.Language == "go" {
		info.HTTPClient, info.StateMgmt, info.Routing = detectGoPatterns(projectDir)
	}

	return info, nil
}

// detectStack reads config files to determine the language and framework.
func detectStack(dir string) Stack {
	s := Stack{}

	// Go
	if exists(filepath.Join(dir, "go.mod")) {
		s.Language = "go"
		s.PackageMngr = "go"
		s.Framework = detectGoFramework(dir)
		return s
	}

	// Rust
	if exists(filepath.Join(dir, "Cargo.toml")) {
		s.Language = "rust"
		s.PackageMngr = "cargo"
		s.Framework = detectRustFramework(dir)
		return s
	}

	// Python
	if exists(filepath.Join(dir, "pyproject.toml")) || exists(filepath.Join(dir, "requirements.txt")) {
		s.Language = "python"
		s.PackageMngr = "pip"
		s.Framework = detectPythonFramework(dir)
		return s
	}

	// JavaScript/TypeScript — check package.json
	if exists(filepath.Join(dir, "package.json")) {
		s.Language = "javascript"
		s.Runtime = "node"
		s.PackageMngr = detectJSPackageManager(dir)

		if exists(filepath.Join(dir, "tsconfig.json")) {
			s.Language = "typescript"
		}

		s.Framework = detectJSFramework(dir)
		return s
	}

	return s
}

// detectGoFramework checks go.mod for known frameworks.
func detectGoFramework(dir string) string {
	content, err := readFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	switch {
	case strings.Contains(content, "github.com/gin-gonic/gin"):
		return "gin"
	case strings.Contains(content, "github.com/labstack/echo"):
		return "echo"
	case strings.Contains(content, "github.com/gofiber/fiber"):
		return "fiber"
	case strings.Contains(content, "net/http"):
		return "net-http"
	default:
		return ""
	}
}

// detectRustFramework checks Cargo.toml for known frameworks.
func detectRustFramework(dir string) string {
	content, err := readFile(filepath.Join(dir, "Cargo.toml"))
	if err != nil {
		return ""
	}
	switch {
	case strings.Contains(content, "actix-web"):
		return "actix"
	case strings.Contains(content, "axum"):
		return "axum"
	case strings.Contains(content, "rocket"):
		return "rocket"
	default:
		return ""
	}
}

// detectPythonFramework checks for known Python frameworks.
func detectPythonFramework(dir string) string {
	content := ""
	for _, f := range []string{"pyproject.toml", "requirements.txt"} {
		c, err := readFile(filepath.Join(dir, f))
		if err == nil {
			content += c
		}
	}
	switch {
	case strings.Contains(content, "fastapi"):
		return "fastapi"
	case strings.Contains(content, "django"):
		return "django"
	case strings.Contains(content, "flask"):
		return "flask"
	default:
		return ""
	}
}

// detectJSFramework checks package.json for known JS/TS frameworks.
func detectJSFramework(dir string) string {
	pkg := readPackageJSON(dir)
	deps := mergeDeps(pkg)

	switch {
	case hasDep(deps, "next"):
		return "next"
	case hasDep(deps, "react-native") || hasDep(deps, "expo"):
		return "react-native"
	case hasDep(deps, "vue"):
		return "vue"
	case hasDep(deps, "@angular/core"):
		return "angular"
	case hasDep(deps, "svelte") || hasDep(deps, "@sveltejs/kit"):
		return "svelte"
	case hasDep(deps, "react"):
		return "react"
	default:
		return ""
	}
}

// detectJSPackageManager checks lockfiles to determine the package manager.
func detectJSPackageManager(dir string) string {
	switch {
	case exists(filepath.Join(dir, "bun.lockb")):
		return "bun"
	case exists(filepath.Join(dir, "pnpm-lock.yaml")):
		return "pnpm"
	case exists(filepath.Join(dir, "yarn.lock")):
		return "yarn"
	default:
		return "npm"
	}
}

// detectPatterns looks at directory structure for architectural patterns.
func detectPatterns(dir string) []Pattern {
	var patterns []Pattern

	entries, err := os.ReadDir(dir)
	if err != nil {
		return patterns
	}

	dirs := make(map[string]bool)
	for _, e := range entries {
		if e.IsDir() {
			dirs[strings.ToLower(e.Name())] = true
		}
	}

	// Clean/Hexagonal architecture
	if dirs["domain"] && dirs["infrastructure"] {
		patterns = append(patterns, Pattern{
			Name:        "clean-architecture",
			Description: "Domain + Infrastructure separation",
			Confidence:  "high",
		})
	} else if dirs["domain"] && dirs["data"] {
		patterns = append(patterns, Pattern{
			Name:        "clean-architecture",
			Description: "Domain + Data layer separation",
			Confidence:  "medium",
		})
	}

	// Feature-based structure
	if dirs["features"] || dirs["modules"] {
		patterns = append(patterns, Pattern{
			Name:        "feature-based",
			Description: "Code organized by feature/module",
			Confidence:  "high",
		})
	}

	// Internal package structure (Go)
	if dirs["internal"] {
		patterns = append(patterns, Pattern{
			Name:        "internal-packages",
			Description: "Go internal package convention",
			Confidence:  "high",
		})
	}

	// src-based structure
	if dirs["src"] {
		patterns = append(patterns, Pattern{
			Name:        "src-based",
			Description: "Source code in src/ directory",
			Confidence:  "high",
		})
	}

	return patterns
}

// detectTesting checks for test configuration and files.
func detectTesting(dir string) (bool, string) {
	// Go — always has go test
	if exists(filepath.Join(dir, "go.mod")) {
		return true, "go test"
	}

	// Rust — always has cargo test
	if exists(filepath.Join(dir, "Cargo.toml")) {
		return true, "cargo test"
	}

	// Python
	if exists(filepath.Join(dir, "pytest.ini")) || exists(filepath.Join(dir, "conftest.py")) {
		return true, "pytest"
	}

	// JS/TS — check package.json scripts
	pkg := readPackageJSON(dir)
	scripts, _ := pkg["scripts"].(map[string]interface{})
	if _, ok := scripts["test"]; ok {
		deps := mergeDeps(pkg)
		switch {
		case hasDep(deps, "vitest"):
			return true, "vitest"
		case hasDep(deps, "jest"):
			return true, "jest"
		default:
			return true, "npm test"
		}
	}

	return false, ""
}

// detectCI checks for CI configuration files.
func detectCI(dir string) bool {
	ciFiles := []string{
		".github/workflows",
		".gitlab-ci.yml",
		"Jenkinsfile",
		".circleci/config.yml",
		".travis.yml",
	}
	for _, f := range ciFiles {
		if exists(filepath.Join(dir, f)) {
			return true
		}
	}
	return false
}

// detectLinting checks for linter configuration.
func detectLinting(dir string) (bool, string) {
	// Go
	if exists(filepath.Join(dir, ".golangci.yml")) || exists(filepath.Join(dir, ".golangci.yaml")) {
		return true, "golangci-lint"
	}

	// Python
	if exists(filepath.Join(dir, "ruff.toml")) || exists(filepath.Join(dir, ".ruff.toml")) {
		return true, "ruff"
	}

	// JS/TS
	if exists(filepath.Join(dir, ".eslintrc.js")) || exists(filepath.Join(dir, ".eslintrc.json")) || exists(filepath.Join(dir, "eslint.config.js")) || exists(filepath.Join(dir, "eslint.config.mjs")) {
		return true, "eslint"
	}

	// Rust
	if exists(filepath.Join(dir, "clippy.toml")) || exists(filepath.Join(dir, ".clippy.toml")) {
		return true, "clippy"
	}

	return false, ""
}

// detectTypeChecking checks for type checking configuration.
func detectTypeChecking(dir string) bool {
	if exists(filepath.Join(dir, "tsconfig.json")) {
		return true
	}
	// Python mypy
	if exists(filepath.Join(dir, "mypy.ini")) || exists(filepath.Join(dir, ".mypy.ini")) {
		return true
	}
	return false
}

// detectJSPatterns reads package.json to detect HTTP client, state management, and routing.
func detectJSPatterns(dir, framework string) (httpClient, stateMgmt, routing string) {
	pkg := readPackageJSON(dir)
	deps := mergeDeps(pkg)

	// HTTP client
	switch {
	case hasDep(deps, "axios"):
		httpClient = "axios"
	case hasDep(deps, "ky"):
		httpClient = "ky"
	default:
		httpClient = "fetch"
	}

	// State management
	switch {
	case hasDep(deps, "zustand"):
		stateMgmt = "zustand"
	case hasDep(deps, "jotai"):
		stateMgmt = "jotai"
	case hasDep(deps, "@reduxjs/toolkit") || hasDep(deps, "redux"):
		stateMgmt = "redux"
	case hasDep(deps, "pinia"):
		stateMgmt = "pinia"
	case hasDep(deps, "vuex"):
		stateMgmt = "vuex"
	case hasDep(deps, "@tanstack/react-query"):
		stateMgmt = "react-query"
	default:
		stateMgmt = "context"
	}

	// Routing
	switch framework {
	case "next":
		routing = "next-router"
	case "react-native":
		if hasDep(deps, "expo-router") {
			routing = "expo-router"
		} else if hasDep(deps, "@react-navigation/native") {
			routing = "react-navigation"
		}
	case "vue":
		if hasDep(deps, "vue-router") {
			routing = "vue-router"
		}
	case "angular":
		routing = "angular-router"
	case "react":
		if hasDep(deps, "react-router-dom") || hasDep(deps, "react-router") {
			routing = "react-router"
		}
	}

	return
}

// detectGoPatterns checks go.mod for HTTP framework and common patterns.
func detectGoPatterns(dir string) (httpClient, stateMgmt, routing string) {
	content, err := readFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", "", ""
	}

	// Go uses net/http by default for HTTP client
	httpClient = "net/http"

	// Routing is usually part of the framework
	switch {
	case strings.Contains(content, "github.com/gin-gonic/gin"):
		routing = "gin"
	case strings.Contains(content, "github.com/labstack/echo"):
		routing = "echo"
	case strings.Contains(content, "github.com/gofiber/fiber"):
		routing = "fiber"
	default:
		routing = "net/http"
	}

	return
}

// — helpers —

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type packageJSON map[string]interface{}

func readPackageJSON(dir string) packageJSON {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}
	return pkg
}

func mergeDeps(pkg packageJSON) map[string]bool {
	deps := make(map[string]bool)
	for _, key := range []string{"dependencies", "devDependencies", "peerDependencies"} {
		if m, ok := pkg[key].(map[string]interface{}); ok {
			for name := range m {
				deps[name] = true
			}
		}
	}
	return deps
}

func hasDep(deps map[string]bool, name string) bool {
	return deps[name]
}

// GenerateSkillContent creates a SKILL.md for a given topic based on project info.
func GenerateSkillContent(info *ProjectInfo, topic, description, trigger string) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("name: %s\n", topic))
	b.WriteString(fmt.Sprintf("description: %s\n", description))
	b.WriteString("sdd_phases: [scan, schematics, execute, debrief]\n")
	b.WriteString(fmt.Sprintf("trigger: %q\n", trigger))
	b.WriteString("---\n\n")
	b.WriteString(fmt.Sprintf("# %s\n\n", strings.Title(topic)))

	b.WriteString("## Stack\n")
	b.WriteString(fmt.Sprintf("- Language: %s\n", info.Stack.Language))
	if info.Stack.Framework != "" {
		b.WriteString(fmt.Sprintf("- Framework: %s\n", info.Stack.Framework))
	}
	if info.Stack.Runtime != "" {
		b.WriteString(fmt.Sprintf("- Runtime: %s\n", info.Stack.Runtime))
	}
	b.WriteString(fmt.Sprintf("- Package manager: %s\n", info.Stack.PackageMngr))
	b.WriteString("\n")

	b.WriteString("## Patterns\n")
	if len(info.Patterns) > 0 {
		for _, p := range info.Patterns {
			b.WriteString(fmt.Sprintf("- **%s**: %s (confidence: %s)\n", p.Name, p.Description, p.Confidence))
		}
	} else {
		b.WriteString("_No architectural patterns detected. Update this section manually._\n")
	}
	b.WriteString("\n")

	b.WriteString("## Rules\n")
	b.WriteString("- Follow existing patterns before introducing new ones\n")
	b.WriteString("- Keep consistent with the detected project conventions\n")
	b.WriteString("\n")

	b.WriteString("## Do not\n")
	b.WriteString("- Introduce patterns that conflict with the existing architecture\n")
	b.WriteString("- Mix conventions from different parts of the codebase\n")
	b.WriteString("\n")

	b.WriteString("## References\n")
	b.WriteString("- _Add links to relevant project documentation here_\n")

	return b.String()
}
