package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func TestSettingsRequestJSON(t *testing.T) {
	model := settingsResourceModel{
		Scope:          types.StringValue("team"),
		Name:           types.StringValue("ccplant/platform"),
		AuthMode:       types.StringValue("bedrock"),
		EnabledPlugins: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("commit@example")}),
		EnvVars: types.MapValueMust(types.StringType, map[string]attr.Value{
			"MANAGED_BY": types.StringValue("terraform"),
		}),
	}

	var diagnostics diag.Diagnostics
	body, ok := settingsRequestJSON(model, &diagnostics)
	if !ok || diagnostics.HasError() {
		t.Fatalf("settingsRequestJSON() diagnostics = %v", diagnostics)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got["auth_mode"] != "bedrock" {
		t.Fatalf("auth_mode = %v, want bedrock", got["auth_mode"])
	}
	env := got["env_vars"].(map[string]any)
	if env["MANAGED_BY"] != "terraform" {
		t.Fatalf("env_vars.MANAGED_BY = %v, want terraform", env["MANAGED_BY"])
	}
}

func TestValidateSettingsScope(t *testing.T) {
	for _, scope := range []string{"user", "team"} {
		var diagnostics diag.Diagnostics
		model := settingsResourceModel{Scope: types.StringValue(scope)}
		if !validateSettingsModel(model, &diagnostics) {
			t.Fatalf("validateSettingsModel(%q) returned false: %v", scope, diagnostics)
		}
	}

	var diagnostics diag.Diagnostics
	model := settingsResourceModel{Scope: types.StringValue("organization")}
	if validateSettingsModel(model, &diagnostics) {
		t.Fatal("validateSettingsModel(organization) returned true")
	}
	if !diagnostics.HasError() {
		t.Fatal("validateSettingsModel(organization) did not return an error diagnostic")
	}
}
