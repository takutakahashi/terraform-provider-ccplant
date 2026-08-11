package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestNormalizeJSON(t *testing.T) {
	got := normalizeJSON([]byte(`{
		"ok": true
	}`))
	if want := `{"ok":true}`; got != want {
		t.Fatalf("normalizeJSON() = %q, want %q", got, want)
	}
}

func TestResourceResponseJSONIsSensitive(t *testing.T) {
	resources := map[string]resource.Resource{
		"webhook":         newWebhookResource(),
		"schedule":        newScheduleResource(),
		"slackbot":        newSlackbotResource(),
		"memory":          newMemoryResource(),
		"session_profile": newSessionProfileResource(),
		"sandbox_policy":  newSandboxPolicyResource(),
	}

	for name, providerResource := range resources {
		t.Run(name, func(t *testing.T) {
			var response resource.SchemaResponse
			providerResource.Schema(context.Background(), resource.SchemaRequest{}, &response)
			if diagnostics := response.Schema.ValidateImplementation(context.Background()); diagnostics.HasError() {
				t.Fatalf("schema validation failed: %v", diagnostics)
			}

			attribute, ok := response.Schema.Attributes["response_json"].(schema.StringAttribute)
			if !ok {
				t.Fatalf("response_json has type %T, want schema.StringAttribute", response.Schema.Attributes["response_json"])
			}
			if !attribute.Sensitive {
				t.Error("response_json must be sensitive")
			}
		})
	}
}

func TestSessionEnvironmentIsSensitive(t *testing.T) {
	tests := []struct {
		name             string
		providerResource resource.Resource
		configAttribute  string
	}{
		{name: "webhook", providerResource: newWebhookResource(), configAttribute: "session_config"},
		{name: "schedule", providerResource: newScheduleResource(), configAttribute: "session_config"},
		{name: "slackbot", providerResource: newSlackbotResource(), configAttribute: "session_config"},
		{name: "session_profile", providerResource: newSessionProfileResource(), configAttribute: "config"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response resource.SchemaResponse
			test.providerResource.Schema(context.Background(), resource.SchemaRequest{}, &response)

			config, ok := response.Schema.Attributes[test.configAttribute].(schema.SingleNestedAttribute)
			if !ok {
				t.Fatalf("%s has type %T, want schema.SingleNestedAttribute", test.configAttribute, response.Schema.Attributes[test.configAttribute])
			}
			environment, ok := config.Attributes["environment"].(schema.MapAttribute)
			if !ok {
				t.Fatalf("environment has type %T, want schema.MapAttribute", config.Attributes["environment"])
			}
			if !environment.Sensitive {
				t.Error("environment must be sensitive")
			}
		})
	}
}
