package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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

// settingsResource manages the named /settings/:name endpoint. User and team
// settings share the same API, but expose distinct Terraform attributes so a
// configuration makes the ownership scope explicit.
type settingsResource struct {
	typeName   string
	identifier string
	client     *client.Client
}

func newSettingsResource(typeName, identifier string) func() resource.Resource {
	return func() resource.Resource {
		return &settingsResource{typeName: typeName, identifier: identifier}
	}
}

func (r *settingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *settingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	identifierDescription := "User ID that owns these settings."
	if r.identifier == "team_id" {
		identifierDescription = "Team ID that owns these settings, in org/team-slug format."
	}

	attrs := map[string]schema.Attribute{
		"body_json": schema.StringAttribute{
			Required:            true,
			Sensitive:           true,
			MarkdownDescription: "JSON settings document sent to agentapi-proxy. It is sensitive because it may contain credentials.",
		},
		"response_json": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Latest sanitized JSON response returned by agentapi-proxy.",
		},
		r.identifier: schema.StringAttribute{
			Required:            true,
			MarkdownDescription: identifierDescription,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages agentapi-proxy %s.", strings.ReplaceAll(r.typeName, "_", " ")),
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
	var name, bodyJSON types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(r.identifier), &name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body_json"), &bodyJSON)...)
	if resp.Diagnostics.HasError() {
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(r.identifier), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("body_json"), bodyJSON)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("response_json"), normalizeJSON(responseBody))...)
}

func (r *settingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(r.identifier), &name)...)
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
	var name, bodyJSON types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(r.identifier), &name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body_json"), &bodyJSON)...)
	if resp.Diagnostics.HasError() {
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(r.identifier), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("body_json"), bodyJSON)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("response_json"), normalizeJSON(responseBody))...)
}

func (r *settingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var name types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(r.identifier), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, r.itemPath(name.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *settingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(r.identifier), req.ID)...)
}

func (r *settingsResource) itemPath(name string) string {
	return "/settings/" + url.PathEscape(name)
}
