package provider

import "testing"

func TestSettingsResourceItemPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "alice", want: "/settings/alice"},
		{name: "ccplant/platform", want: "/settings/ccplant%2Fplatform"},
		{name: "team with spaces", want: "/settings/team%20with%20spaces"},
	}

	r := &settingsResource{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.itemPath(tt.name); got != tt.want {
				t.Fatalf("itemPath(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
