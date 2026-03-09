// Package markers handles idempotent README section management.
package markers

import (
	"strings"
)

const (
	// StartMarker is placed immediately before generated content.
	StartMarker = "<!-- go-readme:start -->"
	// EndMarker is placed immediately after generated content.
	EndMarker = "<!-- go-readme:end -->"

	// LegacyStartMarker is kept for backward compatibility during marker migration.
	LegacyStartMarker = "<!-- readmeaker:start -->"
	// LegacyEndMarker is kept for backward compatibility during marker migration.
	LegacyEndMarker = "<!-- readmeaker:end -->"
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

	start, end, endLen, found := findManagedSection(existing)
	if found {
		before := existing[:start]
		after := existing[end+endLen:]
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
	start, end, _, found := findManagedSection(content)
	if !found {
		return ""
	}

	markerLen := len(StartMarker)
	if strings.HasPrefix(content[start:], LegacyStartMarker) {
		markerLen = len(LegacyStartMarker)
	}

	inner := content[start+markerLen : end]
	return strings.TrimSpace(inner)
}

func findManagedSection(content string) (start int, end int, endLen int, found bool) {
	start = strings.Index(content, StartMarker)
	end = strings.Index(content, EndMarker)
	if start >= 0 && end >= 0 && end > start {
		return start, end, len(EndMarker), true
	}

	start = strings.Index(content, LegacyStartMarker)
	end = strings.Index(content, LegacyEndMarker)
	if start >= 0 && end >= 0 && end > start {
		return start, end, len(LegacyEndMarker), true
	}

	return 0, 0, 0, false
}
