package provider

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var (
	_ resource.Resource                = (*settingsResource)(nil)
	_ resource.ResourceWithConfigure   = (*settingsResource)(nil)
	_ resource.ResourceWithImportState = (*settingsResource)(nil)
)

type settingsResource struct {
	client *client.Client
}

type settingsResourceModel struct {
	Scope                   types.String `tfsdk:"scope"`
	Name                    types.String `tfsdk:"name"`
	Bedrock                 types.Object `tfsdk:"bedrock"`
	MCPServers              types.Map    `tfsdk:"mcp_servers"`
	Marketplaces            types.Map    `tfsdk:"marketplaces"`
	ClaudeCodeOAuthToken    types.String `tfsdk:"claude_code_oauth_token"`
	AuthMode                types.String `tfsdk:"auth_mode"`
	EnabledPlugins          types.Set    `tfsdk:"enabled_plugins"`
	EnvVars                 types.Map    `tfsdk:"env_vars"`
	PreferredTeamID         types.String `tfsdk:"preferred_team_id"`
	GitHubAppInstallationID types.String `tfsdk:"github_app_installation_id"`
	SlackUserID             types.String `tfsdk:"slack_user_id"`
	NotificationChannels    types.Set    `tfsdk:"notification_channels"`
	ExternalSessionManagers types.List   `tfsdk:"external_session_managers"`
	DefaultSessionProfileID types.String `tfsdk:"default_session_profile_id"`
	HasClaudeCodeOAuthToken types.Bool   `tfsdk:"has_claude_code_oauth_token"`
	EnvVarKeys              types.Set    `tfsdk:"env_var_keys"`
	CreatedAt               types.String `tfsdk:"created_at"`
	UpdatedAt               types.String `tfsdk:"updated_at"`
}

func newSettingsResource() resource.Resource { return &settingsResource{} }

func (r *settingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *settingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages typed user- or team-scoped agentapi-proxy settings.",
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{Required: true, MarkdownDescription: "Ownership scope. Must be `user` or `team`.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":  schema.StringAttribute{Required: true, MarkdownDescription: "User ID or org/team-slug owner name.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"bedrock": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
				"enabled":           schema.BoolAttribute{Required: true},
				"model":             schema.StringAttribute{Optional: true},
				"access_key_id":     schema.StringAttribute{Optional: true, Sensitive: true},
				"secret_access_key": schema.StringAttribute{Optional: true, Sensitive: true},
				"role_arn":          schema.StringAttribute{Optional: true},
				"profile":           schema.StringAttribute{Optional: true},
			}},
			"mcp_servers": schema.MapNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"type":    schema.StringAttribute{Required: true},
				"url":     schema.StringAttribute{Optional: true},
				"command": schema.StringAttribute{Optional: true},
				"args":    schema.ListAttribute{Optional: true, ElementType: types.StringType},
				"env":     schema.MapAttribute{Optional: true, Sensitive: true, ElementType: types.StringType},
				"headers": schema.MapAttribute{Optional: true, Sensitive: true, ElementType: types.StringType},
			}}},
			"marketplaces": schema.MapNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{Required: true},
			}}},
			"claude_code_oauth_token":    schema.StringAttribute{Optional: true, Sensitive: true},
			"auth_mode":                  schema.StringAttribute{Optional: true, MarkdownDescription: "Authentication mode: `oauth` or `bedrock`."},
			"enabled_plugins":            schema.SetAttribute{Optional: true, ElementType: types.StringType},
			"env_vars":                   schema.MapAttribute{Optional: true, Sensitive: true, ElementType: types.StringType},
			"preferred_team_id":          schema.StringAttribute{Optional: true},
			"github_app_installation_id": schema.StringAttribute{Optional: true, Sensitive: true},
			"slack_user_id":              schema.StringAttribute{Optional: true},
			"notification_channels":      schema.SetAttribute{Optional: true, ElementType: types.StringType},
			"external_session_managers": schema.ListNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"id":          schema.StringAttribute{Required: true},
				"instance_id": schema.StringAttribute{Optional: true},
				"name":        schema.StringAttribute{Required: true},
				"default":     schema.BoolAttribute{Optional: true},
				"labels":      schema.MapAttribute{Optional: true, ElementType: types.StringType},
			}}},
			"default_session_profile_id":  schema.StringAttribute{Optional: true},
			"has_claude_code_oauth_token": schema.BoolAttribute{Computed: true},
			"env_var_keys":                schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"created_at":                  schema.StringAttribute{Computed: true},
			"updated_at":                  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *settingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *client.Client.")
		return
	}
	r.client = apiClient
}

func (r *settingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan settingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || !validateSettingsModel(plan, &resp.Diagnostics) {
		return
	}
	body, ok := settingsRequestJSON(plan, &resp.Diagnostics)
	if !ok {
		return
	}
	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.Name.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	applySettingsResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state settingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.GetJSON(ctx, r.itemPath(state.Name.ValueString()))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	applySettingsResponse(ctx, responseBody, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *settingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan settingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || !validateSettingsModel(plan, &resp.Diagnostics) {
		return
	}
	body, ok := settingsRequestJSON(plan, &resp.Diagnostics)
	if !ok {
		return
	}
	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.Name.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	applySettingsResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state settingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, r.itemPath(state.Name.ValueString())); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *settingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || (parts[0] != "user" && parts[0] != "team") || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import identifier", "Use scope:name, for example user:alice or team:ccplant/platform.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func validateSettingsModel(model settingsResourceModel, diagnostics *diag.Diagnostics) bool {
	if model.Scope.IsNull() || model.Scope.IsUnknown() || (model.Scope.ValueString() != "user" && model.Scope.ValueString() != "team") {
		diagnostics.AddAttributeError(path.Root("scope"), "Invalid settings scope", "scope must be either user or team.")
		return false
	}
	if !model.AuthMode.IsNull() && !model.AuthMode.IsUnknown() && model.AuthMode.ValueString() != "oauth" && model.AuthMode.ValueString() != "bedrock" {
		diagnostics.AddAttributeError(path.Root("auth_mode"), "Invalid authentication mode", "auth_mode must be either oauth or bedrock.")
		return false
	}
	return true
}

func settingsRequestJSON(model settingsResourceModel, diagnostics *diag.Diagnostics) ([]byte, bool) {
	body := map[string]any{}
	putString(body, "claude_code_oauth_token", model.ClaudeCodeOAuthToken)
	putString(body, "auth_mode", model.AuthMode)
	putString(body, "preferred_team_id", model.PreferredTeamID)
	putString(body, "github_app_installation_id", model.GitHubAppInstallationID)
	putString(body, "slack_user_id", model.SlackUserID)
	putString(body, "default_session_profile_id", model.DefaultSessionProfileID)
	putStringSet(body, "enabled_plugins", model.EnabledPlugins)
	putStringSet(body, "notification_channels", model.NotificationChannels)
	putStringMap(body, "env_vars", model.EnvVars)
	putObject(body, "bedrock", model.Bedrock)
	putObjectMap(body, "mcp_servers", model.MCPServers)
	putObjectMap(body, "marketplaces", model.Marketplaces)
	putObjectList(body, "external_session_managers", model.ExternalSessionManagers)
	result, err := json.Marshal(body)
	if err != nil {
		diagnostics.AddError("Encode settings failed", err.Error())
		return nil, false
	}
	return result, true
}

func putString(body map[string]any, key string, value types.String) {
	if !value.IsNull() && !value.IsUnknown() {
		body[key] = value.ValueString()
	}
}

func putStringSet(body map[string]any, key string, value types.Set) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	items := make([]string, 0, len(value.Elements()))
	for _, item := range value.Elements() {
		items = append(items, item.(types.String).ValueString())
	}
	body[key] = items
}

func putStringMap(body map[string]any, key string, value types.Map) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	items := map[string]string{}
	for name, item := range value.Elements() {
		items[name] = item.(types.String).ValueString()
	}
	body[key] = items
}

func putObject(body map[string]any, key string, value types.Object) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	body[key] = attrMap(value.Attributes())
}

func putObjectMap(body map[string]any, key string, value types.Map) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	items := map[string]any{}
	for name, item := range value.Elements() {
		items[name] = attrMap(item.(types.Object).Attributes())
	}
	body[key] = items
}

func putObjectList(body map[string]any, key string, value types.List) {
	if value.IsNull() || value.IsUnknown() {
		return
	}
	items := make([]any, 0, len(value.Elements()))
	for _, item := range value.Elements() {
		items = append(items, attrMap(item.(types.Object).Attributes()))
	}
	body[key] = items
}

func attrMap(values map[string]attr.Value) map[string]any {
	result := map[string]any{}
	for key, value := range values {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		switch v := value.(type) {
		case types.String:
			result[key] = v.ValueString()
		case types.Bool:
			result[key] = v.ValueBool()
		case types.List:
			items := make([]string, 0, len(v.Elements()))
			for _, item := range v.Elements() {
				items = append(items, item.(types.String).ValueString())
			}
			result[key] = items
		case types.Map:
			items := map[string]string{}
			for name, item := range v.Elements() {
				items[name] = item.(types.String).ValueString()
			}
			result[key] = items
		}
	}
	return result
}

func applySettingsResponse(ctx context.Context, body []byte, model *settingsResourceModel, diagnostics *diag.Diagnostics) {
	var response struct {
		HasToken bool     `json:"has_claude_code_oauth_token"`
		EnvKeys  []string `json:"env_var_keys"`
		Created  string   `json:"created_at"`
		Updated  string   `json:"updated_at"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		diagnostics.AddError("Decode settings response failed", err.Error())
		return
	}
	model.HasClaudeCodeOAuthToken = types.BoolValue(response.HasToken)
	model.CreatedAt = types.StringValue(response.Created)
	model.UpdatedAt = types.StringValue(response.Updated)
	set, diags := types.SetValueFrom(ctx, types.StringType, response.EnvKeys)
	diagnostics.Append(diags...)
	model.EnvVarKeys = set
}

func (r *settingsResource) itemPath(name string) string { return "/settings/" + url.PathEscape(name) }
