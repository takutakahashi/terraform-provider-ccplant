package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	_ resource.Resource                = (*jsonResource)(nil)
	_ resource.ResourceWithConfigure   = (*jsonResource)(nil)
	_ resource.ResourceWithImportState = (*jsonResource)(nil)
)

type jsonResource struct {
	typeName string
	basePath string
	client   *client.Client
}

type jsonResourceModel struct {
	ID           types.String `tfsdk:"id"`
	BodyJSON     types.String `tfsdk:"body_json"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func newJSONResource(typeName, basePath string) func() resource.Resource {
	return func() resource.Resource {
		return &jsonResource{typeName: typeName, basePath: basePath}
	}
}

func (r *jsonResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *jsonResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Manages an agentapi-proxy %s resource with JSON request bodies.", strings.ReplaceAll(r.typeName, "_", " ")),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource ID assigned by agentapi-proxy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"body_json": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "JSON request body sent to agentapi-proxy on create and update.",
			},
			"response_json": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Latest JSON response returned by agentapi-proxy.",
			},
		},
	}
}

func (r *jsonResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *jsonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan jsonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, ok := validateJSON(plan.BodyJSON, path.Root("body_json"), &resp.Diagnostics)
	if !ok {
		return
	}

	responseBody, err := r.client.CreateJSON(ctx, r.basePath, body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}

	id, err := extractID(responseBody)
	if err != nil {
		resp.Diagnostics.AddError("Create response missing ID", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.ResponseJSON = types.StringValue(normalizeJSON(responseBody))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jsonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state jsonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	responseBody, err := r.client.GetJSON(ctx, r.itemPath(state.ID.ValueString()))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}

	state.ResponseJSON = types.StringValue(normalizeJSON(responseBody))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *jsonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan jsonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, ok := validateJSON(plan.BodyJSON, path.Root("body_json"), &resp.Diagnostics)
	if !ok {
		return
	}

	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}

	plan.ResponseJSON = types.StringValue(normalizeJSON(responseBody))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *jsonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state jsonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
		return
	}
}

func (r *jsonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jsonResource) itemPath(id string) string {
	return strings.TrimRight(r.basePath, "/") + "/" + id
}

func validateJSON(value types.String, attrPath path.Path, diagnostics *diag.Diagnostics) ([]byte, bool) {
	if value.IsNull() || value.IsUnknown() {
		diagnostics.AddAttributeError(attrPath, "Missing JSON body", "body_json must be known and non-null.")
		return nil, false
	}

	body := []byte(value.ValueString())
	if !json.Valid(body) {
		diagnostics.AddAttributeError(attrPath, "Invalid JSON body", "body_json must contain a valid JSON object.")
		return nil, false
	}

	return body, true
}

func extractID(body []byte) (string, error) {
	var obj map[string]any
	if err := json.Unmarshal(body, &obj); err != nil {
		return "", err
	}

	for _, key := range []string{"id", "memory_id", "session_profile_id", "webhook_id", "schedule_id", "slackbot_id"} {
		if value, ok := obj[key].(string); ok && value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("response must include an id field")
}

func normalizeJSON(body []byte) string {
	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return string(body)
	}

	return string(normalized)
}
