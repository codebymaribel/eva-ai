package system_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/codebymaribel/eva-ai/internal/system"
)

// TestDetect_ReturnsValidPlatform verifies that Detect() always returns
// a non-nil Platform with a home directory on the current machine.
func TestDetect_ReturnsValidPlatform(t *testing.T) {
	p, err := system.Detect()
	if err != nil {
		t.Fatalf("Detect() returned unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("Detect() returned nil platform")
	}
	if p.HomeDir == "" {
		t.Error("HomeDir should not be empty")
	}
	if p.Arch == "" {
		t.Error("Arch should not be empty")
	}
}

// TestParseOSRelease_ValidFile checks that a well-formed os-release file
// is parsed correctly, picking up both ID and ID_LIKE fields.
func TestParseOSRelease_ValidFile(t *testing.T) {
	// Write a temp fake os-release file
	dir := t.TempDir()
	path := filepath.Join(dir, "os-release")

	content := `
# This is a comment
ID="ubuntu"
ID_LIKE="debian"
VERSION_ID="22.04"
PRETTY_NAME="Ubuntu 22.04 LTS"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// We expose parseOSRelease via an exported helper for testing.
	// See detect_test_helpers.go
	fields, err := system.ParseOSReleaseFile(path)
	if err != nil {
		t.Fatalf("ParseOSReleaseFile() error: %v", err)
	}

	tests := []struct {
		key  string
		want string
	}{
		{"ID", "ubuntu"},
		{"ID_LIKE", "debian"},
		{"VERSION_ID", "22.04"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := fields[tt.key]
			if !ok {
				t.Errorf("key %q not found in parsed fields", tt.key)
			}
			if got != tt.want {
				t.Errorf("fields[%q] = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

// TestParseOSRelease_EmptyFile returns an empty map without error.
func TestParseOSRelease_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "os-release")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	fields, err := system.ParseOSReleaseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("expected empty map, got %v", fields)
	}
}

// TestParseOSRelease_MissingFile returns an error for non-existent paths.
func TestParseOSRelease_MissingFile(t *testing.T) {
	_, err := system.ParseOSReleaseFile("/tmp/does-not-exist-xyz/os-release")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

// TestPlatform_String_Linux checks the String() format for a Linux platform.
func TestPlatform_String_Linux(t *testing.T) {
	p := &system.Platform{
		OS:             system.OSLinux,
		Distro:         system.DistroUbuntu,
		Arch:           "amd64",
		PackageManager: system.PkgApt,
		HomeDir:        "/home/user",
	}

	got := p.String()
	want := "linux/ubuntu (amd64) [apt]"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// TestPlatform_String_Darwin checks the String() format for macOS.
func TestPlatform_String_Darwin(t *testing.T) {
	p := &system.Platform{
		OS:             system.OSDarwin,
		Arch:           "arm64",
		PackageManager: system.PkgBrew,
		HomeDir:        "/Users/user",
	}

	got := p.String()
	want := "darwin (arm64) [brew]"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

// TestPlatform_IsSupported covers both supported and unsupported cases.
func TestPlatform_IsSupported(t *testing.T) {
	tests := []struct {
		name string
		p    system.Platform
		want bool
	}{
		{
			name: "macOS with brew",
			p:    system.Platform{OS: system.OSDarwin, PackageManager: system.PkgBrew},
			want: true,
		},
		{
			name: "ubuntu with apt",
			p:    system.Platform{OS: system.OSLinux, PackageManager: system.PkgApt},
			want: true,
		},
		{
			name: "unknown OS",
			p:    system.Platform{OS: system.OSUnknown, PackageManager: system.PkgNone},
			want: false,
		},
		{
			name: "linux no package manager",
			p:    system.Platform{OS: system.OSLinux, PackageManager: system.PkgNone},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.IsSupported()
			if got != tt.want {
				t.Errorf("IsSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}