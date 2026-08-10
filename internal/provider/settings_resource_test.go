package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func TestValidateSettingsScope(t *testing.T) {
	for _, scope := range []string{"user", "team"} {
		var diagnostics diag.Diagnostics
		if !validateSettingsScope(types.StringValue(scope), &diagnostics) {
			t.Fatalf("validateSettingsScope(%q) returned false: %v", scope, diagnostics)
		}
	}

	var diagnostics diag.Diagnostics
	if validateSettingsScope(types.StringValue("organization"), &diagnostics) {
		t.Fatal("validateSettingsScope(organization) returned true")
	}
	if !diagnostics.HasError() {
		t.Fatal("validateSettingsScope(organization) did not return an error diagnostic")
	}
}
