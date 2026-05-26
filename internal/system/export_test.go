// Package system — test helpers.
// This file exports internal functions so the _test package can access them
// without making them part of the public API.
package system

// ParseOSReleaseFile is the exported version of parseOSRelease,
// exposed only for testing purposes.
func ParseOSReleaseFile(path string) (map[string]string, error) {
	return parseOSRelease(path)
}