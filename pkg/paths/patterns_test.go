package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchesWildcard(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wildcard string
		expected bool
	}{
		// Basic prefix matching (nginx default behavior)
		{"exact prefix match", "/api/users", "/api/*", true},
		{"deep path match", "/api/users/123/details", "/api/*", true},
		{"root wildcard matches everything", "/anything/deep/path", "/*", true},
		{"no match different prefix", "/api/users", "/static/*", false},

		// Exact matching without wildcard
		{"no wildcard requires exact match", "/api/users", "/api/", false},
		{"exact match without wildcard", "/api/users", "/api/users", true},

		// Trailing slash normalization (nginx normalizes trailing slashes)
		{"trailing slash normalized on path", "/api/users/", "/api/users", true},
		{"trailing slash normalized on pattern", "/api/users", "/api/users/", true},
		{"both have trailing slash", "/api/users/", "/api/users/", true},

		// Wildcard at different positions
		{"wildcard at end", "/api/users/123", "/api/*", true},
		{"intermediate wildcard", "/api/v1/users", "/api/*/users", true},
		{"multiple segments after wildcard", "/api/users/123/profile", "/api/*/profile", true},

		// Edge cases
		{"root path with root wildcard", "/", "/*", true},
		{"root path exact", "/", "/", true},
		{"empty segments rejected", "//api//users", "/api/users", false},

		// Longest prefix wins (multiple wildcards)
		{"nested wildcard", "/static/css/main.css", "/static/*", true},
		{"double wildcard", "/a/b/c/d", "/*/c/*", true},

		// Security: path traversal attempts
		{"path traversal blocked in path", "/api/../admin", "/api/*", false},
		{"path traversal blocked in pattern", "/api/users", "/api/../*", false},

		// Case sensitivity (nginx is case-sensitive by default)
		{"case sensitive path", "/API/users", "/api/*", false},
		{"case sensitive pattern", "/api/users", "/API/*", false},

		// Special characters in paths
		{"encoded space", "/api/my%20file", "/api/*", true},
		{"dash in path", "/api/user-123", "/api/*", true},
		{"underscore in path", "/api/user_123", "/api/*", true},
		{"dot in path", "/api/file.txt", "/api/*", true},
		{"contains invalid character in path", "/api/file.txt?", "/api/*", false},
		{"contains invalid character in pattern", "/api/user", "/api/*?", false},

		// Greedy wildcard behavior
		{"wildcard is greedy", "/api/v1/v2/users", "/api/*/users", true},
		{"wildcard matches empty rejected", "/api//users", "/api/*/users", false},

		// No match scenarios
		{"shorter path than pattern", "/api", "/api/users/*", false},
		{"different path segment", "/api/posts", "/api/users/*", false},
		{"missing middle segment", "/api/123", "/api/users/*/profile", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchesWildcard(tt.path, tt.wildcard)
			assert.Equal(t, tt.expected, result,
				"matchesWildcard(%q, %q) = %v, want %v",
				tt.path, tt.wildcard, result, tt.expected)
		})
	}
}

func TestExistsInPaths(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		paths    []string
		expected bool
	}{
		{"exact match", "/api/users", []string{"/api/users", "/api/*"}, true},
		{"wildcard match", "/api/users", []string{"/api/*"}, true},
		{"no match", "/api/users", []string{"/static/*"}, false},
		{"empty paths", "/api/users", []string{}, false},
		{"empty path", "", []string{"/api/users"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExistsInPaths(tt.path, tt.paths)
			assert.Equal(t, tt.expected, result,
				"existsInPaths(%q, %v) = %v, want %v",
				tt.path, tt.paths, result, tt.expected)
		})
	}
}
