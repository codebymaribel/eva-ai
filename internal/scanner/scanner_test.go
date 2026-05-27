package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_GoProject(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.26\n\nrequire github.com/gin-gonic/gin v1.9.0\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.Stack.Language != "go" {
		t.Errorf("Language = %q, want %q", info.Stack.Language, "go")
	}
	if info.Stack.PackageMngr != "go" {
		t.Errorf("PackageMngr = %q, want %q", info.Stack.PackageMngr, "go")
	}
	if info.Stack.Framework != "gin" {
		t.Errorf("Framework = %q, want %q", info.Stack.Framework, "gin")
	}
}

func TestScan_TypeScriptNextProject(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"next":"14.0.0","react":"18.0.0"},"devDependencies":{"typescript":"5.0.0","eslint":"8.0.0","jest":"29.0.0"}}`)
	writeFile(t, dir, "tsconfig.json", `{}`)
	writeFile(t, dir, "yarn.lock", ``)
	writeFile(t, dir, "eslint.config.js", `export default []`)

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.Stack.Language != "typescript" {
		t.Errorf("Language = %q, want %q", info.Stack.Language, "typescript")
	}
	if info.Stack.Framework != "next" {
		t.Errorf("Framework = %q, want %q", info.Stack.Framework, "next")
	}
	if info.Stack.PackageMngr != "yarn" {
		t.Errorf("PackageMngr = %q, want %q", info.Stack.PackageMngr, "yarn")
	}
	if info.HTTPClient != "fetch" {
		t.Errorf("HTTPClient = %q, want %q", info.HTTPClient, "fetch")
	}
	if !info.HasLint || info.Linter != "eslint" {
		t.Errorf("Linter = %q/%v, want eslint/true", info.Linter, info.HasLint)
	}
	if !info.HasTypeChk {
		t.Errorf("HasTypeChk = false, want true")
	}
}

func TestScan_PythonFastAPIProject(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "pyproject.toml", "[project]\nname = \"myapp\"\ndependencies = [\"fastapi\"]\n")
	writeFile(t, dir, "conftest.py", "")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.Stack.Language != "python" {
		t.Errorf("Language = %q, want %q", info.Stack.Language, "python")
	}
	if info.Stack.Framework != "fastapi" {
		t.Errorf("Framework = %q, want %q", info.Stack.Framework, "fastapi")
	}
	if !info.HasTests || info.TestRunner != "pytest" {
		t.Errorf("TestRunner = %q/%v, want pytest/true", info.TestRunner, info.HasTests)
	}
}

func TestScan_RustAxumProject(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "Cargo.toml", "[package]\nname = \"myapp\"\n\n[dependencies]\naxum = \"0.7\"\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.Stack.Language != "rust" {
		t.Errorf("Language = %q, want %q", info.Stack.Language, "rust")
	}
	if info.Stack.Framework != "axum" {
		t.Errorf("Framework = %q, want %q", info.Stack.Framework, "axum")
	}
	if !info.HasTests || info.TestRunner != "cargo test" {
		t.Errorf("TestRunner = %q/%v, want cargo test/true", info.TestRunner, info.HasTests)
	}
}

func TestDetectPatterns_CleanArchitecture(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "domain"), 0755)
	os.MkdirAll(filepath.Join(dir, "infrastructure"), 0755)
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.26\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	found := false
	for _, p := range info.Patterns {
		if p.Name == "clean-architecture" {
			found = true
			if p.Confidence != "high" {
				t.Errorf("clean-architecture confidence = %q, want %q", p.Confidence, "high")
			}
		}
	}
	if !found {
		t.Error("expected clean-architecture pattern to be detected")
	}
}

func TestDetectPatterns_FeatureBased(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "features"), 0755)
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.26\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	found := false
	for _, p := range info.Patterns {
		if p.Name == "feature-based" {
			found = true
		}
	}
	if !found {
		t.Error("expected feature-based pattern to be detected")
	}
}

func TestDetectPatterns_InternalPackages(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "internal"), 0755)
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.26\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	found := false
	for _, p := range info.Patterns {
		if p.Name == "internal-packages" {
			found = true
		}
	}
	if !found {
		t.Error("expected internal-packages pattern to be detected")
	}
}

func TestDetectCI(t *testing.T) {
	tests := []struct {
		name   string
		files  map[string]string
		expect bool
	}{
		{"github-workflows", map[string]string{".github/workflows/ci.yml": ""}, true},
		{"gitlab-ci", map[string]string{".gitlab-ci.yml": ""}, true},
		{"none", map[string]string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for path, content := range tt.files {
				writeFile(t, dir, path, content)
			}
			info, err := Scan(dir)
			if err != nil {
				t.Fatalf("Scan() error: %v", err)
			}
			if info.HasCI != tt.expect {
				t.Errorf("HasCI = %v, want %v", info.HasCI, tt.expect)
			}
		})
	}
}

func TestDetectJSPatterns(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{
		"dependencies": {
			"react": "18.0.0",
			"axios": "1.0.0",
			"zustand": "4.0.0",
			"react-router-dom": "6.0.0"
		}
	}`)

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.HTTPClient != "axios" {
		t.Errorf("HTTPClient = %q, want %q", info.HTTPClient, "axios")
	}
	if info.StateMgmt != "zustand" {
		t.Errorf("StateMgmt = %q, want %q", info.StateMgmt, "zustand")
	}
	if info.Routing != "react-router" {
		t.Errorf("Routing = %q, want %q", info.Routing, "react-router")
	}
}

func TestDetectGoPatterns(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/myapp\n\ngo 1.26\n\nrequire github.com/labstack/echo v4.0.0\n")

	info, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if info.HTTPClient != "net/http" {
		t.Errorf("HTTPClient = %q, want %q", info.HTTPClient, "net/http")
	}
	if info.Routing != "echo" {
		t.Errorf("Routing = %q, want %q", info.Routing, "echo")
	}
}

// helper

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
