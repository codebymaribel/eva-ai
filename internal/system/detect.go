// Package system handles OS and distro detection.
// It determines what platform the user is running on
// so the rest of the tool can make smart decisions
// (e.g. use brew on macOS, apt on Ubuntu, winget on Windows).
package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OS represents the operating system family.
type OS string

const (
	OSDarwin  OS = "darwin"
	OSLinux   OS = "linux"
	OSWindows OS = "windows"
	OSUnknown OS = "unknown"
)

// Distro represents a specific Linux distribution.
type Distro string

const (
	DistroUbuntu   Distro = "ubuntu"
	DistroDebian   Distro = "debian"
	DistroArch     Distro = "arch"
	DistroFedora   Distro = "fedora"
	DistroMint     Distro = "mint"
	DistroPopOS    Distro = "pop"
	DistroManjaro  Distro = "manjaro"
	DistroNone     Distro = ""
	DistroUnknown  Distro = "unknown"
)

// PackageManager represents the system package manager.
type PackageManager string

const (
	PkgBrew   PackageManager = "brew"
	PkgApt    PackageManager = "apt"
	PkgPacman PackageManager = "pacman"
	PkgDnf    PackageManager = "dnf"
	PkgWinget PackageManager = "winget"
	PkgNone   PackageManager = ""
)

// Platform holds all detected information about the current system.
type Platform struct {
	OS             OS
	Distro         Distro
	Arch           string
	PackageManager PackageManager
	HomeDir        string
}

// Detect inspects the current runtime environment and returns a Platform.
// It reads GOOS for the OS family, parses /etc/os-release for Linux distros,
// and resolves the appropriate package manager for the platform.
func Detect() (*Platform, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not resolve home directory: %w", err)
	}

	p := &Platform{
		OS:      detectOS(),
		Arch:    runtime.GOARCH,
		HomeDir: homeDir,
	}

	if p.OS == OSLinux {
		p.Distro = detectLinuxDistro()
	}

	p.PackageManager = resolvePackageManager(p.OS, p.Distro)

	return p, nil
}

// detectOS maps Go's runtime.GOOS to our OS type.
func detectOS() OS {
	switch runtime.GOOS {
	case "darwin":
		return OSDarwin
	case "linux":
		return OSLinux
	case "windows":
		return OSWindows
	default:
		return OSUnknown
	}
}

// detectLinuxDistro reads /etc/os-release to find the distro name.
// It also checks ID_LIKE to handle derivatives like Mint (Ubuntu-based)
// or Manjaro (Arch-based).
func detectLinuxDistro() Distro {
	fields, err := parseOSRelease("/etc/os-release")
	if err != nil {
		return DistroUnknown
	}

	// Check primary ID first
	if id, ok := fields["ID"]; ok {
		if d := mapDistro(id); d != DistroUnknown {
			return d
		}
	}

	// Fall back to ID_LIKE for derivative distros
	// e.g. Linux Mint has ID=linuxmint but ID_LIKE=ubuntu
	if idLike, ok := fields["ID_LIKE"]; ok {
		for _, like := range strings.Fields(idLike) {
			if d := mapDistro(like); d != DistroUnknown {
				return d
			}
		}
	}

	return DistroUnknown
}

// parseOSRelease reads a key=value file (like /etc/os-release)
// and returns a map of the fields.
func parseOSRelease(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fields := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		fields[key] = val
	}

	return fields, scanner.Err()
}

// mapDistro converts a raw distro string to our Distro type.
func mapDistro(raw string) Distro {
	switch strings.ToLower(raw) {
	case "ubuntu":
		return DistroUbuntu
	case "debian":
		return DistroDebian
	case "arch":
		return DistroArch
	case "fedora":
		return DistroFedora
	case "linuxmint", "mint":
		return DistroMint
	case "pop", "pop-os":
		return DistroPopOS
	case "manjaro":
		return DistroManjaro
	default:
		return DistroUnknown
	}
}

// resolvePackageManager returns the best package manager for the platform.
// On Linux it checks whether the binary actually exists on PATH
// before committing to it.
func resolvePackageManager(os OS, distro Distro) PackageManager {
	switch os {
	case OSDarwin:
		return PkgBrew
	case OSWindows:
		return PkgWinget
	case OSLinux:
		return resolveLinuxPkgManager(distro)
	default:
		return PkgNone
	}
}

// resolveLinuxPkgManager picks the package manager based on distro family
// and verifies the binary is available on PATH.
func resolveLinuxPkgManager(distro Distro) PackageManager {
	candidates := linuxPkgCandidates(distro)
	for _, pm := range candidates {
		if commandExists(string(pm)) {
			return pm
		}
	}
	return PkgNone
}

// linuxPkgCandidates returns the ordered list of package managers
// to try for a given distro.
func linuxPkgCandidates(distro Distro) []PackageManager {
	switch distro {
	case DistroUbuntu, DistroDebian, DistroMint, DistroPopOS:
		return []PackageManager{PkgApt}
	case DistroArch, DistroManjaro:
		return []PackageManager{PkgPacman}
	case DistroFedora:
		return []PackageManager{PkgDnf}
	default:
		// Unknown distro: try all common ones
		return []PackageManager{PkgApt, PkgPacman, PkgDnf}
	}
}

// commandExists checks whether a binary is available on the system PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// String returns a human-readable summary of the platform.
func (p *Platform) String() string {
	if p.OS == OSLinux {
		return fmt.Sprintf("%s/%s (%s) [%s]", p.OS, p.Distro, p.Arch, p.PackageManager)
	}
	return fmt.Sprintf("%s (%s) [%s]", p.OS, p.Arch, p.PackageManager)
}

// IsSupported returns true if the platform is one we can install on.
func (p *Platform) IsSupported() bool {
	return p.OS != OSUnknown && p.PackageManager != PkgNone
}