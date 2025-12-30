package paths

import (
	"regexp"
	"slices"
	"strings"
)

func ExistsInPaths(path string, paths []string) bool {
	return slices.ContainsFunc(paths, func(s string) bool {
		if MatchesWildcard(path, s) {
			return true
		}
		return false
	})
}

func MatchesWildcard(path, pattern string) bool {
	// Security: Reject paths with path traversal attempts or empty segments
	if containsSecurityRisk(path) || containsSecurityRisk(pattern) {
		return false
	}

	// Special case: root path "/" matches root wildcard "/*"
	if path == "/" && pattern == "/*" {
		return true
	}

	// Normalize trailing slashes for comparison
	normPath := normalizeTrailingSlash(path)
	normPattern := normalizeTrailingSlash(pattern)

	// Validate pattern and path contain only safe characters
	if !containsOnlyValidChars(normPattern, isValidPatternChar) {
		return false
	}
	if !containsOnlyValidChars(normPath, isValidPathChar) {
		return false
	}

	// If pattern doesn't contain wildcard, require exact match
	if !strings.Contains(normPattern, "*") {
		return normPath == normPattern
	}

	// Convert wildcard pattern to regex:
	// 1. Escape special regex characters
	// 2. Replace \* with .+ to match at least one character
	// 3. Anchor to match entire string
	regexPattern := regexp.QuoteMeta(normPattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\*`, `.+`)
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, normPath)
	return err == nil && matched
}

// containsSecurityRisk checks for path traversal attempts and empty segments
func containsSecurityRisk(s string) bool {
	return strings.Contains(s, "..") || strings.Contains(s, "//")
}

// normalizeTrailingSlash removes trailing slash for consistent comparison
func normalizeTrailingSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}

// isValidPathChar returns true if the character is allowed in paths
func isValidPathChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '/' || ch == '-' || ch == '_' || ch == '.' || ch == '%'
}

// isValidPatternChar returns true if the character is allowed in patterns
func isValidPatternChar(ch rune) bool {
	return isValidPathChar(ch) || ch == '*'
}

// containsOnlyValidChars validates that all characters pass the validator function
func containsOnlyValidChars(s string, isValid func(rune) bool) bool {
	for _, ch := range s {
		if !isValid(ch) {
			return false
		}
	}
	return true
}
