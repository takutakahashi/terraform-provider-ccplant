package provider

import "testing"

func TestExtractID(t *testing.T) {
	id, err := extractID([]byte(`{"id":"abc123"}`))
	if err != nil {
		t.Fatalf("extractID() error = %v", err)
	}
	if id != "abc123" {
		t.Fatalf("id = %q, want abc123", id)
	}
}

func TestExtractIDAlternativeKeys(t *testing.T) {
	id, err := extractID([]byte(`{"session_profile_id":"profile123"}`))
	if err != nil {
		t.Fatalf("extractID() error = %v", err)
	}
	if id != "profile123" {
		t.Fatalf("id = %q, want profile123", id)
	}
}

func TestNormalizeJSON(t *testing.T) {
	got := normalizeJSON([]byte(`{
		"ok": true
	}`))
	if want := `{"ok":true}`; got != want {
		t.Fatalf("normalizeJSON() = %q, want %q", got, want)
	}
}
