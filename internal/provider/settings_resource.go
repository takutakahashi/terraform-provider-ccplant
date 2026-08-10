package provider

import (
	"context"
	"net/url"
	"strings"

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

// settingsResource manages the named /settings/:name endpoint. The scope
// attribute makes the ownership explicit and follows other scoped resources.
type settingsResource struct {
	client *client.Client
}

func newSettingsResource() resource.Resource {
	return &settingsResource{}
}

func (r *settingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *settingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"scope": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Ownership scope for the settings. Must be `user` or `team`.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Settings owner name: a user ID for user scope or an org/team-slug ID for team scope.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"body_json": schema.StringAttribute{
			Required:            true,
			Sensitive:           true,
			MarkdownDescription: "JSON settings document sent to agentapi-proxy. It is sensitive because it may contain credentials.",
		},
		"response_json": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Latest sanitized JSON response returned by agentapi-proxy.",
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages user- or team-scoped agentapi-proxy settings.",
		Attributes:          attrs,
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
	var scope, name, bodyJSON types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("scope"), &scope)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body_json"), &bodyJSON)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !validateSettingsScope(scope, &resp.Diagnostics) {
		return
	}

	body, ok := validateJSON(bodyJSON, path.Root("body_json"), &resp.Diagnostics)
	if !ok {
		return
	}

	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(name.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), scope)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("body_json"), bodyJSON)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("response_json"), normalizeJSON(responseBody))...)
}

func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	responseBody, err := r.client.GetJSON(ctx, r.itemPath(name.ValueString()))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("response_json"), normalizeJSON(responseBody))...)
}

func (r *settingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var scope, name, bodyJSON types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("scope"), &scope)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body_json"), &bodyJSON)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !validateSettingsScope(scope, &resp.Diagnostics) {
		return
	}

	body, ok := validateJSON(bodyJSON, path.Root("body_json"), &resp.Diagnostics)
	if !ok {
		return
	}

	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(name.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope"), scope)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("body_json"), bodyJSON)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("response_json"), normalizeJSON(responseBody))...)
}

func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, r.itemPath(name.ValueString()))
	if err != nil && !client.IsNotFound(err) {
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

func validateSettingsScope(scope types.String, diagnostics *diag.Diagnostics) bool {
	if scope.IsNull() || scope.IsUnknown() {
		diagnostics.AddAttributeError(path.Root("scope"), "Missing settings scope", "scope must be known and non-null.")
		return false
	}
	if scope.ValueString() != "user" && scope.ValueString() != "team" {
		diagnostics.AddAttributeError(path.Root("scope"), "Invalid settings scope", "scope must be either user or team.")
		return false
	}
	return true
}

func (r *settingsResource) itemPath(name string) string {
	return "/settings/" + url.PathEscape(name)
}
