package provider

import "testing"

func TestNormalizeJSON(t *testing.T) {
	got := normalizeJSON([]byte(`{
		"ok": true
	}`))
	if want := `{"ok":true}`; got != want {
		t.Fatalf("normalizeJSON() = %q, want %q", got, want)
	}
}
