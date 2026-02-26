// Package markers handles idempotent README section management.
package markers

import (
	"strings"
)

const (
	// StartMarker is placed immediately before generated content.
	StartMarker = "<!-- readmeaker:start -->"
	// EndMarker is placed immediately after generated content.
	EndMarker = "<!-- readmeaker:end -->"
)

// Replace updates the managed section inside existing README content.
// If markers are present, only the content between them is replaced.
// If no markers are present, the generated section is appended.
// If existing is empty, generated is returned as-is (with markers).
func Replace(existing, generated string) string {
	wrapped := wrap(generated)

	if existing == "" {
		return wrapped
	}

	start := strings.Index(existing, StartMarker)
	end := strings.Index(existing, EndMarker)

	if start >= 0 && end >= 0 && end > start {
		before := existing[:start]
		after := existing[end+len(EndMarker):]
		return before + wrapped + after
	}

	// No markers found – append the managed section.
	return strings.TrimRight(existing, "\n") + "\n\n" + wrapped + "\n"
}

// wrap surrounds content with the start/end markers.
func wrap(content string) string {
	return StartMarker + "\n" + strings.TrimSpace(content) + "\n" + EndMarker
}

// Extract returns only the content between the markers, or "" if none found.
func Extract(content string) string {
	start := strings.Index(content, StartMarker)
	end := strings.Index(content, EndMarker)
	if start < 0 || end < 0 || end <= start {
		return ""
	}
	inner := content[start+len(StartMarker) : end]
	return strings.TrimSpace(inner)
}
