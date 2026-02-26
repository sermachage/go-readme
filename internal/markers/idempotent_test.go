package markers_test

import (
	"strings"
	"testing"

	"github.com/sermachage/go-readme/internal/markers"
)

func TestReplace_EmptyExisting(t *testing.T) {
	got := markers.Replace("", "generated content")
	if !strings.Contains(got, markers.StartMarker) {
		t.Error("expected StartMarker in output")
	}
	if !strings.Contains(got, markers.EndMarker) {
		t.Error("expected EndMarker in output")
	}
	if !strings.Contains(got, "generated content") {
		t.Error("expected generated content in output")
	}
}

func TestReplace_ExistingWithMarkers(t *testing.T) {
	existing := "# My Project\n\n" +
		markers.StartMarker + "\nold content\n" + markers.EndMarker + "\n\nCustom section"

	got := markers.Replace(existing, "new content")

	if strings.Contains(got, "old content") {
		t.Error("old content should have been replaced")
	}
	if !strings.Contains(got, "new content") {
		t.Error("new content should be present")
	}
	if !strings.Contains(got, "# My Project") {
		t.Error("prefix before markers should be preserved")
	}
	if !strings.Contains(got, "Custom section") {
		t.Error("content after markers should be preserved")
	}
}

func TestReplace_ExistingWithoutMarkers(t *testing.T) {
	existing := "# Existing README\n\nSome manual content."

	got := markers.Replace(existing, "generated section")

	if !strings.Contains(got, "# Existing README") {
		t.Error("original content should be preserved")
	}
	if !strings.Contains(got, "generated section") {
		t.Error("generated content should be appended")
	}
	if !strings.Contains(got, markers.StartMarker) {
		t.Error("StartMarker should be added")
	}
}

func TestReplace_Idempotent(t *testing.T) {
	generated := "auto-generated content"
	first := markers.Replace("", generated)
	second := markers.Replace(first, generated)

	if first != second {
		t.Errorf("Replace should be idempotent:\nfirst: %q\nsecond: %q", first, second)
	}
}

func TestExtract(t *testing.T) {
	content := markers.StartMarker + "\nhello world\n" + markers.EndMarker
	got := markers.Extract(content)
	if got != "hello world" {
		t.Errorf("Extract = %q, want %q", got, "hello world")
	}
}

func TestExtract_NoMarkers(t *testing.T) {
	got := markers.Extract("no markers here")
	if got != "" {
		t.Errorf("Extract = %q, want empty string", got)
	}
}
