package bundler

import "strings"

// List of supported bundlers
func GetSupportedBundlers() []string {
	return []string{
		"transeptor",
	}
}

// Check if the bundler is supported
func CheckBundler(bundler string) bool {
	for _, b := range GetSupportedBundlers() {
		if b == strings.ToLower(bundler) {
			return true
		}
	}
	return false
}
