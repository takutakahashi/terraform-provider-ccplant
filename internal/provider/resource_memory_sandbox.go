package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var (
	_ resource.Resource                = (*memoryResource)(nil)
	_ resource.ResourceWithConfigure   = (*memoryResource)(nil)
	_ resource.ResourceWithImportState = (*memoryResource)(nil)
	_ resource.Resource                = (*sandboxPolicyResource)(nil)
	_ resource.ResourceWithConfigure   = (*sandboxPolicyResource)(nil)
	_ resource.ResourceWithImportState = (*sandboxPolicyResource)(nil)
)

type memoryResource struct{ baseResource }

type memoryModel struct {
	baseModel
	Title     types.String `tfsdk:"title"`
	Content   types.String `tfsdk:"content"`
	Scope     types.String `tfsdk:"scope"`
	TeamID    types.String `tfsdk:"team_id"`
	Tags      types.Map    `tfsdk:"tags"`
	OwnerID   types.String `tfsdk:"owner_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

type memoryPayload struct {
	ID        string            `json:"id,omitempty"`
	Title     string            `json:"title,omitempty"`
	Content   string            `json:"content,omitempty"`
	Scope     string            `json:"scope,omitempty"`
	TeamID    string            `json:"team_id,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	OwnerID   string            `json:"owner_id,omitempty"`
	CreatedAt string            `json:"created_at,omitempty"`
	UpdatedAt string            `json:"updated_at,omitempty"`
}

type memoryUpdatePayload struct {
	Title   string             `json:"title,omitempty"`
	Content string             `json:"content,omitempty"`
	Tags    *map[string]string `json:"tags,omitempty"`
}

func newMemoryResource() resource.Resource {
	return &memoryResource{baseResource: baseResource{typeName: "memory", basePath: "/memories"}}
}

func (r *memoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *memoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *memoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy memory entry.",
		Attributes: map[string]schema.Attribute{
			"id":            idAttribute(),
			"response_json": responseJSONAttribute(),
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Memory title.",
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Memory content.",
			},
			"scope": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Ownership scope: `user` or `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"team_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Team identifier. Required by agentapi-proxy when scope is `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"tags": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Tags attached to the memory. Updates replace the complete tag map.",
			},
			"owner_id":   computedStringAttribute("User ID that owns the memory."),
			"created_at": computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at": computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *memoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan memoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := marshalRequest(memoryPayload{
		Title:   plan.Title.ValueString(),
		Content: plan.Content.ValueString(),
		Scope:   plan.Scope.ValueString(),
		TeamID:  optionalString(plan.TeamID),
		Tags:    stringMapFromTerraform(ctx, plan.Tags, &resp.Diagnostics),
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.CreateJSON(ctx, r.basePath, body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *memoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state memoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := ensureID(state.ID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	responseBody, err := r.client.GetJSON(ctx, r.itemPath(id))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *memoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan memoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tags := optionalStringMapPointer(ctx, plan.Tags, &resp.Diagnostics)
	body := marshalRequest(memoryUpdatePayload{
		Title:   plan.Title.ValueString(),
		Content: plan.Content.ValueString(),
		Tags:    tags,
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *memoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state memoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *memoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *memoryResource) applyResponse(ctx context.Context, body []byte, model *memoryModel, diags *diag.Diagnostics) {
	var payload memoryPayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Title = types.StringValue(payload.Title)
	model.Content = types.StringValue(payload.Content)
	model.Scope = types.StringValue(payload.Scope)
	if payload.TeamID == "" {
		model.TeamID = types.StringNull()
	} else {
		model.TeamID = types.StringValue(payload.TeamID)
	}
	model.Tags = mapToTerraformPreserveEmpty(ctx, payload.Tags, model.Tags, diags)
	model.OwnerID = types.StringValue(payload.OwnerID)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}

type sandboxPolicyResource struct{ baseResource }

type sandboxPolicyModel struct {
	baseModel
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	AllowedDomains types.List   `tfsdk:"allowed_domains"`
	DeniedDomains  types.List   `tfsdk:"denied_domains"`
	CountMode      types.Bool   `tfsdk:"count_mode"`
	Scope          types.String `tfsdk:"scope"`
	TeamID         types.String `tfsdk:"team_id"`
	OwnerID        types.String `tfsdk:"owner_id"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

type sandboxPolicyPayload struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	AllowedDomains []string `json:"allowed_domains,omitempty"`
	DeniedDomains  []string `json:"denied_domains,omitempty"`
	CountMode      bool     `json:"count_mode,omitempty"`
	Scope          string   `json:"scope,omitempty"`
	TeamID         string   `json:"team_id,omitempty"`
	OwnerID        string   `json:"owner_id,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

type sandboxPolicyUpdatePayload struct {
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	AllowedDomains []string `json:"allowed_domains"`
	DeniedDomains  []string `json:"denied_domains"`
	CountMode      bool     `json:"count_mode"`
}

func newSandboxPolicyResource() resource.Resource {
	return &sandboxPolicyResource{baseResource: baseResource{typeName: "sandbox_policy", basePath: "/sandbox-policies"}}
}

func (r *sandboxPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *sandboxPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *sandboxPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy sandbox policy.",
		Attributes: map[string]schema.Attribute{
			"id":            idAttribute(),
			"response_json": responseJSONAttribute(),
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Sandbox policy name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Sandbox policy description.",
			},
			"allowed_domains": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Domains allowed by the sandbox policy.",
			},
			"denied_domains": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Domains denied by the sandbox policy.",
			},
			"count_mode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the policy counts matches instead of enforcing them.",
			},
			"scope": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Ownership scope: `user` or `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"team_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Team identifier. Required by agentapi-proxy when scope is `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"owner_id":   computedStringAttribute("User ID that owns the sandbox policy."),
			"created_at": computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at": computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *sandboxPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sandboxPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := marshalRequest(sandboxPolicyPayload{
		Name:           plan.Name.ValueString(),
		Description:    optionalString(plan.Description),
		AllowedDomains: stringListFromTerraform(ctx, plan.AllowedDomains, &resp.Diagnostics),
		DeniedDomains:  stringListFromTerraform(ctx, plan.DeniedDomains, &resp.Diagnostics),
		CountMode:      !plan.CountMode.IsNull() && !plan.CountMode.IsUnknown() && plan.CountMode.ValueBool(),
		Scope:          plan.Scope.ValueString(),
		TeamID:         optionalString(plan.TeamID),
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.CreateJSON(ctx, r.basePath, body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sandboxPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sandboxPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := ensureID(state.ID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	responseBody, err := r.client.GetJSON(ctx, r.itemPath(id))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sandboxPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sandboxPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := marshalRequest(sandboxPolicyUpdatePayload{
		Name:           plan.Name.ValueString(),
		Description:    optionalString(plan.Description),
		AllowedDomains: stringListFromTerraform(ctx, plan.AllowedDomains, &resp.Diagnostics),
		DeniedDomains:  stringListFromTerraform(ctx, plan.DeniedDomains, &resp.Diagnostics),
		CountMode:      !plan.CountMode.IsNull() && !plan.CountMode.IsUnknown() && plan.CountMode.ValueBool(),
	}, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sandboxPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sandboxPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *sandboxPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *sandboxPolicyResource) applyResponse(ctx context.Context, body []byte, model *sandboxPolicyModel, diags *diag.Diagnostics) {
	var payload sandboxPolicyPayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Name = types.StringValue(payload.Name)
	model.Description = types.StringValue(payload.Description)
	model.AllowedDomains = listToTerraformPreserveEmpty(ctx, payload.AllowedDomains, model.AllowedDomains, diags)
	model.DeniedDomains = listToTerraformPreserveEmpty(ctx, payload.DeniedDomains, model.DeniedDomains, diags)
	model.CountMode = types.BoolValue(payload.CountMode)
	model.Scope = types.StringValue(payload.Scope)
	if payload.TeamID == "" {
		model.TeamID = types.StringNull()
	} else {
		model.TeamID = types.StringValue(payload.TeamID)
	}
	model.OwnerID = types.StringValue(payload.OwnerID)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}
