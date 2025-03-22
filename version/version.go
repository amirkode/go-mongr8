/*
Copyright (c) 2023-present the go-mongr8 Authors and Contributors
[@see Authors file]

Licensed under the MIT License
(https://opensource.org/licenses/MIT)
*/
package version

const (
	Major = 0
	Minor = 1
	Patch = 1

	// version name
	Version = "v0.1.1"
)

// returns the full version string with additional info
func String() string {
	return Version
}

// checks if the current version is at least the specified version
func IsAtLeast(major, minor, patch int) bool {
	if Major > major {
		return true
	}
	if Major < major {
		return false
	}
	if Minor > minor {
		return true
	}
	if Minor < minor {
		return false
	}
	return Patch >= patch
}
